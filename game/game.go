package game

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/framework/draws"
	"github.com/hndada/gosu/framework/scene"
	"github.com/hndada/gosu/game/chart"
	"github.com/hndada/gosu/game/format/osr"
)

// ScreenSize is a logical size of in-game screen.
const (
	ScreenSizeX = 1600
	ScreenSizeY = 900
)

// Underscore is for avoiding getting same name with package game
// while letting it be unexported struct.
type _Game struct {
	scene.Scene
}
type NewScenePlay func(fsys fs.FS, cname string, mods interface{}, rf *osr.Format) (_scene scene.Scene, err error)

// Todo: should .zip be extracted throughly?
func ZipFS(name string) fs.FS {
	r, err := zip.OpenReader(name)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		// gen file then return FS
		_, err = io.CopyN(os.Stdout, rc, 68)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
	}
	return r
}
func NewGame(newScenePlays []NewScenePlay) *_Game {
	g := &_Game{}
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(ScreenSizeX, ScreenSizeY)
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(TPS)

	var err error
	g.Scene, err = newScenePlays[1](os.DirFS("asdf - 1223"), "asdf - 1223 (MuangMuangE) [Oni].osu", nil, nil)
	if err != nil {
		panic(err)
	}
	return g
}
func (g *_Game) Update() (err error) {
	args := g.Scene.Update()
	switch args := args.(type) {
	case error:
		return args
	}
	return
}
func (g *_Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(draws.Image{Image: screen})
}
func (g *_Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenSizeX, ScreenSizeY
}

func SetTitle(c chart.Header) {
	title := fmt.Sprintf("gosu | %s - %s [%s] (%s) ", c.Artist, c.MusicName, c.ChartName, c.Charter)
	ebiten.SetWindowTitle(title)
}