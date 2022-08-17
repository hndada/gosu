package mode

import (
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/format/osu"
)

// Now mode consists of main mode + sub mode.
// Key count and scratch mode are sub modes of Piano mode.
// const (
//
//	ModePiano = 1 << (iota + 7) // 128
//	ModeDrum                    // 256
//	ModeJjava                   // 512
//
// )
// const ModeDefault = ModePiano + 4
const (
	ModePiano4 = iota // 1 ~ 4 Key
	ModePiano7        // 5 ~ Key
	ModeDrum
	ModeKaraoke // aka jjava
)
const ModeDefault = ModePiano4
const ModeUnknown = -1

// ChartHeader contains non-play information.
// Chaning ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
// Music: when treating it as media
// Audio: when considering as programming aspect
type ChartHeader struct {
	ChartID       int64 // 6byte: setID, 2byte: subID
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string // diff name
	Producer      string // Name of field may change
	HolderID      int64  // 0: gosu Chart Management

	AudioFilename   string
	PreviewTime     int64
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64
}

func NewChartHeader(f any) ChartHeader {
	var c ChartHeader
	switch f := f.(type) {
	case *osu.Format:
		c = ChartHeader{
			MusicName:     f.Title,
			MusicUnicode:  f.TitleUnicode,
			Artist:        f.Artist,
			ArtistUnicode: f.ArtistUnicode,
			MusicSource:   f.Source,
			ChartName:     f.Version,
			Producer:      f.Creator,

			AudioFilename: f.AudioFilename,
			PreviewTime:   int64(f.PreviewTime),
		}
		var e osu.Event
		e, _ = f.Background()
		c.ImageFilename = e.Filename
		e, _ = f.Video()
		c.VideoFilename, c.VideoTimeOffset = e.Filename, int64(e.StartTime)
	}
	return c
}

func (c ChartHeader) MusicPath(cpath string) string {
	return filepath.Join(filepath.Dir(cpath), c.AudioFilename)
}
func (c ChartHeader) BackgroundPath(cpath string) string {
	return filepath.Join(filepath.Dir(cpath), c.ImageFilename)
}

// mode.Chart is a base chart.
type Chart struct {
	ChartHeader
	TransPoints []*TransPoint
	Mode        int
	Mode2       int // KeyCount, for example.
	Duration    int64
	NoteCounts  []int
}

// Mode determines a mode of chart file by its path.
// Todo: should I make a new type Mode?
func Mode(fpath string) int {
	switch strings.ToLower(filepath.Ext(fpath)) {
	case ".osu":
		mode, keyCount := osu.Mode(fpath)
		switch mode {
		case osu.ModeMania:
			if keyCount <= 4 {
				return ModePiano4
			}
			return ModePiano7
		case osu.ModeTaiko:
			return ModeDrum
		default:
			return ModeUnknown
		}
	case ".ojn", ".bms":
		return ModePiano7
	}
	return ModeUnknown
}
