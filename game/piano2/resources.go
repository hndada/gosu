package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	draws "github.com/hndada/gosu/draws5"
	"github.com/hndada/gosu/game"
)

// Resources is a collection of images and sounds.
// It is not expected to be loaded multiple times.

// The word "Resource" is countable.
// Todo: make the fields unexported?
type Resources struct {
	FieldImage         draws.Image
	BarImage           draws.Image // generated
	HintImage          draws.Image
	NotesFramesList    [4]draws.Frames
	KeyButtonsImages   [2]draws.Image // up, down
	BacklightsImage    draws.Image
	HitLightsFrames    draws.Frames
	HoldLightsFrames   draws.Frames
	JudgmentFramesList [4]draws.Frames
	ComboImages        []draws.Image // 10
	ScoreImages        []draws.Image // 13
	HitSound           []byte
}

func newFieldImage() draws.Image {
	img := draws.CreateImage(game.ScreenSizeX, game.ScreenSizeY)
	return img
}

func newBarImage() draws.Image {
	img := draws.CreateImage(1, 1)
	img.Fill(color.White)
	return img
}

func newHintImage(fsys fs.FS) draws.Image {
	fname := "piano/hint.png"
	return draws.NewImageFromFile(fsys, fname)
}

func newNoteFramesList(fsys fs.FS) [4]draws.Frames {
	var framesList [4]draws.Frames
	for nk, nkn := range []string{"normal", "head", "tail", "body"} {
		name := fmt.Sprintf("piano/note/%s.png", nkn)
		framesList[nk] = draws.NewFramesFromFile(fsys, name)
	}
	return framesList
}

func newKeyButtonImages(fsys fs.FS) [2]draws.Image {
	var imgs [2]draws.Image
	for i, name := range []string{"up", "down"} {
		fname := fmt.Sprintf("piano/key/%s.png", name)
		imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
	return imgs
}

func newBacklightImage(fsys fs.FS) draws.Image {
	fname := "piano/light/back.png"
	return draws.NewImageFromFile(fsys, fname)
}

func newHitLightFrames(fsys fs.FS) draws.Frames {
	fname := "piano/light/hit.png"
	return draws.NewFramesFromFile(fsys, fname)
}

func newHoldLightFrames(fsys fs.FS) draws.Frames {
	fname := "piano/light/hold.png"
	return draws.NewFramesFromFile(fsys, fname)
}

func newJudgmentFramesList(fsys fs.FS) [4]draws.Frames {
	var framesList [4]draws.Frames
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		fname := fmt.Sprintf("piano/judgment/%s.png", name)
		framesList[i] = draws.NewFramesFromFile(fsys, fname)
	}
	return framesList
}

// game/combo.go
func newComboImages(fsys fs.FS) []draws.Image {
	imgs := make([]draws.Image, 10)
	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("combo/%d.png", i)
		imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
	return imgs
}

// game/score.go
func newScoreImages(fsys fs.FS) []draws.Image {
	imgs := make([]draws.Image, 13)
	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("score/%d.png", i)
		imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		fname := fmt.Sprintf("score/%s.png", name)
		imgs[i+10] = draws.NewImageFromFile(fsys, fname)
	}
	return imgs
}

func newHitSound(fsys fs.FS) []byte {
	fname := "piano/hit.wav"
	data, _ := fs.ReadFile(fsys, fname)
	return data
}

func NewResources(fsys fs.FS) *Resources {
	return &Resources{
		FieldImage:         newFieldImage(),
		BarImage:           newBarImage(),
		HintImage:          newHintImage(fsys),
		NotesFramesList:    newNoteFramesList(fsys),
		KeyButtonsImages:   newKeyButtonImages(fsys),
		BacklightsImage:    newBacklightImage(fsys),
		HitLightsFrames:    newHitLightFrames(fsys),
		HoldLightsFrames:   newHoldLightFrames(fsys),
		JudgmentFramesList: newJudgmentFramesList(fsys),
		ComboImages:        newComboImages(fsys),
		ScoreImages:        newScoreImages(fsys),
	}
}
