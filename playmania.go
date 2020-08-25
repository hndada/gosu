package gosu

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/mode/mania"
)

// todo: NewBasePlayScene()
// todo: examples/camera 참고하여 플레이 프레임 안정화
type SceneMania struct { // aka Clavier
	// BasePlayScene
	buffer *beep.Buffer
	g      *Game

	tick int64
	// score float64
	// HP    float64
	// Combo int32

	chart *mania.Chart
	// Cover  *ebiten.Image
	backStage   *ebiten.Image
	speed       float64
	speedFactor float64
	notes       []NoteSprite
	lnotes      []LNSprite // 롱노트 특성상, 2개로 나누는 게 불가피해보임
	viewport    float64
	step        func() float64
}

// lnhead와 lntail 분리 유지
func (g *Game) NewSceneMania(c *mania.Chart, mods mania.Mods) *SceneMania {
	s := &SceneMania{}
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", c.SongName, c.ChartName))
	s.tick = -int64(2 * s.g.MaxTPS())

	f, err := os.Open(s.chart.AbsPath(s.chart.AudioFilename))
	if err != nil {
		log.Fatal(err)
	}
	var streamer beep.StreamSeekCloser
	var streamFormat beep.Format
	switch strings.ToLower(filepath.Ext(s.chart.AudioFilename)) {
	case ".mp3":
		streamer, streamFormat, err = mp3.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	s.buffer = beep.NewBuffer(streamFormat)
	s.buffer.Append(beep.Silence(streamFormat.SampleRate.N(2 * time.Second)))
	s.buffer.Append(streamer)
	_ = streamer.Close()

	if err = speaker.Init(s.buffer.Format().SampleRate,
		s.buffer.Format().SampleRate.N(time.Second/100)); err != nil {
		log.Fatal(err)
	}

	s.chart = c.ApplyMods(mods)

	// 스테이지 + StartAt Offset op걸고
	// ebitenutil.DrawRect(screen, 565, 0, 70*7, 900, color.RGBA{0, 0, 0, 180})
	// ebitenutil.DrawRect(screen, 565, 730, 70*7, 10, color.RGBA{252, 106, 111, 255})

	// var ps []image.Point
	// var fieldWidth float64
	// for _, p := range ps {
	// 	fieldWidth += p.X
	// }
	// 아래 목록은 screen에 처음부터 fixed 된 상태로 그려질 애들
	// hint size: fieldWidth, ps[0].Y (indexing했으니 앞에서 0키면 에러 리턴)
	// (근데 Y 그대로 쓰면 두꺼운 노트 스킨인 경우 어색할 수도 있음)
	// todo: 판정선 가운데에 노트 가운데가 맞을 때 Max가 뜨게
	// main size: fieldWidth, screenHeight
	// bottom size: fieldWidth, scaled src height
	// left, right, HPBarBG size: scaled src width, screenHeight

	// 필드의 중앙이 스크린의 중앙에 오게
	// ColumnStart를 리턴하는 함수로 만드는 게 좋을 듯
	// func X(w, center float64) float64 {
	// 	var halfWidth = float64(screenWidth / 2)
	// 	offset := halfWidth * (center / 100)
	// 	return halfWidth + offset - w/2
	// }

	s.setNoteSprites()
	s.applySpeed(s.g.ScrollSpeed)
	// todo: 값 변경 안생기게 함수로 감쌌는데, performance 확인
	const approachDuration = 1000
	s.step = func() float64 { return float64(s.g.ScreenSize().Y) / float64(s.g.MaxTPS()) / (approachDuration / Millisecond) }
	{
		_bg, err := s.chart.Background()
		if err != nil {
			log.Fatal(err)
		}
		bg, _ := ebiten.NewImageFromImage(_bg, ebiten.FilterDefault)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(ratio(s.g.ScreenSize(), image.Pt(bg.Size()))) // todo: 폭맞춤
		op.ColorM.ChangeHSV(0, 1, 0.30)
		s.backStage.DrawImage(bg, op)
	}

	return s
}

func (s *SceneMania) Update() error {
	// speedFactor 적용
	// todo: 1프레임 사이에 speedFactor가 바뀔 경우, 그 잠깐의 간격에서 발생하는 오차
	s.tick++
	if s.Time() > s.chart.EndTime() {
		s.g.ChangeScene(NewSceneSelect())
	}
	// todo: 플레이 하면서 리플레이 데이터 저장
	// 키보드 입력 채널에서 키 입력 불러오기
	// next note들이 slice에 fetch되어 있는 상태. hit 판정이 나왔을 경우.
	if done {
		s.notes[i].op.ColorM.ChangeHSV(0, 0, 0.5) // gray
	}

	dy := s.step() * s.speed * s.speedFactor
	// todo: op을 Update()에서 바꾸는 게 나는 바람직해 보이는데, 표준은?
	for i := range s.notes {
		s.notes[i].op.GeoM.Translate(0, dy)
	}
	for i := range s.lnotes {
		s.lnotes[i].bodyop.GeoM.Translate(0, dy)
	}
	s.viewport += dy
	return nil
}

// 단순명료하게 전체 노트 그리기; edit 위해서도 좋음
func (s *SceneMania) Draw(screen *ebiten.Image) {
	// BackStage 그리기
	for _, n := range s.notes {
		screen.DrawImage(noteimgs[keys], &n.op)
	}
	for _, n := range s.lnotes {
		screen.DrawImage(noteimgs[keys], &n.bodyop)
	}
	// 키 버튼 그리기
	// 스코어, hp, 콤보, 시간 그리기
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTime: %.1fs", ebiten.CurrentFPS(), float64(s.Time())/1000))
}

func (s *SceneMania) Time() int64 {
	return s.tick * Millisecond / int64(ebiten.MaxTPS())
}

func (s *SceneMania) Init() {
	speaker.Play(s.buffer.Streamer(0, s.buffer.Len()))
}

func ratio(dst, src image.Point) (float64, float64) {
	return float64(dst.X) / float64(src.X),
		float64(dst.Y) / float64(src.Y)
}
