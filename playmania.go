package gosu

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/graphics"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/mania"
	"github.com/moutend/go-hook/pkg/types"
	"image"
	_ "image/jpeg"
	"sort"
)

// todo: 판정선 가운데에 노트 가운데가 맞을 때 Max가 뜨게
// todo: image.Point가 int,int만 받는걸 이제 봤네
// todo: speed factors, (decending/ascending) order로 sort
type NoteSprite struct {
	h        int     // same with second value of NoteSprite.i.Size()
	x        float64 // x is fixed among mania notes
	position float64 // 양수
	y        float64
	i        *ebiten.Image
	op       *ebiten.DrawImageOptions
}

type LNSprite struct {
	head   *NoteSprite
	tail   *NoteSprite
	length float64
	i      *ebiten.Image
	bodyop *ebiten.DrawImageOptions
}

func (s LNSprite) height() float64 {
	return s.tail.position - s.head.position
}

type SceneMania struct { // aka Clavier
	g    *Game
	bg   *ebiten.Image
	bgop *ebiten.DrawImageOptions
	// buffer *beep.Buffer // todo: bgmBuffer
	audioPlayer *AudioPlayer
	// sfxBuffer
	// map[string]*beep.Buffer
	tick int64
	// step func(ms int64) float64

	score float64
	hp    float64
	combo int32

	mods    mania.Mods
	chart   *mania.Chart
	notes   []NoteSprite
	lnotes  []LNSprite // 롱노트 특성상, 2개로 나누는 게 불가피해보임
	endTime int64

	speed float64
	// position     float64
	speedFactors []mode.SpeedFactorPoint
	// sfactorIdx   int

	stage       graphics.ManiaStage
	layout      []types.VKCode
	kbEventChan input.KeyboardEventChannel

	displayScale float64
	hitPosition  float64

	stamps   []timeStamp
	stampIdx int
	// stamp    timeStamp
	logs []keyLog
}

// tools.Stamp를 통해서 구현하려 했다가 element 설정에서 fail
// 각 타입 별로 (float64 등) 만들면 그때 비로소 쓸 수 있을 듯
type timeStamp struct {
	time     int64
	nextTime int64
	position float64
	factor   float64
}

// // decending order
// func SearchTimeStamp(stamps []timeStamp, time int64) int {
// 	return sort.Search(len(stamps), func(i int) bool { return stamps[i].time <= time })
// }

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
	s.endTime = s.chart.EndTime()
	s.displayScale = s.g.ScaleY()
	s.hitPosition = s.g.Settings.HitPosition
	switch {
	default:
		s.speed = s.g.Settings.GeneralSpeed
	}
	// todo: 노트가 언제나 양수 시간에 있다고 상정; 실제로는 노트가 BufferTime보다 뒤에 있을 수 있음
	// initSpeedFactor := mode.SpeedFactorPoint{-BufferTime * Millisecond, 1}
	initSpeedFactor := mode.SpeedFactorPoint{0, 1}
	// s.position = float64(-BufferTime*Millisecond-initSpeedFactor.Time) * initSpeedFactor.Factor
	s.speedFactors = append([]mode.SpeedFactorPoint{initSpeedFactor}, s.chart.TimingPoints.SpeedFactors...)
	{
		stamps := make([]timeStamp, len(s.speedFactors))
		var position float64

		// func SortTimeStamps(stamps []timeStamp) {
		// sort.Slice(stamps, func(i, j int) bool { return stamps[i].time > stamps[j].time })
		// }
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
			stamps[i] = stamp
		}
		s.stamps = stamps
	}
	{
		bg, err := s.chart.Background()
		if err != nil {
			panic(err)
		}
		s.bg, _ = ebiten.NewImageFromImage(bg, ebiten.FilterDefault)
		s.bgop = BackgroundOp(s.g.ScreenSize(), image.Pt(s.bg.Size()))
	}
	// {
	// f, err := os.Open(s.chart.AbsPath(s.chart.AudioFilename))
	// if err != nil {
	// 	panic(err)
	// }
	// var streamer beep.StreamSeekCloser
	// var format beep.Format
	// switch strings.ToLower(filepath.Ext(s.chart.AudioFilename)) {
	// case ".mp3":
	// 	streamer, format, err = mp3.Decode(f)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// _ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2))
	// s.audioPlayer = NewAudioPlayer(streamer, format.SampleRate, s.mods.TimeRate)
	// }
	{
		s.audioPlayer = NewAudioPlayer(s.g.audioContext, s.chart.AbsPath(s.chart.AudioFilename))
	}
	// const BufferTime = 2
	// s.tick = -int64(BufferTime * s.g.MaxTPS())

	// const ApproachDuration = 1500
	// step1ms := float64(s.g.ScreenSize().Y) / ApproachDuration
	// s.step = func(ms int64) float64 { return float64(ms) * step1ms }

	s.stage = s.g.GameSprites.ManiaStages[s.chart.Keys] // for quick access

	s.notes = make([]NoteSprite, len(s.chart.Notes))
	for i, n := range s.chart.Notes {
		var ns NoteSprite
		var sprite graphics.Sprite
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
	var stampIdx int
	stamp := s.stamps[0]
	for i, n := range s.chart.Notes {
		for si := range s.stamps[stampIdx:] {
			if n.Time < s.stamps[stampIdx+si].nextTime && si != 0 {
				stamp = s.stamps[stampIdx+si]
				stampIdx += si
				break
			}
		}
		s.notes[i].position = float64(n.Time-stamp.time)*stamp.factor + stamp.position
	}

	// var current int
	// for _, stamp := range s.stamps {
	// 	// for ; s.chart.Notes[i].Time < stamp.time; i++ {
	// 	for i, n := range s.chart.Notes[current:] {
	// 		// if s.chart.Notes[i].Time >= stamp.time {
	// 		// 	current = i
	// 		// 	break
	// 		// }
	// 		s.notes[i].position = float64(n.Time-stamp.time)*stamp.factor + stamp.position
	// 	}
	// }
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
			ls.bodyop = &ebiten.DrawImageOptions{}
			s.lnotes = append(s.lnotes, ls)
		}
	}

	s.layout = s.g.Settings.KeyLayout[s.chart.Keys]
	{
		var err error
		s.kbEventChan, err = input.NewKeyboardEventChannel()
		if err != nil {
			panic(err)
		}
	}

	var dimness uint8
	switch {
	default:
		dimness = s.g.Settings.GeneralDimness()
	}
	// 별도 함수 없이 여기서 처리하는게 맞을듯
	// 만약 dim 을 바꾸는 입력이 들어왔다면 즉석에서 s.bgop.ColorM.Reset() 날리고 다시 설정.
	s.bgop.ColorM.ChangeHSV(0, 1, float64(dimness)/100)
	return s
}

