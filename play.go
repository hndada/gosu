package gosu

import (
	"fmt"
	"image/color"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ScenePlay: struct, PlayScene: function
type ScenePlay struct {
	Tick int

	Chart       *Chart
	PlayNotes   []*PlayNote
	KeySettings []ebiten.Key // Todo: separate ebiten

	FetchPressed func() []bool
	Pressed      []bool
	LastPressed  []bool

	StagedNotes    []*PlayNote
	Combo          int
	Karma          float64
	KarmaSum       float64
	JudgmentCounts []int

	// In dev
	// ReplayMode   bool
	// ReplayStates []ReplayState
	// ReplayCursor int

	Speed float64
	TransPoint
	NoteSprites []Sprite
	BodySprites []Sprite
	HeadSprites []Sprite
	TailSprites []Sprite

	MusicFile   io.ReadSeekCloser
	MusicPlayer *audio.Player

	Background      Sprite
	FieldSprite     Sprite
	ComboSprites    []Sprite
	ScoreSprites    []Sprite
	JudgmentSprites []Sprite
	ClearSprite     Sprite
	HintSprite      Sprite

	Judgment          Judgment
	JudgmentCountdown int
}

func TickToMsec(tick int) int64 { return int64(1000 * float64(tick) / float64(MaxTPS)) }
func MsecToTick(msec int64) int { return int(msec) * MaxTPS }
func NewScenePlay(c *Chart, cpath string) *ScenePlay {
	s := new(ScenePlay)
	s.Tick = -2 * MaxTPS // Put 2 seconds of waiting
	s.Chart = c
	s.PlayNotes, s.StagedNotes = NewPlayNotes(c) // Todo: add Mods to input param
	// s.KeySettings = KeySettings[s.Chart.Parameter.KeyCount]
	s.JudgmentCounts = make([]int, 5)
	s.LastPressed = make([]bool, c.KeyCount)
	s.Pressed = make([]bool, c.KeyCount)
	s.Karma = 1
	s.TransPoint = TransPoint{
		s.Chart.SpeedFactors[0],
		s.Chart.Tempos[0],
		s.Chart.Volumes[0],
		s.Chart.Effects[0],
	}
	s.NoteSprites = make([]Sprite, c.KeyCount)
	s.BodySprites = make([]Sprite, c.KeyCount)
	var wsum int
	for k, kind := range NoteKindsMap[c.KeyCount] {
		w := int(NoteWidths[c.KeyCount][int(kind)] * Scale()) // w should be integer, since it is a play note's width.
		var fpath string
		fpath = "skin/note/" + fmt.Sprintf("n%d.png", []int{1, 2, 3, 3}) // Todo: 4th note image
		s.NoteSprites[k] = Sprite{
			I: NewImage(fpath),
			W: float64(w),
			H: NoteHeigth * Scale(),
		}
		fpath = "skin/note/" + fmt.Sprintf("l%d.png", []int{1, 2, 3, 3}) // Todo: 4th note image
		s.BodySprites[k] = Sprite{
			I: NewImage(fpath),
			W: float64(w),
			H: NoteHeigth * Scale(), // Long note body does not have to be scaled though.
		}
		wsum += w
	}
	s.HeadSprites = make([]Sprite, c.KeyCount)
	s.TailSprites = make([]Sprite, c.KeyCount)
	copy(s.HeadSprites, s.NoteSprites)
	copy(s.TailSprites, s.NoteSprites)

	// Todo: Scratch should be excluded to width sum.
	x := (ScreenSizeX - wsum) / 2 // x should be integer as well as w
	for k, kind := range NoteKindsMap[c.KeyCount] {
		s.NoteSprites[k].X = float64(x)
		x += int(NoteWidths[c.KeyCount][kind] * Scale())
	}
	f, err := os.Open(c.MusicPath(cpath))
	if err != nil {
		panic(err)
	}
	s.MusicPlayer, err = Context.NewPlayer(f)
	if err != nil {
		panic(err)
	}

	s.Background = Sprite{
		I: NewImage(c.BgPath(cpath)),
		W: float64(ScreenSizeX),
		H: float64(ScreenSizeY),
	}
	{
		w := wsum
		h := ScreenSizeY
		x := (ScreenSizeX - w) / 2
		i := ebiten.NewImage(w, h)
		i.Fill(color.RGBA{0, 0, 0, uint8(255 * FieldDark)})
		s.FieldSprite = Sprite{i, float64(w), float64(h), float64(x), 0}
	}
	s.FieldSprite = Sprite{}
	s.ComboSprites = make([]Sprite, 10)
	for i := 0; i < 10; i++ {
		sp := Sprite{
			I: NewImage(fmt.Sprintf("skin/combo/%d.png", i)),
			W: ComboWidth * Scale(),
		}
		sp.H = float64(sp.I.Bounds().Dy()) * (sp.W / float64(sp.I.Bounds().Dx()))
		sp.Y = ComboPosition - sp.H/2
		s.ComboSprites[i] = sp
	}
	s.ScoreSprites = make([]Sprite, 10)
	for i := 0; i < 10; i++ {
		sp := Sprite{
			I: NewImage(fmt.Sprintf("skin/score/%d.png", i)),
			W: ScoreWidth * Scale(),
		}
		sp.H = float64(sp.I.Bounds().Dy()) * (sp.W / float64(sp.I.Bounds().Dx()))
		s.ComboSprites[i] = sp
	}
	s.JudgmentSprites = make([]Sprite, 5)
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		sp := Sprite{
			I: NewImage(fmt.Sprintf("skin/judgment/%s.png", name)),
			W: JudgmentWidth * Scale(),
		}
		sp.H = float64(sp.I.Bounds().Dy()) * (sp.W / float64(sp.I.Bounds().Dx()))
		sp.X = (float64(ScreenSizeX) - sp.W) / 2
		sp.Y = JudgePosition*Scale() - sp.H/2
		s.JudgmentSprites[i] = sp
	}
	{
		sp := Sprite{
			I: NewImage("skin/play/clear.png"),
		}
		sp.W = float64(sp.I.Bounds().Dx())
		sp.H = float64(sp.I.Bounds().Dy())
		sp.X = (float64(ScreenSizeX) - sp.W) / 2
		sp.Y = (float64(ScreenSizeY) - sp.H) / 2
		s.ClearSprite = sp
	}
	{
		sp := Sprite{
			I: NewImage("skin/play/hint.png"),
		}
		sp.W = float64(wsum)
		sp.H = HintHeight * Scale()
		sp.X = (float64(ScreenSizeX) - sp.W) / 2
		sp.Y = HintPosition*Scale() - sp.H/2
		s.HintSprite = sp
	}
	return s
}

