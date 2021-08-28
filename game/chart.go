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
	BG Sprite

	Parameter map[string]float64
	TimingPoints

	Level float64 // todo: mods may change level
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
	//c.SetBG(bg.Filename)

	c.TimingPoints = newTimingPointsFromOsu(o)
	if dat, err := ioutil.ReadFile(c.AbsPath(c.AudioFilename)); err == nil {
		c.AudioHash = md5.Sum(dat)
	}
	c.Parameter["Scale"] = o.CircleSize
	return &c
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

// temp: ChartHeader가 Sprite를 가진다
// 한편 Sprite는 ScreenSize에 종속이다
// gob 등으로 정보를 재활용하고자 할 때에는 Sprite Reload 등의 작업이 필요할 것으로 예상
func (c *ChartHeader) SetBG(fname string) {
	path := c.AbsPath(fname)
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err) // log.Fatal은 에러가 난 위치를 알려주지 않는 듯
	}
	src, _, err := image.Decode(bytes.NewReader(dat))
	if err != nil {
		panic(err)
	}
	i, _ := ebiten.NewImageFromImage(src, ebiten.FilterDefault)

	sprite := NewSprite(i)

	sw := i.Bounds().Dx()
	sh := i.Bounds().Dy()
	screenX := Settings.ScreenSize.X
	screenY := Settings.ScreenSize.Y
	w, h := sw, sh
	if sw > screenX || sh > screenY { // 스크린이 그림보다 작을 경우 그림 크기 줄이기
		minRatio := screenX / sw
		if minRatio > screenY/sh {
			minRatio = screenY / sh
		}
		w *= minRatio
		h *= minRatio
	}

	x := screenX/2 - w/2
	y := screenY/2 - h/2
	// x, y := bx*ratio, by*ratio // x와 y 둘 중 하나는 스크린 크기와 일치는 보류
	sprite.SetFixedOp(w, h, x, y)
	c.BG = sprite
}
