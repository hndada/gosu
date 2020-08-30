package gosu

import (
	"fmt"
	"github.com/hndada/gosu/config"
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
// todo: image.Point가 int,int만 받는걸 이제 봤네
type SceneMania struct { // aka Clavier
	g      *Game
	bg     *ebiten.Image
	bgop   *ebiten.DrawImageOptions
	buffer *beep.Buffer // todo: bgmBuffer
	// sfxBuffer
	// map[string]*beep.Buffer
	tick int64
	step func(ms int64) float64

	score float64
	hp    float64
	combo int32

	chart   *mania.Chart
	notes   []NoteSprite
	lnotes  []LNSprite // 롱노트 특성상, 2개로 나누는 게 불가피해보임
	endTime int64

	speed      float64
	progress   float64
	sfactorIdx int

	stage config.ManiaStage
}

// lnhead와 lntail 분리 유지
func (g *Game) NewSceneMania(c *mania.Chart, mods mania.Mods) *SceneMania {
	s := &SceneMania{}
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", c.SongName, c.ChartName))
	// 얘도 별도 함수 없이 여기서 처리하는게 맞을듯
	// BaseChart는 ScreenSize에 대해서 알지 못함
	// mode 쪽은 ebiten으로부터 독립적이었으면 좋겠음
	bg, err := s.chart.Background()
	if err != nil {
		log.Fatal(err)
	}
	s.bg, _ = ebiten.NewImageFromImage(bg, ebiten.FilterDefault)

	// s.bgop = &ebiten.DrawImageOptions{}
	// s.bgop.GeoM.Scale(FitRatio(s.g.ScreenSize(), image.Pt(s.bg.Size()), FixRatioModeMin))
	s.bgop = BackgroundOp(s.g.ScreenSize(), image.Pt(s.bg.Size()))

	var dimness uint8
	switch {
	default:
		dimness = s.g.Settings.GeneralDimness
	}
	// 별도 함수 없이 여기서 처리하는게 맞을듯
	// 만약 dim 을 바꾸는 입력이 들어왔다면 즉석에서 s.bgop.ColorM.Reset() 날리고 다시 설정.
	s.bgop.ColorM.ChangeHSV(0, 1, float64(dimness)/100)

	{
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
	}
	const BufferTime = 2
	s.tick = -int64(BufferTime * s.g.MaxTPS())

	const ApproachDuration = 1500
	step1ms := float64(s.g.ScreenSize().Y) / ApproachDuration
	s.step = func(ms int64) float64 { return float64(ms) * step1ms }

	s.chart = c.ApplyMods(mods)
	s.setNoteSprites()
	var speed float64
	switch {
	default:
		speed = s.g.Settings.GeneralSpeed
	}
	s.applySpeed(speed)
	s.endTime = s.chart.EndTime()

	s.stage = s.g.Sprites.ManiaStages[s.chart.Keys] // for quick access
	return s
}

func (s *SceneMania) Update() error {
	if s.Time(s.tick) > s.endTime {
		s.g.ChangeScene(NewSceneSelect())
	}
	// todo: 플레이 하면서 리플레이 데이터 저장
	// 키보드 입력 채널에서 키 입력 불러오기
	// next note들이 slice에 fetch되어 있는 상태. hit 판정이 나왔을 경우.
	// if done {
	// 	s.notes[i].op.ColorM.ChangeHSV(0, 0, 0.5) // gray
	// }

	var dy float64
	lastTime := s.Time(s.tick - 1)
	sfactors := s.chart.TimingPoints.SpeedFactors
	for si, sp := range sfactors[s.sfactorIdx+1:] {
		if sp.Time > s.Time(s.tick) { // equality condition should be excluded
			break
		}
		dy += s.step(sp.Time-lastTime) * s.speed * sp.Factor
		lastTime = sp.Time
		s.sfactorIdx = si
	}
	for i := range s.notes {
		s.notes[i].op.GeoM.Translate(0, dy)
	}
	for i := range s.lnotes {
		s.lnotes[i].bodyop.GeoM.Translate(0, dy)
	}
	s.progress += dy
	s.tick++
	return nil
}

// 단순명료하게 전체 노트 그리기; edit 위해서도 좋음
// (op을 Update()에서 바꾸는 게 나는 바람직해 보이는데, 표준인진 모르겠음)
func (s *SceneMania) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.bg, s.bgop)
	// 스테이지 그리기
	// 키 버튼 그리기
	// 스코어, hp, 콤보, 시간 그리기
	// s.g.Sprites.Score
	// s.g.Sprites.ManiaCombo
	// hp는 마스크 이미지를 씌우면 되지 않을까
	for _, n := range s.notes {
		screen.DrawImage(n.i, n.op)
	}
	for _, n := range s.lnotes {
		screen.DrawImage(n.i, n.bodyop)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTime: %.1fs", ebiten.CurrentFPS(), float64(s.Time(s.tick))/1000))
}

// downcast happens
func (s *SceneMania) Time(tick int64) int64 {
	return Millisecond * tick / int64(s.g.MaxTPS())
}

func (s *SceneMania) Init() {
	speaker.Play(s.buffer.Streamer(0, s.buffer.Len()))
}

// // FitRatio returns fitting resizing ratio with maintaining aspect ratio.
// type FixRatioMode uint8
//
// const (
// 	FixRatioModeMax FixRatioMode = iota
// 	FixRatioModeMin
// )
//
// func FitRatio(dst, src image.Point, mode FixRatioMode) (float64, float64) {
// 	rx := float64(dst.X) / float64(src.X)
// 	ry := float64(dst.Y) / float64(src.Y)
// 	max, min := rx, ry
// 	if max < min {
// 		max, min = min, max
// 	}
// 	switch mode {
// 	case FixRatioModeMax:
// 		return max, max
// 	case FixRatioModeMin:
// 		return min, min
// 	default:
// 		panic("not reach")
// 	}
// }

func BackgroundOp(bg, screen image.Point) *ebiten.DrawImageOptions {
	// 그림의 크기가 스크린보다 클 경우 줄이기
	// rx := float64(bg.X) / float64(screen.X)
	// ry := float64(bg.Y) / float64(screen.Y)
	// var adj float64
	// if rx > 1 || ry > 1 {
	// 	max := rx
	// 	if max < ry {
	// 		max = ry
	// 	}
	// 	adj = 1 / max
	// 	op.GeoM.Scale(adj, adj)
	// }

	op := &ebiten.DrawImageOptions{}
	bx, by := float64(bg.X), float64(bg.Y)
	sx, sy := float64(screen.X), float64(screen.Y)

	rx, ry := sx/bx, sy/by
	var ratio float64 = 1
	if rx < 1 || ry < 1 { // 스크린이 그림보다 작을 경우 그림 크기 줄이기
		min := rx
		if min > ry {
			min = ry
		}
		ratio = min
		op.GeoM.Scale(ratio, ratio)
	}

	// 그림이 모니터의 중앙에 위치하게
	// x와 y 둘 중 하나는 스크린 크기와 일치; 둘 모두 크기가 스크린보다 작거나 같다
	x, y := bx*ratio, by*ratio
	op.GeoM.Translate((sx-x)/2, (sy-y)/2)

	return op
}
