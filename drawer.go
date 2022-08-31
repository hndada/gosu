package gosu

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
)

// Todo: use Effecter in draws.BaseDrawer
type BackgroundDrawer struct {
	Sprite  draws.Sprite
	Dimness *float64
}

func (d BackgroundDrawer) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.ChangeHSV(0, 1, *d.Dimness)
	d.Sprite.Draw(screen, op)
	// op := d.Sprite.Op()
	// screen.DrawImage(d.Sprite.I, op)
}

func NewScoreDrawer() draws.NumberDrawer {
	return draws.NumberDrawer{
		Sprites:     ScoreSprites,
		SignSprites: SignSprites,
		DigitWidth:  ScoreSprites[0].W(),
		ZeroFill:    1,
		Origin:      ScoreSprites[0].Origin(),
	}
}

type Direction int

const (
	Upward   Direction = iota // e.g., Rhythm games using feet.
	Downward                  // e.g., Piano mode.
	Leftward                  // e.g., Drum mode.
	Rightward
)

// TPS should be multiple of 1000, since only one speed value
// goes passed per Update, while unit of TransPoint's time is 1ms.
var TimeStep float64 = 1000 / float64(TPS)

// LaneDrawer's cursor position should be consistent with ScenePlay.
type baseLaneDrawer struct {
	direction Direction

	sprites  []draws.Sprite
	margin   float64 // Half of max sizes of sprites.
	bodyLoss float64 // Head/2 + Tail/2

	hitPosition float64
	maxPosition float64
	minPosition float64

	// beatSpeed  float64 // BPM(ratio) * BeatLengthScale
	speedScale float64

	marker draws.Effecter
}

func newBaseLaneDrawer(
	direction Direction,
	sprites []draws.Sprite,
	hitPosition float64,
	speedScale float64,
	marker draws.Effecter,
) (d baseLaneDrawer) {
	d.direction = direction
	d.sprites = sprites
	d.hitPosition = hitPosition
	d.speedScale = speedScale
	d.marker = marker
	var xMax, yMax float64
	for _, s := range sprites {
		if xMax < s.X() {
			xMax = s.X()
		}
		if yMax < s.Y() {
			yMax = s.Y()
		}
	}
	if len(d.sprites) >= 4 { // Long body enabled.
		switch d.direction {
		case Downward, Upward:
			d.margin = yMax / 2
			d.bodyLoss = sprites[Head].H()/2 + sprites[Tail].H()/2
		case Leftward, Rightward:
			d.margin = xMax / 2
			d.bodyLoss = sprites[Head].W()/2 + sprites[Tail].W()/2
		}
	}
	switch d.direction {
	case Upward:
		d.maxPosition = screenSizeY - d.hitPosition
		d.minPosition = -d.hitPosition
	case Downward:
		d.maxPosition = d.hitPosition
		d.minPosition = -screenSizeY + d.hitPosition
	case Leftward:
		d.maxPosition = screenSizeX - d.hitPosition
		d.minPosition = -d.hitPosition
	case Rightward:
		d.maxPosition = d.hitPosition
		d.minPosition = -screenSizeX + d.hitPosition
	}
	return
}

type LaneObject struct {
	Type     int
	Position float64
	Speed    float64 // Not used in FixedLaneDrawer.
	Next     *LaneObject
	Prev     *LaneObject
	Marked   *bool
	// draws.Effecter
	// Effecter draws.Effecter
}

// DrawLongBody finds sub-image of Body sprite corresponding to current exposed long body
// and scale the sub-image to (exposed length) / (sub-image length).

// Tail's Position is always larger than Head's.
// In other word, Head is always nearer than Tail.
// Start is Head's, and End is Tail's.
func (d baseLaneDrawer) DrawLongBody(screen *ebiten.Image, head *LaneObject) { //headSrc any) {
	// head := headSrc.(*LaneObject)
	tail := head.Next //.(*LaneObject)
	length := tail.Position - head.Position
	length -= -d.bodyLoss
	startPosition := head.Position
	if startPosition < d.minPosition {
		startPosition = d.minPosition
	}
	endPosition := tail.Position
	if endPosition > d.maxPosition {
		endPosition = d.maxPosition
	}
	srcSprite := d.sprites[Body]
	ratio := length / srcSprite.H()
	srcStart := math.Floor((startPosition - head.Position) / ratio)
	srcEnd := math.Ceil((endPosition - head.Position) / ratio)
	op := &ebiten.DrawImageOptions{}
	if d.marker != nil {
		d.marker(op, *head.Marked)
	}
	switch d.direction {
	case Upward, Downward:
		if d.direction == Downward {
			ratio *= -1
			srcStart *= -1
		}
		srcRect := image.Rect(0, int(srcStart), int(srcSprite.W()), int(srcEnd))
		sprite := srcSprite.SubSprite(srcRect)
		op.GeoM.Scale(1, ratio)
		op.GeoM.Translate(0, srcStart)
		sprite.Draw(screen, op)
	case Leftward, Rightward:
		if d.direction == Rightward {
			ratio *= -1
			srcStart *= -1
		}
		srcRect := image.Rect(int(srcStart), 0, int(srcEnd), int(srcSprite.H()))
		sprite := srcSprite.SubSprite(srcRect)
		op.GeoM.Scale(ratio, 1)
		op.GeoM.Translate(srcStart, 0)
		sprite.Draw(screen, op)
	}
}

