package main

import (
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

func main() {
	i, err := load("test.png")
	if err != nil {
		panic(err)
	}
	dc := gg.NewContextForImage(i)
	dc.InvertY()
	gg.SavePNG("test2.png", i)
}

func load(path string) (image.Image, error) {
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
	return i, err
}
