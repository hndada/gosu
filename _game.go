package main

import (
	"errors"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hndada/gosu/game/parser"
	"path/filepath"
)


// (폐기) beatmap index*(box height) 로 y position range 계산
const (
	stageIntro = iota
	stageSelect
	stageMania
	stageResult
	// stageLobby
	// stageRoom
	// stageEdit
)

var (
	button *ebiten.Image
)

type Game struct {
	audioPlayer
	beatmaps []parser.Beatmap
	config
	stage        int
	selectCursor int // (window 포함)
	selectGroup  string
	skin
}

func (g *Game) handleMovement() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.selectGroup = filepath.Dir(g.selectGroup)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		// dirname:= [g.selectCursor]
		g.selectGroup=filepath.Join(g.selectGroup, dirname)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.selectCursor += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.selectCursor -= 1
	}
	// 마우스 오른쪽 버튼 누르면 해당 y포지션으로 날아가기

}

func (g *Game) Update(screen *ebiten.Image) error {
	if button ==nil {
		button, _ = ebiten.NewImage(300, 30, ebiten.FilterDefault)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("game ended by player")
	}
	switch g.stage {
	case stageSelect:
		g.handleMovement() // 매번 movement 체크
	}
	return nil
}

// beatmap idx
// 모든 box 생성 (얼마 안 걸릴듯)
// 커서는 고정 (투덱에서 확인)
// 마우스 오른쪽 버튼 누르면 해당하는 idx의 비트맵으로 커서 잡히기
func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1600, 900
}
