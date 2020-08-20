package mode

import (
	"bytes"
	"crypto/md5"
	"github.com/hndada/rg-parser/osugame/osu"
	"image"
	"io/ioutil"
	"path/filepath"
)

type BaseChart struct {
	Path          string // path of chart source file. It won't be exported to file content.
	ChartID       int64  // 6byte: setID, 2byte: subID
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
	// TimingPoint
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

func (b *BaseChart) Background() (image.Image, error) {
	dat, err := ioutil.ReadFile(b.AbsPath(b.ImageFilename))
	if err != nil {
		return nil, err
	}
	bg, _, err := image.Decode(bytes.NewReader(dat))
	if err != nil {
		return nil, err
	}
	return bg, nil
}
