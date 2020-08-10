package o2jam

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/text/encoding/korean"
	"io"
	"io/ioutil"
)

// .ojn, bpm이랑 다른 notes strip 핸들
// .ojm
// osu, ojn -> .gos chart loader (with interface)
// mania(얘도 이름 변경) 비트맵.go->chart.go 및 수정
// play 및 플레이 화면 만들기

// type Converter interface {
// 	Load() { } // .ojn -> .gos
// }

var euckrDec = korean.EUCKR.NewDecoder()

var genre = []string{"Ballad", "Rock", "Dance", "Techno", "Hip-hop",
	"Soul/R&B", "Jazz", "Funk", "Classical", "Traditional", "Etc"}

type OJN struct {
	Header
	Charts    [3]Chart
	Cover     []byte
	Thumbnail []byte
}

type Header struct {
	SongID           int32
	Signature        [4]byte // "OJN\0"
	EncodeVersion    float32 // Encoder value (9A 99 39 40)
	Genre            int32
	BPM              float32
	Level            [4]int16 // Last 2 bytes are unused
	EventCount       [3]int32 // including samples, BPM and measure event
	NoteCount        [3]int32 // playing notes only
	MeasureCount     [3]int32 // bar count
	StripCount       [3]int32
	OldEncodeVersion int16 // 29
	OldSongID        int16
	OldGenre         [20]byte
	BMPSize          int32
	OldFileVersion   int32
	Title            [64]byte
	Artist           [32]byte
	Noter            [32]byte
	OJMFile          [32]byte
	CoverSize        int32
	Time             [3]int32
	NoteOffset       [3]int32
	CoverOffset      int32
	_                int32
}

type Chart struct {
	Level int16
	Time  int32
	// MeasureCount int32

	MeasureFractions []MeasureFractionEvent
	BPMEvents        []Strip // BPMEvent
	Notes            []Strip // NoteEvent
	Samples          []Strip // NoteEvent
}

// Event: 박자, BPM, 노트, 샘플 모두
// 4바이트짜리 Event가 각 채널에 있는 스트립에 박히다
type Strip struct {
	StripHeader
	Events []uint32
}
type StripHeader struct {
	Measure    int32 // bar
	Channel    int16
	EventCount int16
}
type MeasureFractionEvent float32
type BPMEvent float32
type NoteEvent struct {
	Value     int16 // ignored when 0, otherwise ojm sample number
	VolumePan uint8 // 4bit each. Volume: 1~15=max, 0 is also max
	// Pan: 1~7: left, 0,8: center, 9~15: right
	NoteType int8 // 0: normal 2: ln start 3: ln end 4: ogg sample
}

// raw struct는 Parse 안에서 선언하고 깔끔한 struct를 또 만들까?
// 그런것보다는 그냥 euc-kr string정도만 함수 만들어주고 chart로 load할때 다루는게.
// byte 초과 길이로 저장하려 하면 31바이트 이후는 잘림
// Parse는 load에선 단순 작업만 할 정도로 전처리를 거의 다 끝내줘야함
func Parse(path string) (*OJN, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var ojn OJN
	r := bytes.NewReader(b)
	err = binary.Read(r, binary.LittleEndian, &ojn.Header)
	if err != nil {
		return &ojn, err
	}

	for diff := 0; diff < 3; diff++ {
		strips := make([]Strip, ojn.StripCount[diff])
		r.Seek(int64(ojn.NoteOffset[diff]), io.SeekStart)
		for i := 0; i < int(ojn.StripCount[diff]); i++ {
			var s Strip
			err = binary.Read(r, binary.LittleEndian, &s.StripHeader)
			if err != nil {
				return &ojn, err
			}
			s.Events = make([]uint32, s.EventCount)
			err = binary.Read(r, binary.LittleEndian, &s.Events)
			if err != nil {
				return &ojn, err
			}
			strips[i] = s
		}

		// measureFraction: measure당 1개, 해당 measure에 대해서만 적용. fraction 너머의 note는 찍을 수는 있는데 스킵됨
		// 동일 measure 안에서도 strip grid 수준이 다를 수 있다

		// bms도 내 생각엔 비슷한 구조일 거 같으니 ojn->bms->osu 가 될것 같음
		// NoteCount도 편집될 수 있으니 이걸로는 뭘 할수가 없음
		var c Chart
		c.Level = ojn.Level[diff]
		c.Time = ojn.Time[diff]

		c.MeasureFractions = make([]MeasureFractionEvent, ojn.MeasureCount[diff]+1)
		c.BPMEvents = make([]Strip, ojn.MeasureCount[diff]+1)
		c.Notes = make([]Strip, 0, 7*(ojn.MeasureCount[diff]+1))
		c.Samples = make([]Strip, 0, 32*(ojn.MeasureCount[diff]+1))
		// c.BPMEvents = make([]BPMEvent, 0, ojn.EventCount[diff]-ojn.NoteCount[diff])
		// c.Notes = make([]NoteEvent, 0, ojn.NoteCount[diff])
		// c.Samples = make([]NoteEvent, 0, ojn.EventCount[diff]-ojn.NoteCount[diff])

		for _, s := range strips {
			switch s.Channel {
			case 0: // MeasureFraction
				c.MeasureFractions[s.Measure] = MeasureFractionEvent(s.Events[0])
				if len(s.Events) != 1 {
					panic("unexpected measure fraction count")
				}
			case 1: // BPMChange
				if c.BPMEvents[s.Measure].StripHeader != (StripHeader{}) {
					panic("more than 2 strips in bpm channel")
				}
				c.BPMEvents[s.Measure] = s
			case 2, 3, 4, 5, 6, 7, 8: // Lane
				// c.Notes[s.Measure] = s
				c.Notes = append(c.Notes, s)
			default: // Sample
				c.Samples = append(c.Samples, s)
				// c.Samples[s.Measure] = s

			}
		}
		ojn.Charts[diff] = c
	}
	r.Seek(int64(ojn.CoverOffset), io.SeekStart)
	ojn.Cover = make([]byte, ojn.CoverSize)
	binary.Read(r, binary.LittleEndian, &ojn.Cover)
	ojn.Thumbnail = make([]byte, ojn.BMPSize)
	binary.Read(r, binary.LittleEndian, &ojn.Thumbnail)
	return &ojn, nil
}

func convEUCKR(b []byte) ([]byte, error) {
	s := make([]byte, len(b))
	r := euckrDec.Reader(bytes.NewReader(b))
	_, err := r.Read(s) // size
	if err != nil {
		return s, err
	}
	return bytes.TrimRight(s, string([]byte{0x00})), nil
}
