package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/inpututil"
	"log"
)

/*
음악 재생
vy = 일정

 */

const (
	screenWidth  = 1600
	screenHeight = 900
)

type musicType int
const (
	typeOgg musicType = iota
	typeMP3
)



type Game struct {
	musicPlayer   *Player
	musicPlayerCh chan *Player
	errCh         chan error
}

func NewGame() (*Game, error) {
	audioContext, err := audio.NewContext(sampleRate)
	if err != nil {
		return nil, err
	}

	m, err := NewPlayer(audioContext, typeOgg)
	if err != nil {
		return nil, err
	}

	return &Game{
		musicPlayer:   m,
		musicPlayerCh: make(chan *Player),
		errCh:         make(chan error),
	}, nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	runnableOnUnfocused := ebiten.IsRunnableOnUnfocused()
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		runnableOnUnfocused = !runnableOnUnfocused
	}
	ebiten.SetRunnableOnUnfocused(runnableOnUnfocused)
	select {
	case p := <-g.musicPlayerCh:
		g.musicPlayer = p
	case err := <-g.errCh:
		return err
	default:
	}

	if g.musicPlayer != nil {
		if err := g.musicPlayer.update(); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.musicPlayer != nil {
		g.musicPlayer.draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetMaxTPS(240)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	g, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
