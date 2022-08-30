package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
)

// Strategy of Piano mode
// Calculate position of each note in advance
// Parameter: SpeedBase, BPM Ratio, BeatScale
// Calculate current HitPosition only.
// For other notes, just calculate the difference between HitPosition.
type NoteType int
type BaseNote struct {
	Type         NoteType
	Time         int64
	Time2        int64
	SampleName   string // SampleFilename
	SampleVolume float64

	Next *BaseNote
	Prev *BaseNote // For accessing to Head from Tail.
	// NextTail *BaseNote // For drawing long body faster.
	Position float64 // Scaled x or y value.
}

// TimeStep is expected to be integer.
// TPS should be either multiple or divisor of 1000.
var TimeStep float64 = 1000 / float64(ebiten.MaxTPS())

// NoteLaneDrawer's tick should be consistent with ScenePlay.
type NoteLaneDrawer struct {
	Tick    int
	Sprites [4]draws.Sprite //  map[NoteType]draws.Sprite // []draws.Sprite
	// Note       *BaseNote
	Farthest   *BaseNote
	Nearest    *BaseNote
	Cursor     float64
	HitPostion float64
	Speed      float64 // BPM (or BPM ratio) * BeatScale
	// Sizes      map[NoteType]float64 // Cache for Sprites' sizes. // Todo: Sizes -> halfSizes
	// MaxSize    float64              // Either max width / height. // Todo: remove
	margin   float64 // Half of max sizes of sprites.
	bodyLoss float64 // Head/2 + Tail/2
	// boundFarIn   float64 // Bound for Farthest note being fetched.
	// boundNearOut float64 // Bound for Nearest note being flushed.
	Direction
}
type Direction int

const (
	Downward Direction = iota // e.g., Piano mode.
	Upward                    // e.g., Rhythm games using feet.
	Leftward                  // e.g., Drum mode.
	Rightward
)

// [4]draws.Sprite{Note, Head, Tail, Body}
func NewNoteLaneDrawer(direction Direction) (d NoteLaneDrawer) {
	var xMax, yMax float64
	// margin:=maxSize/2
	// var boundTop, boundBottom float64
	// var boundLeft, boundRight float64
	switch direction {
	case Downward:
		d.boundIn = 0 - yMax/2
		d.boundOut = screenSizeY + yMax/2
	case Leftward:
		d.boundFarIn = screenSizeX + xMax2
		d.boundNearOut = 0 - xMax/2
	}
}
func (d *NoteLaneDrawer) Update(speed float64) {
	// Update should use existing speed, not the new one.
	d.Cursor += speed * TimeStep
	// p1:=d.Farthest.Position-d.Cursor
	for d.ScreenPosition(d.Farthest) >= d.boundIn {
		d.Farthest = d.Farthest.Next
	}
	for d.ScreenPosition(d.Nearest) >= d.boundOut {
		d.Nearest = d.Nearest.Next
	}
	switch d.Direction {
	case Downward:
		for d.ScreenPosition(d.Farthest)+d.margin >= 0 {
			d.Farthest = d.Farthest.Next
		}
		for d.ScreenPosition(d.Nearest)-d.margin >= screenSizeY {
			d.Nearest = d.Nearest.Next
		}
	case Upward:
		for d.ScreenPosition(d.Farthest)-d.margin >= screenSizeY {
			d.Farthest = d.Farthest.Next
		}
		for d.ScreenPosition(d.Nearest)+d.margin >= 0 {
			d.Nearest = d.Nearest.Next
		}
	case Leftward:
		for d.ScreenPosition(d.Farthest)-d.margin >= screenSizeX {
			d.Farthest = d.Farthest.Next
		}
		for d.ScreenPosition(d.Nearest)+d.margin >= 0 {
			d.Nearest = d.Nearest.Next
		}
	case Rightward:
		for d.ScreenPosition(d.Farthest)+d.margin >= 0 {
			d.Farthest = d.Farthest.Next
		}
		for d.ScreenPosition(d.Nearest)-d.margin >= screenSizeX {
			d.Nearest = d.Nearest.Next
		}
	}
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

// Draw from farthest to nearest
func (d NoteLaneDrawer) Draw(screen *ebiten.Image) {
	// var offset float64
	// switch d.Direction {
	// case Downward, Rightward:
	// 	offset = -d.MaxSize / 2
	// case Upward, Leftward:
	// 	offset = d.MaxSize / 2
	// }
	n := d.Note
	// var prev *BaseNote
loop:
	for ; ; n = n.Next {
		switch d.Direction {
		case Downward:
			if d.ScreenPosition(n, -d.MaxSize) < 0 {
				break loop
			}
		case Upward:
			if d.ScreenPosition(n, +d.MaxSize) > screenSizeY {
				break loop
			}
		case Leftward:
			if d.ScreenPosition(n, +d.MaxSize) > screenSizeX {
				break loop
			}
		case Rightward:
			if d.ScreenPosition(n, -d.MaxSize) < 0 {
				break loop
			}
		}
		op := &ebiten.DrawImageOptions{}
		offset := n.Position - d.Cursor
		switch d.Direction {
		case Downward, Upward:
			op.GeoM.Translate(0, offset)
		case Leftward, Rightward:
			op.GeoM.Translate(offset, 0)
		}
		d.Sprites[n.Type].Draw(screen, op)
		if n.Type == Tail {
			d.DrawLongBody(screen, n)
		}
		// prev = n
	}
	if n.Type == Head {
		d.DrawLongBody(screen, n.Next)
	}
}

//	func (d NoteLaneDrawer) ScreenPosition0(n *BaseNote, size float64) float64 {
//		return d.HitPostion - (n.Position - d.Cursor) + size/2
//	}
func (d NoteLaneDrawer) ScreenPosition(n *BaseNote) float64 {
	pos := n.Position - d.Cursor // Relative position of note.
	// switch d.Direction {
	// case Downward:
	// 	return -pos + d.HitPostion
	// case Upward:
	// 	return pos + d.HitPostion
	// case Leftward:
	// 	return pos + d.HitPostion
	// case Rightward:
	// 	return -pos + d.HitPostion
	// }
	switch d.Direction {
	case Downward, Rightward:
		pos = -pos
	}
	return d.HitPostion + pos
}

// func ()
// DrawLongBody finds sub-image of Body sprite corresponding to current exposed long body
// and scale the sub-image to (exposed length) / (sub-image length).
func (d NoteLaneDrawer) DrawLongBody(screen *ebiten.Image, head *BaseNote) {
	tail := head.Next
	length := tail.Position - head.Position
	length -= -d.bodyLoss
	// length -= d.Sizes[Tail]/2 + d.Sizes[Head]/2
	// var start, end float64
	// ground := d.ScreenPosition(head) // -d.Sizes[Head])
	// top := d.ScreenPosition(tail)    // +d.Sizes[Tail])
	ground := head.Position
	top := tail.Position
	switch d.Direction {
	case Downward:
		start := top // Top / ground mapping is based on long body's drawn direction
		end := ground
		if start < 0 {
			start = 0
		}
		if end > screenSizeY {
			end = screenSizeY
		}
		exposed := start - end
		startProportion := (start - ground) / length
		endProportion := (end - ground) / length
		startSub := startProportion * d.Sprites[Body].H()
		endSub := endProportion * d.Sprites[Body].H()
	case Upward:
	case Leftward:

	case Rightward:
	}
}
