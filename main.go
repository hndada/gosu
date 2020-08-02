package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
)

var (
	mFont font.Face
)

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const (
		mFontSize = 12
		dpi       = 300
	)
	mFont = truetype.NewFace(tt, &truetype.Options{
		Size:    mFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

const (
	screenWidth  = 1600
	screenHeigth = 900
	// sampleRate = 44100
)

var (
	introImage *ebiten.Image
)

func init() {
	introImage, _ = ebiten.NewImage(screenWidth, screenHeigth, ebiten.FilterDefault)
	const (
		buttonWidth  = 250
		buttonHeigth = 100
	)
	bigButtons := []string{"gosu!"}
	for i, b := range bigButtons {
		width := buttonWidth * 2
		height := buttonHeigth * 3
		x := i*width + 300
		y := i*height + 300
		ebitenutil.DrawRect(introImage, float64(x), float64(y), float64(width), float64(height), color.White)
		text.Draw(introImage, b, mFont, x+width/4, y+height/2, color.Black)
	}
	smallButtons := []string{"(Multi)", "(Edit)", "Options"}
	for i, b := range smallButtons {
		width := buttonWidth
		height := buttonHeigth
		x := 300 + buttonWidth*2
		y := i*height + 300
		ebitenutil.DrawRect(introImage, float64(x), float64(y), float64(width), float64(height), color.White)
		text.Draw(introImage, b, mFont, x+width/4, y+height/2, color.Black)
	}
}

type Game struct {
}

func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

type Button struct {
	width int
	height int
}

func (b *Button) Draw(screen *ebiten.Image) {
	op:=&ebiten.DrawImageOptions{}

}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	screen.DrawImage(introImage, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeigth
}

// var (
// 	game = &game.Game{}
// )
//
// func init() {
// 	game.beatmaps = make([]beatmap.Beatmap, 0)
// 	songs, err := tools.LoadSongList("test_beatmap", ".osu")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, song := range songs {
// 		b, err := beatmap.ParseBeatmap(song)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		game.beatmaps = append(game.beatmaps, b)
// 	}
// }

func main() {
	ebiten.SetMaxTPS(240)
	ebiten.SetWindowSize(screenWidth, screenHeigth)
	ebiten.SetWindowTitle("gosu")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
