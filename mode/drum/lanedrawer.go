package drum

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Todo: set draw order for performance?
type FloatLaneDrawer struct {
	baseLaneDrawer
	objects []LaneObject
}

// Speed calculation is each mode's task.
func NewFloatLaneDrawer(
	direction Direction,
	hitPosition float64,
	maxPosition float64,
	minPosition float64,
	margin float64,
	bodyLoss float64,

	beatSpeed float64, // BPM(ratio) * BeatLengthScale
	speedScale float64,
	marker func(op *ebiten.DrawImageOptions, marked bool), // draws.Effecter

	objs []LaneObject,
) (d FloatLaneDrawer) {
	d.baseLaneDrawer = newBaseLaneDrawer(
		direction,
		hitPosition,
		maxPosition,
		minPosition,
		margin,
		bodyLoss,

		beatSpeed, // BPM(ratio) * BeatLengthScale
		speedScale,
		marker,
	)
	// Reverse objects slice.
	for i, j := 0, len(objs)-1; i < j; i, j = i+1, j-1 {
		objs[i], objs[j] = objs[j], objs[i]
	}
	d.objects = objs
	d.SetSpeedScale(speedScale)
	return
}

// Need to re-calculate cursor's position when speed scale changes.
func (d *FloatLaneDrawer) SetSpeedScale(speedScale float64) {
	for i := range d.objects {
		pos := d.objects[i].Position()
		pos *= speedScale / d.speedScale
		d.objects[i].SetPosition(pos)
	}
	d.speedScale = speedScale
}

// Speed: BPM (or BPM ratio) * BeatLengthScale * SpeedScale.
func (d *FloatLaneDrawer) Update(beatSpeed, speedScale float64) {
	for i, obj := range d.objects {
		speed := obj.Speed() * d.speedScale
		pos := d.objects[i].Position()
		pos -= speed * TimeStep
		d.objects[i].SetPosition(pos)
	}
	if speedScale != d.speedScale {
		d.SetSpeedScale(speedScale)
	}
}

// Draw from farthest to nearest to make nearer notes exposed
// when being overlapped with farther notes.
func (d FloatLaneDrawer) Draw(screen *ebiten.Image) {
	heads := make([]LaneObject, 0)
	for _, obj := range d.objects {
		if obj.IsHead() {
			head := obj
			tail := obj.Next()
			if head.Position() <= d.maxPosition+d.margin &&
				tail.Position() >= d.minPosition-d.margin {
				heads = append(heads, obj)
			}
		}
		if obj.Position() > d.maxPosition+d.margin ||
			obj.Position() < d.minPosition-d.margin {
			continue
		}
		sprite := obj.Sprite()
		offset := obj.Position()
		switch d.direction {
		case Downward, Upward:
			sprite.Move(0, offset)
		case Leftward, Rightward:
			sprite.Move(offset, 0)
		}
		op := &ebiten.DrawImageOptions{}
		if d.marker != nil {
			d.marker(op, obj.Marked())
		}
		sprite.Draw(screen, op)
	}
	for _, head := range heads {
		d.DrawLongBody(screen, head)
	}
}

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