// In FixedLaneDrawer, lane itself moves.
// Hence all notes move same amount.
// Piano mode uses FixedLaneDrawer.
type FixedLaneDrawer struct {
	baseLaneDrawer
	cursor   float64
	farthest *LaneObject
	nearest  *LaneObject
}

func NewFixedLaneDrawer(
	direction Direction,
	sprites []draws.Sprite,
	hitPosition float64,
	// beatSpeed,  // Speed calculation is each mode's task.
	speedScale float64,
	marker draws.Effecter,

	startTime int64,
	tp *TransPoint,
	leading *LaneObject,
) (d FixedLaneDrawer) {
	d.baseLaneDrawer = newBaseLaneDrawer(
		direction, sprites, hitPosition, speedScale, marker)
	d.cursor = tp.Position
	d.cursor -= float64(tp.Time-startTime) * tp.BPM
	d.SetSpeedScale(speedScale)
	d.farthest = leading
	d.nearest = leading
	return
}

// Need to re-calculate cursor's position when speed scale changes.
func (d *FixedLaneDrawer) SetSpeedScale(speedScale float64) {
	d.cursor /= d.speedScale
	d.speedScale = speedScale
	d.cursor *= d.speedScale
}

// Speed: BPM (or BPM ratio) * BeatLengthScale * SpeedScale.
func (d *FixedLaneDrawer) Update(beatSpeed, speedScale float64) {
	speed := beatSpeed * d.speedScale
	d.cursor += speed * TimeStep
	for d.farthest.Position-d.cursor <= d.maxPosition {
		d.farthest = d.farthest.Next //.(*LaneObject)
	}
	for d.nearest.Position-d.cursor <= d.minPosition {
		d.nearest = d.nearest.Next //.(*LaneObject)
	}
	if speedScale != d.speedScale {
		d.SetSpeedScale(speedScale)
	}
}

// Draw from farthest to nearest.
// So that nearer notes are exposed when overlapped with farther notes.
func (d FixedLaneDrawer) Draw(screen *ebiten.Image) {
	obj := d.farthest
	for ; obj != d.nearest; obj = d.farthest.Prev { //.(*LaneObject) {
		sprite := d.sprites[obj.Type]
		offset := obj.Position - d.cursor
		switch d.direction {
		case Downward, Upward:
			sprite.Move(0, offset)
		case Leftward, Rightward:
			sprite.Move(offset, 0)
		}
		op := &ebiten.DrawImageOptions{}
		if d.marker != nil {
			d.marker(op, *obj.Marked)
		}
		sprite.Draw(screen, op)
		if obj.Type == Head {
			d.DrawLongBody(screen, obj)
		}
	}
	if obj.Type == Tail {
		d.DrawLongBody(screen, obj.Prev)
	}
}

// Todo: set draw order for performance?
type FloatLaneDrawer struct {
	baseLaneDrawer
	Objects []*LaneObject
}

// Speed calculation is each mode's task.
func NewFloatLaneDrawer(
	direction Direction,
	sprites []draws.Sprite,
	hitPosition float64,
	speedScale float64,
	marker draws.Effecter,

	objs []*LaneObject,
) (d FloatLaneDrawer) {
	d.baseLaneDrawer = newBaseLaneDrawer(
		direction, sprites, hitPosition, speedScale, marker)
	// Reverse objects slice.
	for i, j := 0, len(objs)-1; i < j; i, j = i+1, j-1 {
		objs[i], objs[j] = objs[j], objs[i]
	}
	d.Objects = objs
	d.SetSpeedScale(speedScale)
	return
}

// Need to re-calculate cursor's position when speed scale changes.
func (d *FloatLaneDrawer) SetSpeedScale(speedScale float64) {
	for i := range d.Objects {
		d.Objects[i].Position /= d.speedScale
		d.Objects[i].Position *= speedScale
	}
	d.speedScale = speedScale
}

// Speed: BPM (or BPM ratio) * BeatLengthScale * SpeedScale.
func (d *FloatLaneDrawer) Update(beatSpeed, speedScale float64) {
	for i, obj := range d.Objects {
		speed := obj.Speed * d.speedScale
		d.Objects[i].Position -= speed * TimeStep
	}
	if speedScale != d.speedScale {
		d.SetSpeedScale(speedScale)
	}
}

