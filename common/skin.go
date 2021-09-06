package common

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/engine/ui"
)

// skin -> spritesheet
// 마지막으로 불러온 스킨 불러오기: 처음 / 오류 발생 시 defaultSkin
const (
	ScoreComma = iota + 10
	ScoreDot
	ScorePercent
)

// image.Image가 아닌 *ebiten.Image로 해야 이미지 자체가 한 번만 로드 됨
var Skin struct {
	Number1    [13]*ebiten.Image // including dot, comma, percent
	Number2    [13]*ebiten.Image
	HPBar      *ebiten.Image
	HPBarColor *ebiten.Image // todo: Animation
	BoxLeft    *ebiten.Image
	BoxRight   *ebiten.Image
	BoxMiddle  *ebiten.Image
	// Cursor      *ebiten.Image
	// CursorSmoke *ebiten.Image
	DefaultBG *ebiten.Image
}

func LoadImage(path string) (*ebiten.Image, error) {
	// temp: @2x 빠르게 적용
	var hdPath string
	if !strings.Contains(path, "@2x.") {
		switch filepath.Ext(path) {
		case ".png":
			hdPath = strings.Replace(path, ".png", "@2x.png", 1)
		case ".jpg":
			hdPath = strings.Replace(path, ".jpg", "@2x.jpg", 1)
		case ".jpeg":
			hdPath = strings.Replace(path, ".jpeg", "@2x.jpeg", 1)
		}
	}
	f, err := os.Open(hdPath)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Open(path)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	ei := ebiten.NewImageFromImage(i)
	return ei, nil
}

func LoadSkin(cwd string) {
	var err error
	dir := filepath.Join(cwd, "skin")
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("score-%d.png", i)
		path := filepath.Join(dir, name)
		Skin.Number1[i], err = ui.LoadImageHD(path)
		if err != nil {
			panic(err)
		}
	}
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("combo-%d.png", i)
		path := filepath.Join(dir, name)
		Skin.Number2[i], err = ui.LoadImageHD(path)
		if err != nil {
			panic(err)
		}
	}

	var path string
	path = filepath.Join(dir, "scorebar-bg.png")
	Skin.HPBar, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "scorebar-colour.png")
	Skin.HPBarColor, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "button-left.png")
	Skin.BoxLeft, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "button-middle.png")
	Skin.BoxMiddle, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "button-right.png")
	Skin.BoxRight, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "menu-background.jpg")
	Skin.DefaultBG, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
}
