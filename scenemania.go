package gosu

import (
	"bytes"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/mode/mania"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// 플레이 하면서 리플레이 데이터 저장
// input: ebiten으로 간단히, 나중에 별도 라이브러리.
// score, hp
// todo: 노트 이미지 미리 렌더 후 그리기만 하기
// sync with mp3, position
type SceneMania struct { // aka Clavier
	Notes       []NoteImageInfo
	C           mania.Chart
	Tick        func() float64
	ScrollSpeed float64
	Time        float64
	Score       float64
	HP          float64

	Cover        *ebiten.Image
	Streamer     beep.StreamSeekCloser
	StreamFormat beep.Format
}

type NoteImageInfo struct {
	x, y, w, h float64
	clr        color.RGBA
}

// mania.NewChart, NewSceneMania 둘 다 오래 걸릴듯
// loading하는 동안 SceneMania 말고 game 쪽에서 block
// todo: 음악 켜지는 순간에 렉걸림 -> 스트리머에 1500ms 공백 파일을 그냥 넣자
// todo: scene_play 위에 각 모드 올리기?
func (s *SceneMania) Update(g *Game) error {
	// if s.Time > 0 && s.Streamer.Position() == 0 {
	// 	speaker.Init(s.StreamFormat.SampleRate, s.StreamFormat.SampleRate.N(time.Second/10))
	// 	nanoTime := time.Duration(int64(s.Time * 1e6))
	// 	s.Streamer.Seek(s.StreamFormat.SampleRate.N(nanoTime))
	// 	speaker.Play(s.Streamer)
	// }

	// if s.Streamer.Position() == 0 {
	// 	go func() {
	// 		time.Sleep(time.Millisecond * time.Duration(s.BufferTime()))
	// 		speaker.Play(s.Streamer)
	// 	}()
	// 	// bufferTimer := time.NewTimer(time.Millisecond * 1500)
	// 	// go func() {
	// 	// 	<-bufferTimer.C
	// 	// 	speaker.Play(s.Streamer)
	// 	// }()
	// }
	if s.Streamer.Position() == 0 {
		err := speaker.Init(s.StreamFormat.SampleRate, s.StreamFormat.SampleRate.N(time.Second/10))
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			n := s.StreamFormat.SampleRate.N(time.Millisecond * time.Duration(s.BufferTime()))
			speaker.Play(beep.Seq(beep.Silence(n), s.Streamer))
		}()
	}

	const endTime = 5.0 * 1000 // float64(s.C.Notes[len(s.C.Notes)-1].Time)
	s.Time += s.Tick()
	for i := range s.Notes {
		s.Notes[i].y += s.Tick() * s.ScrollSpeed
	}
	if s.Time > endTime {
		s.Streamer.Close()
		ebiten.SetWindowTitle("gosu")
		g.NextScene = &SceneResult{}
		g.TransCountdown = g.MaxTransCountDown()
	}
	return nil
}

// todo:범위 넘어간 애들은 Rect 안그리기 -> 오히려 fps 불안정
// todo: 비트맵 로딩 timeout 15초
// field를 미리 전부 그려놔야 할까?
// view(er)
// buffered channel은 block되지 않는다
// Loading 이라는 별도의 Lock을 둔 이상, 특별히 채널은 필요없는거 아닌가?

// game이 Scene을 가지고 있고 scene이 다시 game을 수정하는 게 뭔가 이상한 듯
func (s *SceneMania) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	const ratio = float64(1600) / 1920
	op.GeoM.Scale(ratio, ratio)
	op.ColorM.ChangeHSV(0, 1, 0.30)
	screen.DrawImage(s.Cover, op)

	ebitenutil.DrawRect(screen, 565, 0, 70*7, 900, color.RGBA{0, 0, 0, 180})
	ebitenutil.DrawRect(screen, 565, 730, 70*7, 10, color.RGBA{252, 106, 111, 255})
	for i := range s.Notes {
		ebitenutil.DrawRect(screen, s.Notes[i].x, s.Notes[i].y, s.Notes[i].w, s.Notes[i].h, s.Notes[i].clr)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTime: %.1fs", ebiten.CurrentFPS(), s.Time/1000))
}

// 노트, 그냥 네모 그리고 색깔 채워넣기
// &SceneMania{}로 하고 chart 로딩을 할까
// todo: bufferTime 1500ms 정도로 설정하면 fps 5나옴
// todo: mp3, 버퍼에 올리자. 오래 걸려도 되니까
// todo: examples/camera 참고하여 플레이 프레임 안정화
// todo: 곡목록 select 본격화
func (s *SceneMania) BufferTime() int64 { return 0 }
func NewSceneMania(g *Game, c *mania.Chart) (s *SceneMania) {
	g.Loading = true
	s = &SceneMania{}
	c = mania.NewChart(`C:\Users\hndada\Documents\GitHub\hndada\gosu\mode\mania\test\test_ln.osu`)
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", c.Title, c.ChartName)) // todo: can I change window title ?
	const w = 70
	const noteHeight = 25
	// const bufferTime = 1500
	s.C = *c
	s.Notes = make([]NoteImageInfo, len(c.Notes))
	s.ScrollSpeed = g.ScrollSpeed
	s.Tick = func() float64 { return 1000 / float64(g.MaxTPS) } // 스피드값 1 기준 초당 1000픽셀 내려오게 해야함
	s.Time = -float64(s.BufferTime())
	for i, n := range c.Notes {
		var y, h float64
		x := float64(n.Key*w + 565)
		switch n.Type {
		case mania.TypeNote:
			y = -float64(n.Time+s.BufferTime())*s.ScrollSpeed + 730
			h = noteHeight * s.ScrollSpeed
		case mania.LNHead:
			y = -float64(n.Time2+s.BufferTime())*s.ScrollSpeed + 730
			h = float64(n.Time2-n.Time+noteHeight) * s.ScrollSpeed
		}
		s.Notes[i] = NoteImageInfo{x, y, w, h, noteColor(n, c.Keys)}
	}
	b, err := ioutil.ReadFile("C:\\Users\\hndada\\Documents\\GitHub\\hndada\\gosu\\mode\\mania\\test\\" + c.ImageFilename)
	if err != nil {
		panic(err)
	}
	cover, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	s.Cover, _ = ebiten.NewImageFromImage(cover, ebiten.FilterDefault)
	g.Loading = false

	f, err := os.Open("C:\\Users\\hndada\\Documents\\GitHub\\hndada\\gosu\\mode\\mania\\test\\" + c.AudioFilename)
	if err != nil {
		log.Fatal(err)
	}
	s.Streamer, s.StreamFormat, err = mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// defer streamer.Close()
	// done := make(chan bool)
	// speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// 	done <- true
	// })))
	// <-done
	return
}
