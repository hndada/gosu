package gosu

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/mania"
	"github.com/moutend/go-hook/pkg/types"
	_ "github.com/silbinarywolf/preferdiscretegpu"
	"image"
	_ "image/jpeg"
	"sort"
)

// 판정선 가운데에 노트 가운데가 맞을 때 Max가 뜨게 (거의 다 구현한듯)
// 내 생각에 지금 느린건 (혹은 느리다고 보이는건) audio가 Time을 제대로 안내주기 때문인거 같음
// 이전 값에 상관없이 언제나 다시 그리므로 applySpeed()가 따로 필요 없음
// 최종 이미지는 언제나 사이즈가 int, int이므로 image.Point로 다뤄도 됨

// todo: timing points, (decending/ascending) order로 sort -> rg-parser에서
type NoteSprite struct {
	h        int     // same with second value of NoteSprite.i.Size()
	x        float64 // x is fixed among mania notes
	position float64 // positive value
	y        float64
	i        *ebiten.Image
	op       *ebiten.DrawImageOptions
}

type LNSprite struct {
	head   *NoteSprite
	tail   *NoteSprite
	length float64

	height float64
	i      *ebiten.Image
	bodyop *ebiten.DrawImageOptions
}

type SceneMania struct { // aka Clavier
	g            *Game
	mods         mania.Mods
	chart        *mania.Chart
	speedFactors []mode.SpeedFactorPoint
	stamps       []timeStamp
	stage        mania.Stage // for quick access
	notes        []NoteSprite
	lnotes       []LNSprite // 롱노트 특성상, 2개로 나누는 게 불가피해보임

	bg   *ebiten.Image
	bgop *ebiten.DrawImageOptions

	speed        float64
	hitPosition  float64
	displayScale float64

	audioPlayer *AudioPlayer
	// sfxBuffer map[string]*beep.Buffer
	kbEventChan input.KeyboardEventChannel
	layout      []types.VKCode
	endTime     int64

	tick     int64
	score    float64
	hp       float64
	combo    int32
	stampIdx int
	logs     []keyLog
}

// tools.Stamp를 통해서 구현하려 했다가 element 설정에서 fail
// 각 타입 별로 (float64 등) 만들면 그때 비로소 쓸 수 있을 듯
type timeStamp struct {
	time     int64
	nextTime int64
	position float64
	factor   float64
}

type keyLog struct {
	time  int64
	state []bool
}

// time series; acending order
func SearchKeyLog(logs []keyLog, time int64) int {
	idx := sort.Search(len(logs), func(i int) bool { return logs[i].time >= time })
	if idx < len(logs) && logs[idx].time == time {
		return idx
	}
	return -1
}

// 없다면, 추가하고 Sort해야함
func SortKeyLogs(logs []keyLog) {
	sort.Slice(logs, func(i, j int) bool { return logs[i].time < logs[j].time })
}

