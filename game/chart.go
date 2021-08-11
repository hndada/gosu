package game

import (
	"bufio"
	"crypto/md5"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hndada/rg-parser/osugame/osu"
)

const Millisecond = 1000

const (
	ModeStandard = iota
	ModeTaiko
	ModeCatch
	ModeMania
)

// TransPoint를 Base에 두지 않는다면, ChartHeader로 바꾸어도 된다고 생각
type BaseChart struct {
	Path          string // path of chart source file
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

	Level float64 // 모드 별 레벨의 필요성
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
		Producer:      o.Creator, // field name may change

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

func OsuMode(path string) int {
	const defaultMode = ModeStandard
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	line, err := r.ReadString('\n')
	for err == nil {
		if strings.HasPrefix(line, "Mode: ") {
			s := strings.Split(line, ": ")
			if len(s) < 2 {
				return defaultMode
			}
			mode, err := strconv.ParseInt(string(s[1][0]), 10, 0)
			if err != nil {
				return defaultMode
			}
			return int(mode)
		}
		line, err = r.ReadString('\n')
	}
	return defaultMode
}
func (b *BaseChart) AbsPath(filename string) string {
	return filepath.Join(filepath.Dir(b.Path), filename)
}
