package gosu

import (
	"fmt"
	"github.com/hndada/gosu/graphics"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/moutend/go-hook/pkg/types"
	"image"
	_ "image/jpeg"
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
	sfactors   []mode.SpeedFactorPoint
	sfactorIdx int

	stage       graphics.ManiaStage
	layout      []types.VKCode
	kbEventChan input.KeyboardEventChannel

	log map[int64][]bool
	logTime
}

type logTime []int64

func (l logTime) isLogged(t int64) bool {
	for i := len(l) - 1; i >= 0; i-- {
		if l[i] == t {
			return true
		}
	}
	return false
}

// lnhead와 lntail 분리 유지
func (g *Game) NewSceneMania(c *mania.Chart, mods mania.Mods) *SceneMania {
	s := &SceneMania{}
	s.g = g
	// 얘도 별도 함수 없이 여기서 처리하는게 맞을듯
	// BaseChart는 ScreenSize에 대해서 알지 못함
	// mode 쪽은 ebiten으로부터 독립적이었으면 좋겠음
	s.chart = c.ApplyMods(mods)
	bg, err := s.chart.Background()
	if err != nil {
		panic(err)
	}
	s.bg, _ = ebiten.NewImageFromImage(bg, ebiten.FilterDefault)

	// s.bgop = &ebiten.DrawImageOptions{}
	// s.bgop.GeoM.Scale(FitRatio(s.g.ScreenSize(), image.Pt(s.bg.Size()), FixRatioModeMin))
	s.bgop = BackgroundOp(s.g.ScreenSize(), image.Pt(s.bg.Size()))

	var dimness uint8
	switch {
	default:
		dimness = s.g.Settings.GeneralDimness()
	}
	// 별도 함수 없이 여기서 처리하는게 맞을듯
	// 만약 dim 을 바꾸는 입력이 들어왔다면 즉석에서 s.bgop.ColorM.Reset() 날리고 다시 설정.
	s.bgop.ColorM.ChangeHSV(0, 1, float64(dimness)/100)

	{
		f, err := os.Open(s.chart.AbsPath(s.chart.AudioFilename))
		if err != nil {
			panic(err)
		}
		var streamer beep.StreamSeekCloser
		var streamFormat beep.Format
		switch strings.ToLower(filepath.Ext(s.chart.AudioFilename)) {
		case ".mp3":
			streamer, streamFormat, err = mp3.Decode(f)
			if err != nil {
				panic(err)
			}
		}
		s.buffer = beep.NewBuffer(streamFormat)
		s.buffer.Append(beep.Silence(streamFormat.SampleRate.N(2 * time.Second)))
		s.buffer.Append(streamer)
		_ = streamer.Close()
	}
	const BufferTime = 2
	s.tick = -int64(BufferTime * s.g.MaxTPS())

	const ApproachDuration = 1500
	step1ms := float64(s.g.ScreenSize().Y) / ApproachDuration
	s.step = func(ms int64) float64 { return float64(ms) * step1ms }

	// todo: 노트가 언제나 양수 시간에 있다고 상정; 실제로는 노트가 BufferTime보다 뒤에 있을 수 있음
	initSpeedFactor := mode.SpeedFactorPoint{-BufferTime * Millisecond, 1}
	// s.progress = float64(-BufferTime*Millisecond-initSpeedFactor.Time) * initSpeedFactor.Factor
	s.sfactors = append([]mode.SpeedFactorPoint{initSpeedFactor}, s.chart.TimingPoints.SpeedFactors...)

	s.stage = s.g.GameSprites.ManiaStages[s.chart.Keys] // for quick access
	s.setNoteSprites()
	var speed float64
	switch {
	default:
		speed = s.g.Settings.GeneralSpeed
	}
	s.applySpeed(speed)
	s.endTime = s.chart.EndTime()

	s.layout = s.g.Settings.KeyLayout[s.chart.Keys]
	s.kbEventChan, err = input.NewKeyboardEventChannel()
	if err != nil {
		panic(err)
	}
	return s
}

