package mania

import (
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	_ "image/jpeg"
	"log"
	"os"
	"time"
)

// todo: game.go refactoring
// title -> 간이 song select -> play -> result -> title
// mp3 플레이어, Scene에 저장; 연동
// sync with mp3, position
// 곡선택: 맵정보패널

// 플레이: input (ebiten으로 간단히, 나중에 별도 라이브러리.)
// 점수계산: 1/n -> my score system
// 리플레이 실행 - 스코어/hp 시뮬레이터

// scene: abstract; contents in the screen
// screen: mere image data after all; screen is the result

// 모든 scene에 sceneManager가 하는 일을 embed하면 없어도 되지 않을까?

// Chart struct는 계산한 값 살리기 등을 위해서도 gob로 저장

// 업데이트 이후 Draw라고는 말 못함. 둘은 그냥 1/60마다 한번씩 실행됨

func main() {
	g := game.NewGame()
	c := NewChart(`./test/test_ln.osu`)
	g.Scene = NewSceneMania(g.Options, c)

	f, err := os.Open("./test/" + c.AudioFilename)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// done := make(chan bool)
	speaker.Play(streamer)
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
