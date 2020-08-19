package mode

import (
	"crypto/md5"
	"github.com/hndada/rg-parser/osugame/osu"
)

type BaseChart struct {
	ChartID       int64 // 6byte: setID, 2byte: subID
	SongName      string
	SongUnicode   string
	Artist        string
	ArtistUnicode string
	SongSource    string
	ChartName     string // diff name
	Producer      string
	HolderID      int64 // 0: gosu Chart Management

	AudioFilename string
	AudioHash     [md5.Size]byte // for checking music data update
	PreviewTime   int64
	ImageFilename string
	// VideoFilename string
	// VideoOffset   int64

	Parameter map[string]float64
}

func NewBaseChartFromOsu(o *osu.Format) BaseChart {
	b := BaseChart{
		SongName:      o.Title,
		SongUnicode:   o.TitleUnicode,
		Artist:        o.Artist,
		ArtistUnicode: o.ArtistUnicode,
		SongSource:    o.Source,
		ChartName:     o.Version,
		Producer:      o.Creator, // 변경될 수 있음

		AudioFilename: o.AudioFilename,
		PreviewTime:   int64(o.PreviewTime),
		// ImageFilename: o.Background().Filename,
		Parameter: make(map[string]float64),
	}
	// AudioHash
	// TimingPoint
	b.Parameter["Scale"] = o.CircleSize
	return b
}

// interface로 만들만한 게 없음