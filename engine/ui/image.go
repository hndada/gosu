package ui

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// LoadImage loads from file to ebiten.Image
// The function internally loads LoadImageImage, which loads from file to image.Image
func LoadImage(path string) (*ebiten.Image, error) {
	i, err := LoadImageImage(path)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(i), nil
}

func LoadImageHD(path string) (*ebiten.Image, error) {
	i, err := LoadImageImageHD(path)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(i), nil
}

// LoadImageImage loads from file to image.Image
func LoadImageImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return i, nil
}

//LoadImageImageHD tries @2x load first
func LoadImageImageHD(path string) (image.Image, error) {
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(path, ext)
	name = strings.TrimSuffix(name, "@2x")
	path2x := name + "@2x" + ext
	path1x := name + ext
	// fmt.Println(name, path2x, path1x)
	i, err := LoadImageImage(path2x)
	if err != nil {
		return LoadImageImage(path1x)
	}
	return i, err
}
