package gosu

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/parse/osr"
	"github.com/hndada/gosu/parse/osu"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Todo: should be ChartHeader instead of Chart
// Todo: might integrate database here
type ChartInfo struct {
	Chart *Chart
	Path  string
	Level float64
	Box   Sprite
}

type SceneSelect struct {
	ChartInfos []ChartInfo
	Cursor     int
	Background Sprite
	Hold       int
	HoldKey    ebiten.Key

	ReplayMode    bool
	IndexToMD5Map map[int][md5.Size]byte
	MD5ToIndexMap map[[md5.Size]byte]int
	Replays       []*osr.Format

	PlaySoundMove   func()
	PlaySoundSelect func()
}

const (
	bw  = 450 // Box width
	bh  = 50  // Box height
	pop = bw / 10
)

func NewSceneSelect() *SceneSelect {
	s := &SceneSelect{
		ChartInfos:    make([]ChartInfo, 0, 50),
		IndexToMD5Map: make(map[int][16]byte),
		MD5ToIndexMap: map[[16]byte]int{},
		Replays:       make([]*osr.Format, 0, 10),
	}
	dirs, err := os.ReadDir(MusicPath)
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dpath := filepath.Join(MusicPath, dir.Name())
		fs, err := os.ReadDir(dpath)
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			if f.IsDir() { // There may be directory e.g., SB
				continue
			}
			fpath := filepath.Join(dpath, f.Name())
			b, err := os.ReadFile(fpath)
			if err != nil {
				panic(err)
			}
			switch strings.ToLower(filepath.Ext(fpath)) {
			case ".osu":
				o, err := osu.Parse(b)
				if err != nil {
					panic(err)
				}
				if o.Mode == ModeMania {
					c, err := NewChartFromOsu(o)
					if err != nil {
						panic(err)
					}
					info := ChartInfo{
						Chart: c,
						Path:  fpath,
						Level: c.Level(),
						Box: Sprite{
							I: NewBox(c, c.Level()),
							W: bw,
							H: bh,
						},
					}
					s.ChartInfos = append(s.ChartInfos, info)
					// box's x value is not fixed.
					// box's y value is not fixed.
				}
			}
		}
	}
	sort.Slice(s.ChartInfos, func(i, j int) bool {
		if s.ChartInfos[i].Chart.MusicName == s.ChartInfos[j].Chart.MusicName {
			return s.ChartInfos[i].Level < s.ChartInfos[j].Level
		}
		return s.ChartInfos[i].Chart.MusicName < s.ChartInfos[j].Chart.MusicName
	})
	for i, ci := range s.ChartInfos {
		d, err := os.ReadFile(ci.Path)
		if err != nil {
			panic(err)
		}
		s.MD5ToIndexMap[md5.Sum(d)] = i
		s.IndexToMD5Map[i] = md5.Sum(d)
	}

	fs, err := os.ReadDir("replay")
	if err != nil {
		panic(err)
	}
	for _, f := range fs {
		if f.IsDir() || filepath.Ext(f.Name()) != ".osr" {
			continue
		}
		rd, err := os.ReadFile(filepath.Join("replay", f.Name()))
		if err != nil {
			panic(err)
		}
		rf, err := osr.Parse(rd)
		if err != nil {
			panic(err)
		}
		s.Replays = append(s.Replays, rf)
	}

	s.UpdateBackground()
	s.HoldKey = HoldKeyNone
	s.Hold = threshold1
	_, apMove := NewAudioPlayer("skin/default-hover.wav")
	s.PlaySoundMove = apMove.PlaySoundEffect
	_, apSelect := NewAudioPlayer("skin/restart.wav")
	s.PlaySoundSelect = apSelect.PlaySoundEffect
	return s
}
func (s *SceneSelect) UpdateBackground() {
	s.Background = RandomDefaultBackground
	if len(s.ChartInfos) == 0 {
		return
	}
	info := s.ChartInfos[s.Cursor]
	img := NewImage(info.Chart.BackgroundPath(info.Path))
	if img != nil {
		s.Background.I = img
	}
}

const (
	border = 3
)
const (
	dx = 20 // dot x
	dy = 30 // dot y
)

var borderColor = color.RGBA{172, 49, 174, 255} // Purple

func NewBox(c *Chart, lv float64) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, bw, bh))
	draw.Draw(img, img.Bounds(), &image.Uniform{borderColor}, image.Point{}, draw.Src)
	inRect := image.Rect(border, border, bw-border, bh-border)
	draw.Draw(img, inRect, &image.Uniform{color.White}, image.Point{}, draw.Src)
	t := fmt.Sprintf("(%dK Lv %.1f) %s [%s]", c.KeyCount, lv, c.MusicName, c.ChartName)
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(dx * 64), Y: fixed.Int26_6(dy * 64)},
	}
	d.DrawString(t)
	return ebiten.NewImageFromImage(img)
}

