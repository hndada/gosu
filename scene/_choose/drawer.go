package choose

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/scene"
)

var grayCover draws.Sprite

func init() {
	i := draws.NewImage(ScreenSizeX, ScreenSizeY)
	i.Fill(color.NRGBA{128, 128, 128, 128})
	grayCover = draws.NewSpriteFromSource(i)
}

type SearchDrawer struct {
	draws.Timer
	draws.Sprite
	query *string
}

func NewSearchDrawer(query *string) SearchDrawer {
	const (
		x = ScreenSizeX - RowWidth
		y = 25
	)
	i := draws.NewImage(RowWidth, 50)
	i.Fill(color.NRGBA{153, 217, 234, 192})
	s := draws.NewSpriteFromSource(i)
	s.Locate(x, y, draws.LeftTop)
	return SearchDrawer{
		Timer:  draws.NewTimer(0, draws.ToTick(1000, TPS)),
		query:  query,
		Sprite: s,
	}
}
func (d *SearchDrawer) Update() {
	d.Ticker()
}
func (d SearchDrawer) Draw(dst draws.Image) {
	t := *d.query
	if t == "" {
		t = "Type for search..."
	}
	text.Draw(dst.Image, t, scene.Face16, int(d.X), int(d.Y)+25, color.White)
}

type LoadingDrawer struct {
	draws.Timer
}

func NewLoadingDrawer() LoadingDrawer {
	i := draws.NewImage(ScreenSizeX, ScreenSizeY)
	i.Fill(color.NRGBA{128, 128, 128, 128})
	return LoadingDrawer{
		Timer: draws.NewTimer(0, draws.ToTick(600, TPS)),
	}
}
func (d *LoadingDrawer) Update() {
	d.Ticker()
}
func (d LoadingDrawer) Draw(dst draws.Image) {
	const (
		x = ScreenSizeX/2 - 100
		y = ScreenSizeY/2 + 30
	)
	grayCover.Draw(dst, draws.Op{})
	t := "Loading"
	age := float64(d.Tick) / float64(d.Period) // Todo: generalize at draws package
	c := int(3*age + 1)
	t += strings.Repeat(".", c)
	text.Draw(dst.Image, t, scene.Face24, x, y, color.White)
}

type KeySettingsDrawer struct {
	mode int
	keys []input.Key
}

func (d KeySettingsDrawer) Draw(dst draws.Image) {
	const (
		x = ScreenSizeX/2 - 100
		y = ScreenSizeY/2 + 30
	)
	mname := []string{"Piano4", "Piano7", "Drum"}[d.mode]
	keys := strings.Join(input.KeysToNames(d.keys), " ")
	t := mname + ": " + keys
	text.Draw(dst.Image, t, scene.Face24, x, y, color.White)
}
