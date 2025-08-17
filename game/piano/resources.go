package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/draws"
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
	HitSound           []byte
	ComboImages        []draws.Image // 10
	ScoreImages        []draws.Image // 13
}

func loadFieldImage() draws.Image {
	img := draws.CreateImage(game.ScreenSizeX, game.ScreenSizeY)
	return img
}

func loadBarImage() draws.Image {
	img := draws.CreateImage(1, 1)
	img.Fill(color.White)
	return img
}

func loadHintImage(fsys fs.FS) draws.Image {
	fname := "piano/hint.png"
	return draws.NewImageFromFile(fsys, fname)
}

func loadNoteFramesList(fsys fs.FS) [4]draws.Frames {
	var framesList [4]draws.Frames
	for nk, nkn := range []string{"normal", "head", "tail", "body"} {
		name := fmt.Sprintf("piano/note/%s.png", nkn)
		framesList[nk] = draws.NewFramesFromFile(fsys, name)
	}
	return framesList
}

func loadKeyButtonImages(fsys fs.FS) [2]draws.Image {
	var imgs [2]draws.Image
	for i, name := range []string{"up", "down"} {
		fname := fmt.Sprintf("piano/key/%s.png", name)
		imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
	return imgs
}

func loadBacklightImage(fsys fs.FS) draws.Image {
	fname := "piano/light/back.png"
	return draws.NewImageFromFile(fsys, fname)
}

func loadHitLightFrames(fsys fs.FS) draws.Frames {
	fname := "piano/light/hit.png"
	return draws.NewFramesFromFile(fsys, fname)
}

func loadHoldLightFrames(fsys fs.FS) draws.Frames {
	fname := "piano/light/hold.png"
	return draws.NewFramesFromFile(fsys, fname)
}

func loadJudgmentFramesList(fsys fs.FS) [4]draws.Frames {
	var framesList [4]draws.Frames
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		fname := fmt.Sprintf("piano/judgment/%s.png", name)
		framesList[i] = draws.NewFramesFromFile(fsys, fname)
	}
	return framesList
}

func loadHitSound(fsys fs.FS) []byte {
	fname := "piano/hit.wav"
	data, _ := fs.ReadFile(fsys, fname)
	return data
}

func NewResources(fsys fs.FS) *Resources {
	return &Resources{
		FieldImage:         loadFieldImage(),
		BarImage:           loadBarImage(),
		HintImage:          loadHintImage(fsys),
		NotesFramesList:    loadNoteFramesList(fsys),
		KeyButtonsImages:   loadKeyButtonImages(fsys),
		BacklightsImage:    loadBacklightImage(fsys),
		HitLightsFrames:    loadHitLightFrames(fsys),
		HoldLightsFrames:   loadHoldLightFrames(fsys),
		JudgmentFramesList: loadJudgmentFramesList(fsys),
		HitSound:           loadHitSound(fsys),
		ComboImages:        game.LoadComboImages(fsys),
		ScoreImages:        game.LoadScoreImages(fsys),
	}
}