// lnhead와 lntail 분리 유지
func (g *Game) NewSceneMania(c *mania.Chart, mods mania.Mods) *SceneMania {
	s := &SceneMania{}
	s.g = g
	s.mods = mods
	s.chart = c.ApplyMods(s.mods)
	// todo: 노트가 언제나 양수 시간에 있다고 상정; 실제로는 노트가 BufferTime보다 뒤에 있을 수 있음
	initSpeedFactor := mode.SpeedFactorPoint{0, 1}
	s.speedFactors = append([]mode.SpeedFactorPoint{initSpeedFactor}, s.chart.TimingPoints.SpeedFactors...)
	s.stamps = make([]timeStamp, len(s.speedFactors))
	var position float64
	for i, sf := range s.speedFactors {
		var stamp timeStamp
		stamp.time = sf.Time
		stamp.position = position
		stamp.factor = sf.Factor
		if i < len(s.speedFactors)-1 {
			nextTime := s.speedFactors[i+1].Time
			stamp.nextTime = nextTime
			position += float64(nextTime-sf.Time) * sf.Factor
		} else {
			stamp.nextTime = 9223372036854775807 // max int64
		}
		s.stamps[i] = stamp
	}
	s.stage = s.g.sprites.mania.Stages[s.chart.Keys]
	s.notes = make([]NoteSprite, len(s.chart.Notes))
	for i, n := range s.chart.Notes {
		var ns NoteSprite
		var sprite mode.Sprite
		switch n.Type {
		case mania.TypeNote:
			sprite = s.stage.Notes[n.Key]
		case mania.TypeLNHead:
			sprite = s.stage.LNHeads[n.Key]
		case mania.TypeLNTail:
			sprite = s.stage.LNTails[n.Key]
		}
		ns.x = float64(sprite.Position().X)
		ns.i = sprite.Image()
		_, ns.h = ns.i.Size()
		ns.op = &ebiten.DrawImageOptions{}
		s.notes[i] = ns
	}
	// range stamps를 outer loop로 두고 짜도 큰 차이 없을 듯
	var stampIdx int
	stamp := s.stamps[0]
	for i, n := range s.chart.Notes {
		for si := range s.stamps[stampIdx:] {
			if n.Time < s.stamps[stampIdx+si].nextTime {
				if si != 0 {
					stamp = s.stamps[stampIdx+si]
					stampIdx += si
				}
				break
			}
		}
		s.notes[i].position = float64(n.Time-stamp.time)*stamp.factor + stamp.position
	}
	s.lnotes = make([]LNSprite, 0, s.chart.NumLN())
	lastLNHeads := make([]int, s.chart.Keys)
	for i, n := range s.chart.Notes {
		switch n.Type {
		case mania.TypeLNHead:
			lastLNHeads[n.Key] = i
		case mania.TypeLNTail:
			var ls LNSprite
			ls.head = &s.notes[lastLNHeads[n.Key]]
			ls.tail = &s.notes[i]
			ls.i = s.stage.LNBodys[n.Key][0].Image()
			ls.length = ls.tail.position - ls.head.position
			ls.bodyop = &ebiten.DrawImageOptions{}
			s.lnotes = append(s.lnotes, ls)
		}
	}

	var err error
	s.bg, err = s.chart.Background()
	if err != nil {
		panic(err)
	}
	s.bgop = mode.BackgroundOp(s.g.settings.ScreenSize(), image.Pt(s.bg.Size()))
	var dimness uint8
	switch {
	default:
		dimness = s.g.settings.GeneralDimness()
	}
	// dim 을 바꾸는 입력이 들어왔다면 별도 함수 없이 즉석에서 s.bgop.ColorM.Reset() 날리고 다시 설정.
	s.bgop.ColorM.ChangeHSV(0, 1, float64(dimness)/100)

	switch {
	default:
		s.speed = s.g.settings.mania.GeneralSpeed
	}
	s.hitPosition = s.g.settings.mania.HitPosition
	s.displayScale = s.g.settings.ScaleY()

	s.audioPlayer = NewAudioPlayer(s.g.audioContext, s.chart.AbsPath(s.chart.AudioFilename))
	s.kbEventChan, err = input.NewKeyboardEventChannel()
	if err != nil {
		panic(err)
	}
	s.layout = s.g.settings.mania.KeyLayout[s.chart.Keys]
	s.endTime = s.chart.EndTime()

	// const BufferTime = 2
	// s.tick = -int64(BufferTime * s.g.MaxTPS())
	s.logs = make([]keyLog, 0)
	return s
}