// TPS affects only on Update(), not on Draw()
func (s *ScenePlay) Update() {
	if s.IsFinished() {
		if s.MusicPlayer != nil {
			s.MusicFile.Close()
			s.MusicPlayer = nil // Todo: need a test
		}
		return
	}
	if s.Tick == 0 {
		s.MusicPlayer.Play()
	}
	for s.Time() < s.SpeedFactor.Next.Time {
		s.SpeedFactor = s.SpeedFactor.Next
	}
	for s.Time() < s.Tempo.Next.Time {
		s.Tempo = s.Tempo.Next
	}
	for s.Time() < s.Volume.Next.Time {
		s.Volume = s.Volume.Next
	}
	for s.Time() < s.Effect.Next.Time {
		s.Effect = s.Effect.Next
	}
	s.Pressed = s.FetchPressed()
	for k, n := range s.StagedNotes {
		if n == nil {
			continue
		}
		if n.Type != Tail && s.KeyAction(k) == Hit {
			n.PlaySE()
		}
		td := n.Time - s.Time() // Time difference; negative values means late hit
		if n.Scored {
			if n.Type != Tail {
				panic("non-tail note has not flushed")
			}
			if td < Miss.Window { // Keep Tail being staged until nearly ends
				s.StagedNotes[n.Key] = n.Next
			}
			continue
		}
		var worst Judgment
		if j := Verdict(n.Type, s.KeyAction(n.Key), td); j.Window != 0 {
			s.Score(n, j)
			if worst.Window < j.Window {
				worst = j
			}
		}
		if s.Judgment.Window < worst.Window {
			s.Judgment = worst
			s.JudgmentCountdown = MsecToTick(1000)
		} else {
			s.JudgmentCountdown--
		}
	}
	s.Tick++
}
func (s ScenePlay) Time() int64 {
	return int64(float64(s.Tick) / float64(MaxTPS) * 1000)
}

func (s ScenePlay) IsFinished() bool {
	return s.Time() > 3000+s.PlayNotes[len(s.PlayNotes)-1].Time
}

func (s *ScenePlay) Draw(screen *ebiten.Image) {
	{
		op := s.Background.Op()
		op.ColorM.ChangeHSV(0, 1, BgDimness)
		screen.DrawImage(s.Background.I, op)
	}
	screen.DrawImage(s.FieldSprite.I, s.FieldSprite.Op())
	if s.IsFinished() {
		screen.DrawImage(s.ClearSprite.I, s.ClearSprite.Op())
	} else {
		s.DrawNotes(screen)
		if s.Combo > 0 {
			s.DrawCombo(screen)
		}
		if s.JudgmentCountdown > 0 { // Draw the same judgment for a while.
			for i, j := range Judgments {
				if j.Window == s.Judgment.Window {
					sp := s.JudgmentSprites[i]
					screen.DrawImage(sp.I, sp.Op())
					break
				}
			}
		}
		s.DrawScore(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"CurrentFPS: %.2f\nCurrentTPS: %.2f\nTime: %.3fs\n"+
			"Score: %.0f\nKarma: %.2f\nCombo: %d\n"+
			"judge: %v", ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(s.Time())/1000,
		s.CurrentScore(), s.Karma, s.Combo,
		s.JudgmentCounts))
}

// DrawCombo supposes each number image has different size.
// Wait, we loaded number image with adjusting size.
func (s *ScenePlay) DrawCombo(screen *ebiten.Image) {
	gap := ComboGap * Scale()
	var wsum int
	vs := make([]int, 0)
	for v := s.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian
		// vs = append([]int{v % 10}, vs...) // Big endian
		wsum += int(s.ComboSprites[v%10].W - gap)
	}
	wsum += int(gap)
	x := float64(ScreenSizeX+wsum) / 2
	for _, v := range vs {
		x -= s.ComboSprites[v].W + gap
		y := ComboPosition*Scale() - s.ComboSprites[v].H/2
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(s.ComboSprites[v].I, op)
	}
}

func (s *ScenePlay) DrawScore(screen *ebiten.Image) {
	var wsum int
	vs := make([]int, 0)
	for v := int(s.CurrentScore()); v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian
		// vs = append([]int{v % 10}, vs...) // Big endian
		wsum += int(s.ComboSprites[v%10].W)
	}
	x := float64(ScreenSizeX)
	for _, v := range vs {
		x -= s.ScoreSprites[v].W // ScoreWidth * Scale()
		y := 0.0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(s.ScoreSprites[v].I, op)
	}
}
func (s *ScenePlay) KeyAction(k int) KeyAction {
	return CurrentKeyAction(s.LastPressed[k], s.Pressed[k])
}
