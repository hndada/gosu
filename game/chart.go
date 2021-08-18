package game

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten"
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
type ChartHeader struct {
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

func NewChartHeaderFromOsu(o *osu.Format, path string) *ChartHeader {
	c := ChartHeader{
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
	c.ImageFilename = bg.Filename

	c.TimingPoints = newTimingPointsFromOsu(o)
	if dat, err := ioutil.ReadFile(c.AbsPath(c.AudioFilename)); err == nil {
		c.AudioHash = md5.Sum(dat)
	}
	c.Parameter["Scale"] = o.CircleSize
	return &c
}

func (c *ChartHeader) Background() (*ebiten.Image, error) {
	dat, err := ioutil.ReadFile(c.AbsPath(c.ImageFilename))
	if err != nil {
		return nil, err
	}
	src, _, err := image.Decode(bytes.NewReader(dat))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(src, ebiten.FilterDefault)
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

func (c ChartHeader) AbsPath(filename string) string {
	return filepath.Join(filepath.Dir(c.Path), filename)
}
