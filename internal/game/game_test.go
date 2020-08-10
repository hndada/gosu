package game

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	g := NewGame()
	g.Scene = &Title{}
	// c := mania.NewChart(`./test/test_ln.osu`)
	// g.NextScene = mania.NewSceneMania(g.Options, c)

	// f, err := os.Open("./test/" + c.AudioFilename)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// streamer, format, err := mp3.Decode(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer streamer.Close()

	// speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// // done := make(chan bool)
	// speaker.Play(streamer)
	// speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// 	done <- true
	// })))

	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight) // fixed in prototype
	ebiten.SetWindowTitle("gosu")
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetMaxTPS(g.MaxTPS)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
	// <-done
}