// 리플레이 구조: 마지막 status 시간, 레이아웃 키state
func (s *SceneMania) Update() error {
	if s.Time(s.tick) > s.endTime || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if err := s.kbEventChan.Close(); err != nil {
			return err
		}
		speaker.Close()
		s.g.changeScene(s.g.NewSceneSelect())
		return nil
	}

	for _, e := range s.kbEventChan.Dequeue() {
		for key, keycode := range s.layout {
			if keycode == e.KeyCode {
				if !s.logTime.isLogged(e.Time) {
					s.log[e.Time] = make([]bool, s.chart.Keys)
				}
				if e.State == input.KeyStateDown {
					s.log[e.Time][key] = true
				} else {
					s.log[e.Time][key] = false
				}
				break
			}
		}
		// todo: 스코어 처리 함수에 s.log와 s.logTime 보내기
		// next note들이 slice에 fetch되어 있는 상태.
		// if done {
		// 	s.notes[i].op.ColorM.ChangeHSV(0, 0, 0.5) // gray
		// }
	}

	var dy float64
	now := s.Time(s.tick)
	lastTime := s.Time(s.tick - 1)
	for si, sp := range s.sfactors[s.sfactorIdx:] { // todo: timing points sorting -> rg-parser에서
		if si == len(s.sfactors)-1 {
			dy += s.step(now-lastTime) * s.speed * sp.Factor
			s.sfactorIdx = si
		} else {
			nextsp := s.sfactors[si+1]
			if nextsp.Time > now { // equality condition should be excluded
				s.sfactorIdx = si
				break
			}
			dy += s.step(nextsp.Time-lastTime) * s.speed * sp.Factor
			lastTime = nextsp.Time
		}
	}
	// for si, sp := range sfactors[s.sfactorIdx+1:] {
	// 	if sp.Time > s.Time(s.tick) { // equality condition should be excluded
	// 		break
	// 	}
	// 	dy += s.step(sp.Time-lastTime) * s.speed * sp.Factor
	// 	lastTime = sp.Time
	// 	s.sfactorIdx = si
	// }
	dy *= 1.7 // todo: why?
	for i := range s.notes {
		s.notes[i].op.GeoM.Translate(0, dy)
	}
	for i := range s.lnotes {
		s.lnotes[i].bodyop.GeoM.Translate(0, dy)
	}
	s.progress += dy
	s.tick++
	// fmt.Printf("%d: %f\n", lastTime, dy)
	return nil
}

// 단순명료하게 전체 노트 그리기; edit 위해서도 좋음
// (op을 Update()에서 바꾸는 게 나는 바람직해 보이는데, 표준인진 모르겠음)
func (s *SceneMania) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.bg, s.bgop)
	screen.DrawImage(s.stage.Fixed.Image(), &ebiten.DrawImageOptions{})
	// 키 버튼 그리기
	// 스코어, hp, 콤보, 시간 그리기
	// s.g.Sprites.Score
	// s.g.Sprites.ManiaCombo
	// hp는 마스크 이미지를 씌우면 되지 않을까
	for _, n := range s.lnotes {
		screen.DrawImage(n.i, n.bodyop)
	}
	for _, n := range s.notes {
		screen.DrawImage(n.i, n.op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTime: %.1fs", ebiten.CurrentFPS(), float64(s.Time(s.tick))/1000))
}

// downcast happens
func (s *SceneMania) Time(tick int64) int64 {
	return Millisecond * tick / int64(s.g.MaxTPS())
}

func (s *SceneMania) Init() {
	if err := speaker.Init(s.buffer.Format().SampleRate,
		s.buffer.Format().SampleRate.N(time.Second/10)); err != nil {
		panic(err)
	}
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", s.chart.MusicName, s.chart.ChartName))
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

func BackgroundOp(screen, bg image.Point) *ebiten.DrawImageOptions {
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
