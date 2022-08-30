package gosu

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osu"
)

type NoteType int

const (
	Normal NoteType = iota
	Head
	Tail
	Body
)

// Strategy of Piano mode
// Calculate position of each note in advance
// Parameter: SpeedBase, BPM Ratio, BeatScale
// Calculate current HitPosition only.
// For other notes, just calculate the difference between HitPosition.
type Note struct {
	Type         NoteType
	Time         int64
	Time2        int64
	SampleName   string // SampleFilename
	SampleVolume float64

	Next *Note
	Prev *Note // For accessing to Head from Tail.
	// NextTail *Note // For drawing long body faster.
	Position float64 // Scaled x or y value.
	Marked   bool

	Key int // Not used in Drum mode.
}

func NewNote(f any, mode, subMode int) []*Note {
	ns := make([]*Note, 0, 2)
	switch f := f.(type) {
	case osu.HitObject:
		n := &Note{
			Type:  Normal,
			Time:  int64(f.Time),
			Time2: int64(f.Time),
			// Key:          f.Column(keyCount),
			SampleName:   f.HitSample.Filename,
			SampleVolume: float64(f.HitSample.Volume) / 100,
		}
		if mode == ModeTypePiano4 || mode == ModeTypePiano7 {
			n.Key = f.Column(subMode)
		}
		if f.NoteType&osu.ComboMask == osu.HitTypeHoldNote {
			n.Type = Head
			n.Time2 = int64(f.EndTime)
			n2 := &Note{
				Type:  Tail,
				Time:  n.Time2,
				Time2: n.Time,
				// Key:   n.Key,
				// Tail has no sample sound.
			}
			if mode == ModeTypePiano4 || mode == ModeTypePiano7 {
				n2.Key = f.Column(subMode)
			}
			ns = append(ns, n, n2)
		} else {
			ns = append(ns, n)
		}
	}
	return ns
}

// TimeStep is expected to be integer.
// TPS should be either multiple or divisor of 1000.
var TimeStep float64 = 1000 / float64(ebiten.MaxTPS())

type BaseLaneDrawer struct {
	Tick int
	// Farthest   *Note
	// Nearest    *Note
	Cursor     float64
	HitPostion float64
	Speed      float64 // BPM (or BPM ratio) * BeatScale
	Direction
	maxPosition float64
	minPosition float64
}

// NoteLaneDrawer's tick should be consistent with ScenePlay.
type NoteLaneDrawer struct {
	BaseLaneDrawer
	Sprites  [4]draws.Sprite
	Farthest *Note
	Nearest  *Note
	margin   float64 // Half of max sizes of sprites.
	bodyLoss float64 // Head/2 + Tail/2
}
type BarDrawer struct {
	BaseLaneDrawer
	Sprite   draws.Sprite
	Bars     []Bar
	Farthest int
	Nearest  int
}

// Update should use existing speed, not the new one.
func (d *BarDrawer) Update(speed float64) {
	d.Cursor += speed * TimeStep
	fmt.Println(d.Cursor)
	// var boundFarIn, boundNearOut float64 // Bounds for farthest, nearest each.
	for d.Bars[d.Farthest].Position-d.Cursor <= d.maxPosition {
		d.Farthest++
	}
	for d.Bars[d.Nearest].Position-d.Cursor <= d.minPosition {
		d.Nearest++
	}
	d.Speed = speed
}
func (d BarDrawer) Draw(screen *ebiten.Image) {
	for i := d.Farthest; i >= d.Nearest; i-- {
		op := &ebiten.DrawImageOptions{}
		offset := d.Bars[i].Position - d.Cursor
		switch d.Direction {
		case Downward, Upward:
			op.GeoM.Translate(0, offset)
		case Leftward, Rightward:
			op.GeoM.Translate(offset, 0)
		}
		d.Sprite.Draw(screen, op)
	}
}

// // NoteLaneDrawer's tick should be consistent with ScenePlay.
//
//	type NoteLaneDrawer struct {
//		Tick       int
//		Sprites    [4]draws.Sprite //  map[NoteType]draws.Sprite // []draws.Sprite
//		Farthest   *Note
//		Nearest    *Note
//		Cursor     float64
//		HitPostion float64
//		Speed      float64 // BPM (or BPM ratio) * BeatScale
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
type Direction int

const (
	Upward   Direction = iota // e.g., Rhythm games using feet.
	Downward                  // e.g., Piano mode.
	Leftward                  // e.g., Drum mode.
	Rightward
)

func (d *BaseLaneDrawer) SetDirection(direction Direction) {
	d.Direction = direction
	switch d.Direction {
	case Upward:
		d.maxPosition = screenSizeY - d.HitPostion
		d.minPosition = -d.HitPostion
	case Downward:
		d.maxPosition = d.HitPostion
		d.minPosition = -screenSizeY + d.HitPostion
	case Leftward:
		d.maxPosition = screenSizeX - d.HitPostion
		d.minPosition = -d.HitPostion
	case Rightward:
		d.maxPosition = d.HitPostion
		d.minPosition = -screenSizeX + d.HitPostion
	}
	// fmt.Printf("%+v\n", d)
}

// func NewBaseLaneDrawer(sprites []draws.Sprite, direction Direction) (d BaseLaneDrawer) {
// 	d.Sprites = sprites
// 	d.Direction = direction

// }

