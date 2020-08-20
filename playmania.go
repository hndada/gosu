package gosu

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/mode/mania"
)

// todo: 노트 이미지 미리 렌더 후 그리기만 하기
// 채보를 미리 전부 그려놔야 할까?
// todo: 플레이 하면서 리플레이 데이터 저장
type SceneMania struct { // aka Clavier
	BasePlayScene
	Chart       *mania.Chart
	Cover       *ebiten.Image
	ScrollSpeed float64
}

// todo: 로딩일 때 기다리는 로직
// Loading 이라는 별도의 Lock을 둔 이상, 특별히 채널은 필요없는거 아닌가?
// 비트맵 로딩 15초 후 timeout
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

// 판정선 가운데에 노트 가운데가 맞을 때 Max가 뜨게
// 모드에 맞추어서 Notes 와 TransPoint 변경

// 스피드값 1 기준 초당 1000픽셀 내려오게 해야함
func (g *Game) NewSceneMania(c *mania.Chart, mods mania.Mods) *SceneMania {
	s := &SceneMania{}
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", c.SongName, c.ChartName)) // todo: NewBasePlayScene()
	s.Tick = func() float64 { return 1000 / float64(g.MaxTPS) }                   // todo: NewBasePlayScene()
	s.Time = -PlaySceneBufferTime                                                 // todo: NewBasePlayScene()
	f, err := os.Open(s.Chart.AbsPath(s.Chart.AudioFilename))                     // todo: NewBasePlayScene()
	if err != nil {
		log.Fatal(err)
	}
	s.Streamer, s.StreamFormat, err = mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	s.Chart = c.ApplyMods(mods)
	for i, n := range s.Chart.Notes {
		var y, h float64
		x := float64(n.Key*w + 565)
		switch n.Type {
		case mania.TypeNote:
			y = (-float64(n.Time)+PlaySceneBufferTime)*s.ScrollSpeed + 730
			h = noteHeight * s.ScrollSpeed
		case mania.LNHead:
			y = (-float64(n.Time2)+PlaySceneBufferTime)*s.ScrollSpeed + 730
			h = float64(n.Time2-n.Time+noteHeight) * s.ScrollSpeed
		}
		s.Notes[i] = NoteImageInfo{x, y, w, h, noteColor(n, c.Keys)}
	}

	bg, err := s.Chart.Background()
	if err != nil {
		log.Fatal(err) // todo: normal log
	}
	s.Cover, _ = ebiten.NewImageFromImage(bg, ebiten.FilterDefault)
	// s.ScrollSpeed = g.ScrollSpeed // todo: TransPoint로 설정. 변속 같은 거도 조정해야 하니 좀 더 잘 다뤄야함
	return s
}