// Draw from farthest to nearest.
// So that nearer notes are exposed when overlapped with farther notes.
func (d FloatLaneDrawer) Draw(screen *ebiten.Image) {
	heads := make([]*LaneObject, 0)
	for _, obj := range d.Objects {
		if obj.Type == Head {
			head := obj
			tail := obj.Next //.(LaneObject)
			if head.Position <= d.maxPosition &&
				tail.Position >= d.minPosition {
				heads = append(heads, obj)
			}
		}
		if obj.Position > d.maxPosition ||
			obj.Position < d.minPosition {
			continue
		}
		sprite := d.sprites[obj.Type]
		offset := obj.Position
		switch d.direction {
		case Downward, Upward:
			sprite.Move(0, offset)
		case Leftward, Rightward:
			sprite.Move(offset, 0)
		}
		op := &ebiten.DrawImageOptions{}
		if d.marker != nil {
			d.marker(op, *obj.Marked)
		}
		sprite.Draw(screen, op)
	}
	for _, head := range heads {
		d.DrawLongBody(screen, head)
	}
}

// type BaseLaneDrawer struct {
// 	Tick int
// 	// Farthest   *Note
// 	// Nearest    *Note
// 	Cursor     float64
// 	HitPostion float64
// 	Speed      float64 // BPM (or BPM ratio) * BeatLengthScale
// 	Direction
// 	maxPosition float64
// 	minPosition float64
// }

// type NoteLaneDrawer struct {
// 	BaseLaneDrawer
// 	Sprites  [4]draws.Sprite
// 	Farthest *Note
// 	Nearest  *Note
// 	margin   float64 // Half of max sizes of sprites.
// 	bodyLoss float64 // Head/2 + Tail/2
// }
// type BarDrawer struct {
// 	BaseLaneDrawer
// 	Sprite   draws.Sprite
// 	Bars     []Bar
// 	Farthest int
// 	Nearest  int
// 	count    int
// }

// // Update should use existing speed, not the new one.
//
//	func (d *BarDrawer) Update(speed float64) {
//		d.Cursor += speed * TimeStep
//		var a, b int
//		// var boundFarIn, boundNearOut float64 // Bounds for farthest, nearest each.
//		for d.Bars[d.Farthest].Position-d.Cursor <= d.maxPosition {
//			d.Farthest++
//			a++
//		}
//		for d.Bars[d.Nearest].Position-d.Cursor <= d.minPosition {
//			d.Nearest++
//			b++
//		}
//		d.Speed = speed
//		if d.count%1000 == 0 {
//			fmt.Println(d.maxPosition, d.minPosition, a, b)
//			fmt.Println(d.Farthest, d.Bars[d.Farthest])
//			fmt.Println(d.Nearest, d.Bars[d.Nearest])
//		}
//		d.count++
//		if d.count > 100000 {
//			os.Exit(1)
//		}
//	}
// func (d BarDrawer) Draw(screen *ebiten.Image) {
// 	for i := d.Farthest; i >= d.Nearest; i-- {
// 		op := &ebiten.DrawImageOptions{}
// 		offset := d.Bars[i].Position - d.Cursor
// 		switch d.Direction {
// 		case Downward, Upward:
// 			op.GeoM.Translate(0, offset)
// 		case Leftward, Rightward:
// 			op.GeoM.Translate(offset, 0)
// 		}
// 		d.Sprite.Draw(screen, op)
// 	}
// }

// // NoteLaneDrawer's tick should be consistent with ScenePlay.
//
//	type NoteLaneDrawer struct {
//		Tick       int
//		Sprites    [4]draws.Sprite //  map[NoteType]draws.Sprite // []draws.Sprite
//		Farthest   *Note
//		Nearest    *Note
//		Cursor     float64
//		HitPostion float64
//		Speed      float64 // BPM (or BPM ratio) * BeatLengthScale
//		Direction
//		// Sizes      map[NoteType]float64 // Cache for Sprites' sizes. // Todo: Sizes -> halfSizes
//		// MaxSize    float64              // Either max width / height. // Todo: remove
//		margin   float64 // Half of max sizes of sprites.
//		bodyLoss float64 // Head/2 + Tail/2
//		// boundFarIn   float64 // Bound for Farthest note being fetched.
//		// boundNearOut float64 // Bound for Nearest note being flushed.
//		maxPosition float64
//		minPosition float64
//	}