// [4]draws.Sprite{Note, Head, Tail, Body}
// , direction Direction
func NewNoteLaneDrawer(sprites [4]draws.Sprite) (d NoteLaneDrawer) {
	// d.BaseLaneDrawer.SetDirection(direction)
	d.Sprites = sprites
	var xMax, yMax float64
	for _, s := range sprites {
		if xMax < s.X() {
			xMax = s.X()
		}
		if yMax < s.Y() {
			yMax = s.Y()
		}
	}
	// d.Direction = direction
	// switch d.Direction {
	// case Upward:
	// 	d.maxPosition = screenSizeY - d.HitPostion
	// 	d.minPosition = -d.HitPostion
	// case Downward:
	// 	d.maxPosition = d.HitPostion
	// 	d.minPosition = -screenSizeY + d.HitPostion
	// case Leftward:
	// 	d.maxPosition = screenSizeX - d.HitPostion
	// 	d.minPosition = -d.HitPostion
	// case Rightward:
	// 	d.maxPosition = d.HitPostion
	// 	d.minPosition = -screenSizeX + d.HitPostion
	// }
	switch d.Direction {
	case Downward, Upward:
		d.margin = yMax / 2
		d.bodyLoss = sprites[1].H()/2 + sprites[2].H()/2 // Todo: should I put enum variables?
	case Leftward, Rightward:
		d.margin = xMax / 2
		d.bodyLoss = sprites[1].W()/2 + sprites[2].W()/2
	}

	// switch d.Direction {
	// case Downward:
	// 	d.boundFarIn = 0 - d.margin
	// 	d.boundNearOut = screenSizeY + d.margin
	// case Upward:
	// 	d.boundFarIn = screenSizeY + d.margin
	// 	d.boundNearOut = 0 - d.margin
	// case Leftward:
	// 	d.boundFarIn = screenSizeX + d.margin
	// 	d.boundNearOut = 0 - d.margin
	// case Rightward:
	// 	d.boundFarIn = 0 - d.margin
	// 	d.boundNearOut = screenSizeX + d.margin
	// }
	return
}

// Update should use existing speed, not the new one.
func (d *NoteLaneDrawer) Update(speed float64) {
	d.Cursor += speed * TimeStep
	// var boundFarIn, boundNearOut float64 // Bounds for farthest, nearest each.
	for d.Farthest.Position-d.Cursor <= d.maxPosition {
		d.Farthest = d.Farthest.Next
	}
	for d.Nearest.Position-d.Cursor <= d.minPosition {
		d.Nearest = d.Nearest.Next
	}
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
	d.Speed = speed
}

// Draw from farthest to nearest.
// So that nearer notes are exposed when overlapped with farther notes.
func (d NoteLaneDrawer) Draw(screen *ebiten.Image) {
	n := d.Farthest
	for ; n != d.Nearest; n = d.Farthest.Prev {
		sprite := d.Sprites[n.Type]
		// op := &ebiten.DrawImageOptions{}
		offset := n.Position - d.Cursor
		switch d.Direction {
		case Downward, Upward:
			sprite.Move(0, offset)
		case Leftward, Rightward:
			sprite.Move(offset, 0)
		}
		// sprite.Draw(screen, op)
		sprite.Draw(screen, nil)
		if n.Type == Head {
			d.DrawLongBody(screen, n)
		}
	}
	if n.Type == Tail {
		d.DrawLongBody(screen, n.Prev)
	}
}

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

// DrawLongBody finds sub-image of Body sprite corresponding to current exposed long body
// and scale the sub-image to (exposed length) / (sub-image length).

// Tail's Position is always larger than Head's.
// In other word, Head is always nearer than Tail.
// Start is Head's, and End is Tail's.
func (d NoteLaneDrawer) DrawLongBody(screen *ebiten.Image, head *Note) {
	tail := head.Next
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
	srcSprite := d.Sprites[Body]
	ratio := length / srcSprite.H()
	srcStart := math.Floor((startPosition - head.Position) / ratio)
	srcEnd := math.Ceil((endPosition - head.Position) / ratio)
	switch d.Direction {
	case Upward, Downward:
		if d.Direction == Downward {
			ratio *= -1
			srcStart *= -1
		}
		op := &ebiten.DrawImageOptions{}
		srcRect := image.Rect(0, int(srcStart), int(srcSprite.W()), int(srcEnd))
		sprite := srcSprite.SubSprite(srcRect)
		op.GeoM.Scale(1, ratio)
		op.GeoM.Translate(0, srcStart)
		sprite.Draw(screen, op)
		// start := top // Top / ground mapping is based on long body's drawn direction
		// end := ground
		// if start < 0 {
		// 	start = 0
		// }
		// if end > screenSizeY {
		// 	end = screenSizeY
		// }
		// exposed := start - end
		// startProportion := (start - ground) / length
		// endProportion := (end - ground) / length
		// startSub := startProportion * d.Sprites[Body].H()
		// endSub := endProportion * d.Sprites[Body].H()
	case Leftward, Rightward:
		if d.Direction == Rightward {
			ratio *= -1
			srcStart *= -1
		}
		op := &ebiten.DrawImageOptions{}
		srcRect := image.Rect(int(srcStart), 0, int(srcEnd), int(srcSprite.H()))
		sprite := srcSprite.SubSprite(srcRect)
		op.GeoM.Scale(ratio, 1)
		op.GeoM.Translate(srcStart, 0)
		sprite.Draw(screen, op)
	}
}
