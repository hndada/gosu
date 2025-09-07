// The following code is generated from ChatGPT.

// Package o2jam provides a minimal OJN (O2Jam note) file parser.
//
// It implements enough of the public community docs to extract:
//   - Song metadata (title/artist/noter/genre/BPM/levels/time)
//   - Per-difficulty note sections (packages, events)
//   - Cover/JPEG and optional BMP thumbnail offsets/sizes
//
// References consulted while implementing this:
//   - open2jam "The OJN Documentation" (header/notes layout)
//   - open2jam "The Notes section" (channels, event layout)
//
// This parser aims to be tolerant: unknown channels are preserved as raw bytes.
// Strings inside OJN are typically CP949/KS X 1001 encoded; we try to decode
// them to UTF-8. If decode fails, you still get the raw bytes as fallback.
package o2jam

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// File signature present in header.signature
var ojnMagic = [4]byte{'o', 'j', 'n', 0x00}

// Genre codes as per community docs.
const (
	GenreBallad      = 0
	GenreRock        = 1
	GenreDance       = 2
	GenreTechno      = 3
	GenreHipHop      = 4
	GenreSoulRNB     = 5
	GenreJazz        = 6
	GenreFunk        = 7
	GenreClassical   = 8
	GenreTraditional = 9
	GenreEtc         = 10
)

// Channel identifiers
const (
	ChanMeasureFraction = 0
	ChanBPM             = 1
	ChanLane1           = 2
	ChanLane2           = 3
	ChanLane3           = 4
	ChanLane4           = 5 // middle
	ChanLane5           = 6
	ChanLane6           = 7
	ChanLane7           = 8
	ChanAutoplayStart   = 9 // 9..22 auto-play samples
)

// Note types
const (
	NoteNormal   = 0
	NoteLongHead = 2
	NoteLongTail = 3
	NoteOggFlag  = 4 // indicates Notetool OGG use (compat only)
)

// Header represents the fixed-size header at the start of an OJN file.
// Layout mirrors the community struct; all numeric fields are little-endian.
// Some historical/"old_*" fields are retained for completeness.
type Header struct {
	SongID        uint32
	Signature     [4]byte // must be {'o','j','n',0}
	EncodeVersion float32 // usually 2.9
	Genre         uint32
	BPM           float32
	Level         [4]uint16 // [easy, normal, hard, (unused/SHD)]
	EventCount    [3]uint32
	NoteCount     [3]uint32
	MeasureCount  [3]uint32
	PackageCount  [3]uint32
	OldEncodeVer  uint16
	OldSongID     uint16
	OldGenreStr   [20]byte
	BMPSize       uint32 // 8x8 thumbnail right after cover
	FileVersion   uint32
	TitleRaw      [64]byte
	ArtistRaw     [32]byte
	NoterRaw      [32]byte
	OJMFileRaw    [32]byte
	CoverSize     uint32
	TimeSeconds   [3]uint32 // duration per difficulty
	NoteOffset    [3]uint32 // byte offsets to ex/nm/hd sections
	CoverOffset   uint32    // byte offset to cover (JPEG)
}

// Title/Artist/Noter/OJMFile return decoded UTF-8 strings.
func (h Header) Title() string   { return decodeCP949Trim(h.TitleRaw[:]) }
func (h Header) Artist() string  { return decodeCP949Trim(h.ArtistRaw[:]) }
func (h Header) Noter() string   { return decodeCP949Trim(h.NoterRaw[:]) }
func (h Header) OJMFile() string { return decodeCP949Trim(h.OJMFileRaw[:]) }

// Difficulty tags
const (
	DiffEasy   = 0
	DiffNormal = 1
	DiffHard   = 2
)

// PackageHeader precedes a run of events for a single (measure, channel).
type PackageHeader struct {
	Measure uint32
	Channel uint16
	Events  uint16
}

