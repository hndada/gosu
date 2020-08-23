package gosu

import (
	"fmt"
	"github.com/hndada/gosu/config"
	"image"
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
	Chart    *mania.Chart
	Cover    *ebiten.Image
	Velocity float64
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
	// 노트 이미지가 전부 불러와졌다고 가정
	s.Tick++

	for i := range s.Notes {
		s.Notes[i].y += s.Tick() * scrollSpeed
	}
	if s.Time() > endTime {
		s.Streamer.Close()
		g.ChangeScene(NewSceneSelect())
	}

	// next note들이 slice에 fetch되어 있는 상태. hit 판정이 나왔을 경우.
	if done {
		ops[i].ColorM.ChangeHSV(0, 0, 0.5) // gray
	}

	return nil
}

// op은 노트와 일대일 대응되므로 별도의 타입이나 idx 필요 없음
// 롱노트 body는 어떻게 해야할지 걱정
func loadNotes() {
	// TransPoint 좌르륵 읽기
	// 다음 TransPoint-현 TransPoint
	ops := make([]ebiten.DrawImageOptions, len(notes))
	var notes []mania.Note

	// 필드의 중앙이 스크린의 중앙에 오게
	// noteimgs := config.NoteImages(7)

	var ps []image.Point
	var fieldWidth float64
	for _, p := range ps {
		fieldWidth += p.X
	}
	// 아래 목록은 screen에 처음부터 fixed 된 상태로 그려질 애들
	// hint size: fieldWidth, ps[0].Y (indexing했으니 앞에서 0키면 에러 리턴) (근데 Y 그대로 쓰면 두꺼운 노트 스킨인 경우 어색할 수도 있음)
	// main size: fieldWidth, screenHeight
	// bottom size: fieldWidth, scaled src height
	// left, right, HPBarBG size: scaled src width, screenHeight

	X(fieldWidth)

	for i, n := range notes {
		n.Key
	}
}

func X(w, center float64) float64 {
	var halfWidth = float64(screenWidth / 2)
	offset := halfWidth * (center / 100)
	return halfWidth + offset - w/2
}

// 맨날 Update()에 game을 넘겨주지 말고 scene struct의 필드 값으로 박자
// 이미지를 관리할게 아니라 op을 관리해야하네
func (s *SceneMania) Draw(screen *ebiten.Image) {
	// BasePlayScene 그리기
	op := &ebiten.DrawImageOptions{}
	const ratio = float64(1600) / 1920
	op.GeoM.Scale(ratio, ratio)
	op.ColorM.ChangeHSV(0, 1, 0.30)
	screen.DrawImage(s.Cover, op)

	// edit를 위해서도, 단순하게 짜기: 전체 그리기
	for i := range ops {
		ops[i].GeoM.Translate(0, tickTime*speed*speedFactor)
		screen.DrawImage(noteimgs[keys], &ops[i])
	}

	// 키 버튼 그리기

	// ebitenutil.DrawRect(screen, 565, 0, 70*7, 900, color.RGBA{0, 0, 0, 180})
	// ebitenutil.DrawRect(screen, 565, 730, 70*7, 10, color.RGBA{252, 106, 111, 255})
	// for i := range s.Notes {
	// 	ebitenutil.DrawRect(screen, s.Notes[i].x, s.Notes[i].y, s.Notes[i].w, s.Notes[i].h, s.Notes[i].clr)
	// }

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
		s.Notes[i] = NoteImageInfo{x, y, w, h, settings.noteColor(n, c.Keys)}
	}

	bg, err := s.Chart.Background()
	if err != nil {
		log.Fatal(err) // todo: normal log
	}
	s.Cover, _ = ebiten.NewImageFromImage(bg, ebiten.FilterDefault)
	// s.ScrollSpeed = g.ScrollSpeed // todo: TransPoint로 설정. 변속 같은 거도 조정해야 하니 좀 더 잘 다뤄야함
	return s
}