const HoldKeyNone = -1

// Require holding for a while to move a cursor
var (
	threshold1 = MsecToTick(100)
	threshold2 = MsecToTick(80)
)

// FetchReplay returns first MD5-matching replay format.
// Todo: need to rewrite
func (s SceneSelect) FetchReplay(i int) *osr.Format {
	md5 := s.IndexToMD5Map[i]
outer:
	for _, r := range s.Replays {
		for x := 0; x < 16; x++ {
			ui, err := strconv.ParseUint(string(r.BeatmapMD5[x*2:(x+1)*2]), 16, 8)
			if err != nil {
				panic(err)
			}
			if ui != uint64(md5[x]) {
				continue outer
			}
		}
		return r
	}
	return nil
}

// Default HoldKey value is 0, which is Key0.
func (s *SceneSelect) Update(g *Game) {
	if s.HoldKey == HoldKeyNone {
		s.Hold++
		if s.Hold > threshold1 {
			s.Hold = threshold1
		}
	} else {
		if ebiten.IsKeyPressed(s.HoldKey) {
			s.Hold++
		} else {
			s.Hold = 0
		}
	}
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyEnter), ebiten.IsKeyPressed(ebiten.KeyNumpadEnter):
		s.PlaySoundSelect()
		info := s.ChartInfos[s.Cursor]
		if s.ReplayMode {
			g.Scene = NewScenePlay(info.Chart, info.Path, s.FetchReplay(s.Cursor), true)
		} else {
			g.Scene = NewScenePlay(info.Chart, info.Path, nil, true)
		}
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		s.HoldKey = ebiten.KeyArrowDown
		if s.Hold < threshold1 {
			break
		}
		s.PlaySoundMove()
		s.Hold = 0
		s.Cursor++
		s.Cursor %= len(s.ChartInfos)
		s.UpdateBackground()
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		s.HoldKey = ebiten.KeyArrowUp
		if s.Hold < threshold1 {
			break
		}
		s.PlaySoundMove()
		s.Hold = 0
		s.Cursor--
		if s.Cursor < 0 {
			s.Cursor += len(s.ChartInfos)
		}
		s.UpdateBackground()
	case ebiten.IsKeyPressed(ebiten.KeyQ):
		s.HoldKey = ebiten.KeyQ
		if s.Hold < threshold2 {
			break
		}
		s.Hold = 0
		BaseSpeed -= 0.1
		if BaseSpeed < 0.1 {
			BaseSpeed = 0.1
		}
	case ebiten.IsKeyPressed(ebiten.KeyW):
		s.HoldKey = ebiten.KeyW
		if s.Hold < threshold2 {
			break
		}
		s.Hold = 0
		BaseSpeed += 0.1
		if BaseSpeed > 2 {
			BaseSpeed = 2
		}
	case ebiten.IsKeyPressed(ebiten.KeyA):
		s.HoldKey = ebiten.KeyA
		if s.Hold < threshold2 {
			break
		}
		s.Hold = 0
		Volume -= 0.05
		if Volume < 0 {
			Volume = 0
		}
	case ebiten.IsKeyPressed(ebiten.KeyS):
		s.HoldKey = ebiten.KeyS
		if s.Hold < threshold2 {
			break
		}
		s.Hold = 0
		Volume += 0.05
		if Volume > 1 {
			Volume = 1
		}
	case ebiten.IsKeyPressed(ebiten.KeyZ):
		s.HoldKey = ebiten.KeyZ
		if s.Hold < threshold1 {
			break
		}
		s.Hold = 0
		s.ReplayMode = !s.ReplayMode
	default:
		s.HoldKey = HoldKeyNone
	}
}

// Currently topmost and bottommost boxes are not adjoined.
// May add extra effect to box arrangement.
// x -= y / 5, for example.
func (s SceneSelect) Draw(screen *ebiten.Image) {
	s.Background.Draw(screen)
	for i := range s.ChartInfos {
		y := (i-s.Cursor)*bh + screenSizeY/2 - bh/2
		if y > screenSizeY || y+bh < 0 {
			continue
		}
		x := screenSizeX - bw + pop
		if i == s.Cursor {
			x = screenSizeX - bw
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(s.ChartInfos[i].Box.I, op)
	}
	{
		sprite := GeneralSkin.CursorSprites[0]
		x, y := ebiten.CursorPosition()
		sprite.X, sprite.Y = float64(x), float64(y)
		sprite.Draw(screen)
	}
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("BaseSpeed (Press Q/W): %.0f\n(Exposure time: %.0fms)\n\nVolume (Press A/S): %d%%\nHold:%d\nReplay mode (Press Z): %v\n", // %.1f
			BaseSpeed*100, ExposureTime(BaseSpeed), int(Volume*100), s.Hold, s.ReplayMode))
}