// 리플레이 구조: 마지막 status 시간, 레이아웃 키state
func (s *SceneMania) Update() error {
	if s.Time() > s.endTime || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if err := s.kbEventChan.Close(); err != nil {
			return err
		}
		s.audioPlayer.Close()
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
		// next note들이 slice에 fetch되어 있는 상태.
		// if done {
		// 	s.notes[i].op.ColorM.ChangeHSV(0, 0, 0.5) // gray
		// }
	}

	now := s.Time()
	// stamp := s.stamps[SearchTimeStamp(s.stamps, now)]
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
		s.lnotes[i].length = n.height() * s.speed * s.displayScale
	}
	// for si, sp := range s.speedFactors[s.sfactorIdx:] { // todo: timing points sorting -> rg-parser에서
	// 	if si == len(s.speedFactors)-1 {
	// 		dy += s.step(now-lastTime) * s.speed * sp.Factor
	// 		s.sfactorIdx = si
	// 	} else {
	// 		nextsp := s.speedFactors[si+1]
	// 		if nextsp.Time > now { // equality condition should be excluded
	// 			s.sfactorIdx = si
	// 			break
	// 		}
	// 		dy += s.step(nextsp.Time-lastTime) * s.speed * sp.Factor
	// 		lastTime = nextsp.Time
	// 	}
	// }

	// for si, sp := range speedFactors[s.sfactorIdx+1:] {
	// 	if sp.Time > s.Time(s.tick) { // equality condition should be excluded
	// 		break
	// 	}
	// 	dy += s.step(sp.Time-lastTime) * s.speed * sp.Factor
	// 	lastTime = sp.Time
	// 	s.sfactorIdx = si
	// }
	// dy *= 1.7

	// s.position += dy
	s.tick++
	return nil
}

// 단순명료하게 전체 노트 그리기; edit 위해서도 좋음
// op를 Draw에서 새로 set
// todo: DrawLN 넣는 순간 엄청 느려짐, 빼도 느림
func (s *SceneMania) Draw(screen *ebiten.Image) {
	// 키 버튼 그리기
	// 스코어, hp, 콤보, 시간 그리기
	// hp는 마스크 이미지를 씌우면 되지 않을까
	// screen.DrawImage(s.bg, s.bgop)
	// screen.DrawImage(s.stage.Fixed.Image(), &ebiten.DrawImageOptions{})
	// for _, n := range s.lnotes {
	// 	n.DrawLN(screen)
	// }
	for _, n := range s.notes {
		// if n.y < -100 || n.y > 1080 {
		// 	continue
		// }
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(n.x, n.y)
		screen.DrawImage(n.i, op)
	}
	// for i, n := range s.notes {
	// 	s.notes[i].op.GeoM.Reset()
	// 	s.notes[i].op.GeoM.Translate(n.x, n.y)
	// 	screen.DrawImage(n.i, s.notes[i].op)
	// }
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTime: %.1fs", ebiten.CurrentFPS(), float64(s.Time())/1000))
}

func (s *SceneMania) Init() {
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", s.chart.MusicName, s.chart.ChartName))
	s.audioPlayer.Play()
}

// always follows audio's time
func (s *SceneMania) Time() int64 {
	return s.audioPlayer.Time().Milliseconds()
}

// func (s *SceneMania) Time() int64 {
// 	return Millisecond * s.tick / int64(s.g.MaxTPS())
// }

// todo: Update에서 그리기
func (s LNSprite) DrawLN(screen *ebiten.Image) {
	_, h := s.i.Size()
	count, remainder := int(s.length)/h, int(s.length)%h+1
	s.bodyop.GeoM.Reset()
	s.bodyop.GeoM.Translate(s.tail.x, s.tail.y)

	firstRect := s.i.Bounds()
	firstRect.Min = image.Pt(0, h-remainder)
	firstImg, _ := ebiten.NewImageFromImage(s.i.SubImage(firstRect), ebiten.FilterDefault)
	screen.DrawImage(firstImg, s.bodyop)
	s.bodyop.GeoM.Translate(0, float64(remainder))

	for c := 0; c < count; c++ {
		screen.DrawImage(s.i, s.bodyop)
		s.bodyop.GeoM.Translate(0, float64(h))
	}
}