// Update should use existing speed, not the new one.
// func (d *LaneDrawer) Update(speed float64) {
// 	d.cursor += speed * TimeStep
// 	// var boundFarIn, boundNearOut float64 // Bounds for farthest, nearest each.
// 	for d.farthest.Position-d.cursor <= d.maxPosition {
// 		d.farthest = d.farthest.Next
// 	}
// 	for d.nearest.Position-d.cursor <= d.minPosition {
// 		d.nearest = d.nearest.Next
// 	}
// 	d.Speed = speed
// }
// for d.ScreenPosition(d.Farthest) >= d.boundFarIn {
// 	d.Farthest = d.Farthest.Next
// }
// for d.ScreenPosition(d.Nearest) >= d.boundNearOut {
// 	d.Nearest = d.Nearest.Next
// }
// case Upward:
// 	for d.ScreenPosition(d.Farthest)-d.margin >= screenSizeY {
// 		d.Farthest = d.Farthest.Next
// 	}
// 	for d.ScreenPosition(d.Nearest)+d.margin >= 0 {
// 		d.Nearest = d.Nearest.Next
// 	}
// case Leftward:
// 	for d.ScreenPosition(d.Farthest)-d.margin >= screenSizeX {
// 		d.Farthest = d.Farthest.Next
// 	}
// 	for d.ScreenPosition(d.Nearest)+d.margin >= 0 {
// 		d.Nearest = d.Nearest.Next
// 	}
// case Rightward:
// 	for d.ScreenPosition(d.Farthest)+d.margin >= 0 {
// 		d.Farthest = d.Farthest.Next
// 	}
// 	for d.ScreenPosition(d.Nearest)-d.margin >= screenSizeX {
// 		d.Nearest = d.Nearest.Next
// 	}
// }
// switch d.Direction {
// case Downward:
// 	for d.Note.Position-d.MaxSize/2 >= screenSizeY {
// 		d.Note = d.Note.Next
// 	}
// case Leftward:
// 	for d.Note.Position+d.MaxSize/2 < 0 {
// 		d.Note = d.Note.Next
// 	}
// }

// func (d NoteLaneDrawer) ScreenPosition(n *Note) float64 {
// 	pos := n.Position - d.Cursor // Relative position of note.
// 	switch d.Direction {
// 	case Downward, Rightward:
// 		pos *= -1
// 	case Upward, Leftward:
// 		pos *= 1
// 	}
// 	return d.HitPostion + pos
// }

// type BarLineDrawer struct {
// 	Times      []int64
// 	Cursor     int     // Index of closest bar line.
// 	Offset     float64 // Bar line is drawn at bottom, not at the center.
// 	Sprite     draws.Sprite
// 	Horizontal bool
// }
// type LaneDrawer struct {
// 	Sprite draws.Sprite
// 	Object *any
// 	Bound  func() bool
// }
// type LongBodyDrawer struct {
// 	Sprite draws.Sprite
// 	Object *any
// 	Bound  func() bool
// }

// func (d *BarLineDrawer) Update(position func(time int64) float64) {
// 	bound := screenSizeY
// 	if d.Horizontal {
// 		bound = screenSizeX
// 	}
// 	t := d.Times[d.Cursor]
// 	// Bar line and Hint are anchored at the bottom.
// 	for d.Cursor < len(d.Times)-1 &&
// 		int(position(t)+d.Offset) >= bound {
// 		d.Cursor++
// 		t = d.Times[d.Cursor]
// 	}
// }
// func (d BarLineDrawer) Draw(screen *ebiten.Image, position func(time int64) float64) {
// 	for _, t := range d.Times[d.Cursor:] {
// 		sprite := d.Sprite
// 		sprite.Y = position(t) + d.Offset
// 		if sprite.Y < 0 {
// 			break
// 		}
// 		sprite.Draw(screen)
// 	}
// }

// func (d *ScoreDrawer) Update(score float64) {
// 	d.DelayedScore.Set(score)
// 	d.DelayedScore.Update()
// }

// func (d ScoreDrawer) Draw(screen *ebiten.Image) {
// 	var wsum int
// 	vs := make([]int, 0)
// 	for v := int(math.Ceil(d.DelayedScore.Delayed)); v > 0; v /= 10 {
// 		vs = append(vs, v%10) // Little endian
// 		// wsum += int(d.Sprites[v%10].W)
// 		wsum += int(d.Sprites[0].W)
// 	}
// 	if len(vs) == 0 {
// 		vs = append(vs, 0) // Little endian
// 		wsum += int(d.Sprites[0].W)
// 	}
// 	x := float64(screenSizeX) - d.Sprites[0].W/2
// 	for _, v := range vs {
// 		// x -= d.Sprites[v].W
// 		x -= d.Sprites[0].W
// 		sprite := d.Sprites[v]
// 		sprite.X = x + (d.Sprites[0].W - sprite.W/2)
// 		sprite.Draw(screen)
// 	}
// }

//	type ScoreDrawer struct {
//		DelayedScore ctrl.Delayed
//		Sprites      []draws.Sprite
//	}
//
// ScoreDrawer.Update(int(math.Ceil(delayedScore)))
