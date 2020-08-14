package game

import (
	"github.com/hndada/gosu/mode/mania"
	"image"
	"testing"
)

func TestRenderPanel(t *testing.T) {
	c := mania.NewChart(`C:\Users\hndada\Documents\GitHub\hndada\gosu\mode\mania\test\test_ln.osu`)
	cp := NewChartPanel(c, image.Pt(200, 200))
	cp.Render()

	// // w, h := cp.Image.Size()
	// // img := cp.Image.SubImage(image.Rect(0, 0, w, h))
	// // fmt.Println(img.At(10, 10))
	//
	// f, err := os.Create("panel.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// f.Close()
	// err = png.Encode(f, cp.Image)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
