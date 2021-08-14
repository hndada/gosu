package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

// suppose all number images has same size
func main() {
	combos := make([]image.Image, 10)
	for i := range combos {
		fname := fmt.Sprintf("score-%d.png", i)
		b, err := ioutil.ReadFile(fname)
		if err != nil {
			log.Fatal(err)
		}
		r := bytes.NewReader(b)
		img, err := png.Decode(r)
		combos[i] = img
	}
	tileSize := combos[0].Bounds().Size()
	tx, ty := tileSize.X, tileSize.Y
	mergedImg := image.NewNRGBA(image.Rect(0, 0, 5*tx, 2*ty))
	for i, c := range combos {
		x, y := (i%5)*tx, (i/5)*ty
		x2, y2 := x+tx, y+ty
		draw.Draw(mergedImg, image.Rect(x, y, x2, y2), c, c.Bounds().Min, draw.Over)
	}
	f, err := os.Create("score.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, mergedImg)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}
}