func (s *SceneMania) Update() error {
	if s.Time() > s.endTime || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if err := s.kbEventChan.Close(); err != nil {
			return err
		}
		_ = s.audioPlayer.Close()
		s.g.changeScene(s.g.NewSceneSelect())
		return nil
	}
	for _, e := range s.kbEventChan.Dequeue() {
		for key, keycode := range s.layout {
			if keycode == e.KeyCode {
				if SearchKeyLog(s.logs, e.Time) == -1 {
					s.logs[e.Time].state = make([]bool, s.chart.Keys)
				}
				if e.State == input.KeyStateDown {
					s.logs[e.Time].state[key] = true
				} else {
					s.logs[e.Time].state[key] = false
				}
				break
			}
		}
		// todo: 스코어 처리 함수에 s.log와 s.logTime 보내기
		// 리플레이 구조: 마지막 status 시간, 레이아웃 키state
		// next note들이 slice에 fetch되어 있는 상태.
		// if done {
		// 	s.notes[i].op.ColorM.ChangeHSV(0, 0, 0.5) // gray
		// }
	}

	now := s.Time()
	// fmt.Printf("tick:%d, time: %dms\n", s.tick, now)
	var stamp timeStamp
	for si := range s.stamps[s.stampIdx:] {
		if now < s.stamps[s.stampIdx+si].nextTime {
			stamp = s.stamps[s.stampIdx+si]
			s.stampIdx += si
			break
		}
	}

	measure := float64(now-stamp.time)*stamp.factor + stamp.position
	for i, n := range s.notes {
		pos := (n.position-measure)*s.speed - s.hitPosition
		s.notes[i].y = -pos*s.displayScale + float64(n.h)/2
	}
	for i, n := range s.lnotes {
		s.lnotes[i].height = n.length * s.speed * s.displayScale
	}
	s.tick++
	return nil
}

func (s *SceneMania) Draw(screen *ebiten.Image) {
	// 키 버튼 그리기
	// 스코어, hp, 콤보, 시간 그리기; hp는 마스크 이미지를 씌우면 되지 않을까
	screen.DrawImage(s.bg, s.bgop)
	screen.DrawImage(s.stage.Fixed.Image(), &ebiten.DrawImageOptions{})
	for _, n := range s.lnotes {
		n.DrawLN(screen)
	}
	// op를 매번 생성하는 게 더 빠를까? 근데 그럴 것 같지는 않아
	// 화면 범위 바깥은 Draw 생략해야할까?
	for i, n := range s.notes {
		s.notes[i].op.GeoM.Reset()
		s.notes[i].op.GeoM.Translate(n.x, n.y)
		screen.DrawImage(n.i, s.notes[i].op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		`CurrentFPS: %.2f
CurrentTPS: %.2f
Time: %.3fs`, ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(s.Time())/1000))
}

func (s LNSprite) DrawLN(screen *ebiten.Image) {
	_, h := s.i.Size()
	count, remainder := int(s.height)/h, int(s.height)%h+1
	s.bodyop.GeoM.Reset()
	s.bodyop.GeoM.Translate(s.tail.x, s.tail.y)

	firstRect := s.i.Bounds()
	firstRect.Min = image.Pt(0, h-remainder)
	screen.DrawImage(s.i.SubImage(firstRect).(*ebiten.Image), s.bodyop)
	s.bodyop.GeoM.Translate(0, float64(remainder))

	for c := 0; c < count; c++ {
		screen.DrawImage(s.i, s.bodyop)
		s.bodyop.GeoM.Translate(0, float64(h))
	}
}
func (s *SceneMania) Init() {
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", s.chart.MusicName, s.chart.ChartName))
	s.audioPlayer.Play()
}

// always follows audio's time
// func (s *SceneMania) Time() int64 {
// 	return s.audioPlayer.Time().Milliseconds()
// }

// 이 방법을 하려면 tps가 게임 중에 변하지 않아야 함
// CurrentTPS가 약간 딱 떨어지지 않는 게 마음에 걸리지만, 곧 보충되어 결과적으로 일정히 유지 된다고 상정하겠음
// -> Audio에서 Time 따오는 게 제일 정확. 그런데 지금 오디오가 내주는 시간이 버퍼에 의해 정확하지 않음
func (s *SceneMania) Time() int64 {
	return Millisecond * s.tick / int64(s.g.settings.MaxTPS())
}

// 노트 효율적으로 하강시키기 위해 시도했던 방법 중 하나
// 모종의 이유로 update가 누락되어도 오디오가 재생되는 거에 맞춰서 스크린 그려지게 하려면 이 방법은 쓰면 안됨
// const ApproachDuration = 1500
// step1ms := float64(s.g.ScreenSize().Y) / ApproachDuration
// s.step = func(ms int64) float64 { return float64(ms) * step1ms }
