package ui

import (
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

// process flow: from image.Image to image.Image

// LoadImageImage loads from file path to image.Image
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
	name := strings.TrimRight(path, ext)
	name = strings.TrimRight(name, "@2x")
	path2x := name + "@2x" + ext
	path1x := name + ext

	i, err := LoadImageImage(path2x)
	if err != nil {
		return LoadImageImage(path1x)
	}
	return i, err
}