// Event is a discriminated union for per-channel event payloads.
type Event struct {
	// Position is the 0-based index inside the package. To get a normalized
	// position inside the measure, divide by PackageHeader.Events.
	Position int

	// One of the following is set depending on Channel.
	Fraction *float32   // ChanMeasureFraction: fraction of the measure actually used
	BPM      *float32   // ChanBPM: new BPM at this position
	Note     *NoteEvent // ChanLaneX or Autoplay: note/sample trigger
	Raw      []byte     // unknown channel fallback
}

// NoteEvent payload for lane/autoplay channels.
type NoteEvent struct {
	Value  uint16 // 0 = padding/no note; otherwise sample index (WAV/OGG separate)
	Volume uint8  // 0 = max; 1..15 = quieter
	Pan    uint8  // 1..7 left->center, 0 or 8 center, 9..15 center->right
	Type   uint8  // NoteNormal/NoteLongHead/NoteLongTail/NoteOggFlag
}

// MeasureChannel groups a package and its events.
type MeasureChannel struct {
	Header PackageHeader
	Events []Event
}

// Section represents one difficulty's note section.
type Section struct {
	Difficulty   int // 0/1/2 (E/N/H)
	Offset       uint32
	PackageCount uint32
	Packages     []MeasureChannel
}

// OJN aggregates a parsed file.
type OJN struct {
	Header Header
	Easy   Section
	Normal Section
	Hard   Section
	// Cover holds the raw JPEG (and optional BMP thumbnail if present).
	CoverJPEG []byte
	ThumbBMP  []byte // may be empty
}

// Parse reads an OJN from r. It does not keep r positioned; callers control it.
func Parse(r io.ReaderAt) (*OJN, error) {
	br := io.NewSectionReader(r, 0, 1<<62)
	var h Header
	if err := binary.Read(br, binary.LittleEndian, &h); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if h.Signature != ojnMagic {
		return nil, fmt.Errorf("invalid signature: %v", h.Signature)
	}

	obj := &OJN{Header: h}
	// Parse three difficulty sections using offsets and package counts.
	secs := []struct {
		idx int
		sec *Section
	}{
		{DiffEasy, &obj.Easy}, {DiffNormal, &obj.Normal}, {DiffHard, &obj.Hard},
	}
	for i, s := range secs {
		sec := s.sec
		sec.Difficulty = s.idx
		sec.Offset = h.NoteOffset[i]
		sec.PackageCount = h.PackageCount[i]
		pkgs, err := readPackages(br, int64(sec.Offset), sec.PackageCount)
		if err != nil {
			return nil, fmt.Errorf("read section %d: %w", i, err)
		}
		sec.Packages = pkgs
	}

	// Read cover JPEG and optional 8x8 BMP thumbnail following it.
	if h.CoverOffset != 0 && h.CoverSize > 0 {
		obj.CoverJPEG = make([]byte, h.CoverSize)
		if _, err := br.ReadAt(obj.CoverJPEG, int64(h.CoverOffset)); err != nil {
			return nil, fmt.Errorf("read cover jpeg: %w", err)
		}
		if h.BMPSize > 0 {
			obj.ThumbBMP = make([]byte, h.BMPSize)
			if _, err := br.ReadAt(obj.ThumbBMP, int64(h.CoverOffset)+int64(h.CoverSize)); err != nil {
				// Non-fatal; some charts omit or misreport it.
				obj.ThumbBMP = nil
			}
		}
	}
	return obj, nil
}

