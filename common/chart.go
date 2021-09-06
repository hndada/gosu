package common

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/engine/ui"
	"github.com/hndada/rg-parser/osugame/osu"
)

const (
	ModeStandard = iota
	ModeTaiko
	ModeCatch
	ModeMania
)

const defaultMode = ModeStandard

type ChartHeader struct {
	ChartPath     string
	ChartID       int64 // 6byte: setID, 2byte: subID
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string // diff name
	Producer      string
	HolderID      int64 // 0: gosu Chart Management

	AudioFilename   string
	AudioHash       [md5.Size]byte // for checking music data update
	PreviewTime     int64
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64

	Parameter map[string]float64
	Level     float64
}

// Sprite는 ScreenSize에 종속이다.
// gob 등으로 정보를 재활용하고자 할 때에는 Sprite Reload 등의 작업이 필요할 것으로 예상
// path is needed for lazy load: BG, Video
func NewChartHeaderFromOsu(o *osu.Format, path string) *ChartHeader {
	c := ChartHeader{
		ChartPath:     path,
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
	if dat, err := ioutil.ReadFile(c.Path(c.AudioFilename)); err == nil {
		c.AudioHash = md5.Sum(dat)
	}
	{
		e, ok := o.Background()
		if !ok {
			panic("failed to load bg")
		}
		c.ImageFilename = e.Filename
	}
	{
		e, ok := o.Video()
		if ok {
			c.VideoFilename = e.Filename
			c.VideoTimeOffset = int64(e.StartTime)
		}
	}
	switch o.General.Mode {
	case ModeStandard, ModeCatch:
		c.Parameter["CircleSize"] = o.CircleSize
	case ModeMania:
		c.Parameter["KeyCount"] = o.CircleSize
	}
	return &c
}

func (c ChartHeader) Path(fname string) string {
	d := filepath.Dir(c.ChartPath)
	return filepath.Join(d, fname)
}
func (c ChartHeader) BG(dimness float64) ui.FixedSprite {
	var src *ebiten.Image
	path := c.Path(c.ImageFilename) // chart's own background file path
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		src = Skin.DefaultBG
	} else {
		i, _, err := image.Decode(bytes.NewReader(dat))
		if err != nil {
			panic(err)
		}
		src = ebiten.NewImageFromImage(i)
	}
	sprite := ui.NewFixedSprite(src)
	sw := src.Bounds().Dx()
	sh := src.Bounds().Dy()
	screenX := Settings.ScreenSize.X
	screenY := Settings.ScreenSize.Y
	w, h := sw, sh
	ratioW, ratioH := float64(screenX)/float64(sw), float64(screenY)/float64(sh)
	minRatio := ratioW
	if minRatio > ratioH {
		minRatio = ratioH
	}
	// BG가 스크린보다 크든 작든 min ratio 곱해지면 딱 맞춰짐
	w = int(float64(w) * minRatio)
	h = int(float64(h) * minRatio)
	x := screenX/2 - w/2
	y := screenY/2 - h/2
	sprite.W = w
	sprite.H = h
	sprite.X = x
	sprite.Y = y
	sprite.Dimness = dimness
	sprite.Fix()
	return sprite
}

func (c ChartHeader) AudioPath() string {
	if c.AudioFilename == "virtual" { // keysound only
		return ""
	}
	return c.Path(c.AudioFilename)
}

func DefaultBG() ui.FixedSprite {
	const dimness = 1
	return ChartHeader{}.BG(dimness) // default background goes returned when error occurs
}

// Use when want to know the mode with no parsing whole .osu file
// If path's directing file isn't .osu, OsuMode panics.
func OsuMode(path string) int {
	if strings.ToLower(filepath.Ext(path)) != ".osu" {
		panic("not .osu file")
	}
	file, err := os.Open(path)
	if err != nil {
		panic(err)
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
