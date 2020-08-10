package mode

import (
	"crypto/md5"
	"github.com/hndada/gosu/parser/osu"
)

// interface 로 할까?

// human-readable
// 실제 제작자와 관리자가 다를 수 있다
// 예시: Muang이 Genesis 의 오투잼 채보를 올렸다
// producer: Genesis
// holder: Muang

// hash for replay: 채보 raw 데이터만
type BaseChart struct {
	ChartID       int64 // 6byte: setID, 2byte: subID
	Title         string
	TitleUnicode  string
	Artist        string
	ArtistUnicode string
	Source        string
	ChartName     string // version; diff name
	Producer      string
	HolderID      int64 // 0: Gosu Chart Management

	AudioFilename string
	AudioLeadIn   int64
	AudioHash     [md5.Size]byte // for checking music data update
	PreviewTime   int64
	ImageFilename string
	VideoFilename string
	VideoOffset   int64

	Parameter map[string]float64
}

// editor 관련 데이터는 전용 포맷으로. 
// 배포도 editing 겸용/제외 옵션 추가
// tags는 text file 로 관리

// hash, 최소한의 변화에만 반응하게 하고 싶다
// A라는 쉬운 곡이 B라는 어려운 곡을 사칭하는 것을 가정하고 hash target정하기

func NewBaseChart(o *osu.OSU) BaseChart {
	c := BaseChart{
		Title:         o.Metadata["Title"].(string),
		TitleUnicode:  o.Metadata["TitleUnicode"].(string),
		Artist:        o.Metadata["Artist"].(string),
		ArtistUnicode: o.Metadata["ArtistUnicode"].(string),
		Source:        o.Metadata["Source"].(string),
		ChartName:     o.Metadata["Version"].(string),
		Producer:      o.Metadata["Creator"].(string), // 변경될 수 있음

		AudioFilename: o.General["AudioFilename"].(string),
		AudioLeadIn:   int64(o.General["AudioLeadIn"].(int)),
		PreviewTime:   int64(o.General["PreviewTime"].(int)),
		ImageFilename: o.Image.Filename,
		VideoFilename: o.Video.Filename,
		VideoOffset:   o.Video.StartTime,
	}
	// AudioHash
	// TimingPoint
	c.Parameter = make(map[string]float64)
	c.Parameter["Scale"] = o.Difficulty["CircleSize"].(float64)
	return c
}