func readPackages(r io.ReaderAt, offset int64, count uint32) ([]MeasureChannel, error) {
	br := io.NewSectionReader(r, offset, 1<<62)
	pkgs := make([]MeasureChannel, 0, count)
	for i := uint32(0); i < count; i++ {
		var ph PackageHeader
		if err := binary.Read(br, binary.LittleEndian, &ph); err != nil {
			return nil, fmt.Errorf("pkg %d header: %w", i, err)
		}
		mc := MeasureChannel{Header: ph}
		mc.Events = make([]Event, 0, ph.Events)
		for e := 0; e < int(ph.Events); e++ {
			var ev Event
			ev.Position = e
			switch ph.Channel {
			case ChanMeasureFraction:
				var f float32
				if err := binary.Read(br, binary.LittleEndian, &f); err != nil {
					return nil, fmt.Errorf("pkg %d event %d: %w", i, e, err)
				}
				ev.Fraction = &f
			case ChanBPM:
				var f float32
				if err := binary.Read(br, binary.LittleEndian, &f); err != nil {
					return nil, fmt.Errorf("pkg %d event %d: %w", i, e, err)
				}
				ev.BPM = &f
			default:
				// Lane/autoplay or unknown: 4-byte payload
				buf := make([]byte, 4)
				if _, err := io.ReadFull(br, buf); err != nil {
					return nil, fmt.Errorf("pkg %d event %d: %w", i, e, err)
				}
				if ph.Channel >= ChanLane1 && ph.Channel <= 22 {
					ne := NoteEvent{
						Value:  binary.LittleEndian.Uint16(buf[0:2]),
						Volume: buf[2] & 0x0F,        // lower nibble 0=max, 1..15 quieter
						Pan:    (buf[2] >> 4) & 0x0F, // upper nibble (0/8=center)
						Type:   buf[3],
					}
					ev.Note = &ne
				} else {
					ev.Raw = buf
				}
			}
			mc.Events = append(mc.Events, ev)
		}
		pkgs = append(pkgs, mc)
	}
	return pkgs, nil
}

// Helpers

func decodeCP949Trim(b []byte) string {
	// Trim at first NUL
	if n := bytes.IndexByte(b, 0x00); n >= 0 {
		b = b[:n]
	}
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return ""
	}
	// Best-effort CP949 -> UTF-8; fall back to byte->string
	rd := transform.NewReader(bytes.NewReader(b), korean.EUCKR.NewDecoder())
	out, err := io.ReadAll(rd)
	if err != nil {
		return string(b)
	}
	return strings.TrimSpace(string(out))
}

// Convenience: ParseFromReader wraps an io.Reader by buffering to memory to
// satisfy io.ReaderAt. Prefer Parse on a ReaderAt for large files.
func ParseFromReader(r io.Reader) (*OJN, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Parse(bytes.NewReader(buf))
}

// Convenience: ParseFile reads an OJN from disk.
func ParseFile(path string) (*OJN, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(bytes.NewReader(f))
}

// ScanChannelEvents flattens packages into normalized events for a channel.
// For each package on the given channel, it yields (measure, t, event), where
// t is a 0..1 position within the measure.
type ChannelEvent struct {
	Measure uint32
	T       float64 // 0..1 within measure
	Event   Event
}

func (s *Section) ScanChannelEvents(channel uint16) []ChannelEvent {
	var out []ChannelEvent
	for _, p := range s.Packages {
		if p.Header.Channel != channel || p.Header.Events == 0 {
			continue
		}
		div := float64(p.Header.Events)
		for _, ev := range p.Events {
			out = append(out, ChannelEvent{Measure: p.Header.Measure, T: float64(ev.Position) / div, Event: ev})
		}
	}
	return out
}

// Validate performs a few basic sanity checks; returns nil if OK.
func (o *OJN) Validate() error {
	if o.Header.Signature != ojnMagic {
		return errors.New("invalid magic")
	}
	// Offsets must be within file; cannot check size from ReaderAt here.
	// Ensure package counts match actual data lengths per section.
	for _, sec := range []*Section{&o.Easy, &o.Normal, &o.Hard} {
		if uint32(len(sec.Packages)) != sec.PackageCount {
			return fmt.Errorf("section %d: package count mismatch header=%d actual=%d",
				sec.Difficulty, sec.PackageCount, len(sec.Packages))
		}
	}
	return nil
}

// Example usage (not built by default):
//
//  func main() {
//      f, _ := os.Open("o2ma1234.ojn")
//      defer f.Close()
//      ojn, err := o2jam.Parse(f)
//      if err != nil { log.Fatal(err) }
//      fmt.Println(ojn.Header.Title(), ojn.Header.Artist())
//      notes := ojn.Hard.ScanChannelEvents(ChanLane4)
//      fmt.Printf("Hard middle-lane events: %d\n", len(notes))
//  }
