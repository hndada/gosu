package mode

import (
	"crypto/md5"
	"github.com/hndada/rg-parser/osugame/osu"
	"io/ioutil"
	"path/filepath"
)

type BaseChart struct {
	Path          string // path of chart source file. It won't be exported to file content.
	ChartID       int64  // 6byte: setID, 2byte: subID
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
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
	TimingPoints
}

// todo: tidy pointer up
func NewBaseChart(path string) (*BaseChart, error) {
	var b = &BaseChart{}
	switch filepath.Ext(path) {
	case ".osu":
		o, err := osu.Parse(path)
		if err != nil {
			return b, err
		}
		*b = *NewBaseChartFromOsu(o, path)
	}
	return b, nil
}

func NewBaseChartFromOsu(o *osu.Format, path string) *BaseChart {
	b := BaseChart{
		Path:          path,
		MusicName:     o.Title,
		MusicUnicode:  o.TitleUnicode,
		Artist:        o.Artist,
		ArtistUnicode: o.ArtistUnicode,
		MusicSource:   o.Source,
		ChartName:     o.Version,
		Producer:      o.Creator, // 변경될 수 있음

		AudioFilename: o.AudioFilename,
		PreviewTime:   int64(o.PreviewTime),
		Parameter:     make(map[string]float64),
	}
	bg, ok := o.Background()
	if !ok {
		panic("failed to load bg")
	}
	b.ImageFilename = bg.Filename
	b.TimingPoints = newTimingPointsFromOsu(o)
	if dat, err := ioutil.ReadFile(b.AbsPath(b.AudioFilename)); err == nil {
		b.AudioHash = md5.Sum(dat)
	}
	b.Parameter["Scale"] = o.CircleSize
	return &b
}

// interface로 만들만한 게 없음
// TransPoint를 Base에 두지 않는다면, ChartHeader로 바꾸어도 된다고 생각

func (b *BaseChart) AbsPath(filename string) string {
	return filepath.Join(filepath.Dir(b.Path), filename)
}
