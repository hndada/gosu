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
	buffer *beep.Buffer
	bg     *ebiten.Image
	bgop   *ebiten.DrawImageOptions
	tick   int64
	step   func(ms int64) float64

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
}

// lnhead와 lntail 분리 유지
func (g *Game) NewSceneMania(c *mania.Chart, mods mania.Mods) *SceneMania {
	s := &SceneMania{}
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", c.SongName, c.ChartName))
	{
		bg, err := s.chart.Background()
		if err != nil {
			log.Fatal(err)
		}
		s.bg, _ = ebiten.NewImageFromImage(bg, ebiten.FilterDefault)
		s.bgop = &ebiten.DrawImageOptions{}
		s.bgop.GeoM.Scale(ratio(s.g.ScreenSize(), image.Pt(s.bg.Size()))) // todo: 폭맞춤
		s.bgop.ColorM.ChangeHSV(0, 1, 0.30)
	}
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
	s.applySpeed(s.g.ScrollSpeed)
	s.endTime = s.chart.EndTime()
	return s
}

func (s *SceneMania) Update() error {
	if s.Time(s.tick) > s.endTime {
		s.g.ChangeScene(NewSceneSelect())
	}
	// todo: 플레이 하면서 리플레이 데이터 저장
	// 키보드 입력 채널에서 키 입력 불러오기
	// next note들이 slice에 fetch되어 있는 상태. hit 판정이 나왔을 경우.
	if done {
		s.notes[i].op.ColorM.ChangeHSV(0, 0, 0.5) // gray
	}

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
	screen.DrawImage(s.g.Skin.Mania.Stage.Image, s.g.Skin.Mania.Stage.Op)
	for _, n := range s.notes {
		var img *ebiten.Image
		switch n.noteType {
		case mania.TypeNote:
			img = config.NoteImgs[n.kind]
		case mania.TypeLNHead:
			img = config.LNHeadImgs[n.kind]
		case mania.TypeLNTail:
			img = config.LNTailImgs[n.kind]
		}
		screen.DrawImage(img, &n.op)
	}
	for _, n := range s.lnotes {
		screen.DrawImage(config.LNBodyImgs[n.kind], &n.bodyop)
	}
	// 키 버튼 그리기
	// 스코어, hp, 콤보, 시간 그리기
	// hp는 마스크 이미지를 씌우면 되지 않을까
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTime: %.1fs", ebiten.CurrentFPS(), float64(s.Time(s.tick))/1000))
}

// downcast happens
func (s *SceneMania) Time(tick int64) int64 {
	return Millisecond * tick / int64(s.g.MaxTPS())
}

func (s *SceneMania) Init() {
	speaker.Play(s.buffer.Streamer(0, s.buffer.Len()))
}

func ratio(dst, src image.Point) (float64, float64) {
	return float64(dst.X) / float64(src.X),
		float64(dst.Y) / float64(src.Y)
}
