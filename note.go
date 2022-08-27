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
	Tick       int
	Sprites    map[NoteType]draws.Sprite // []draws.Sprite
	Note       *BaseNote
	Cursor     float64
	HitPostion float64
	Speed      float64              // BPM (or BPM ratio) * BeatScale
	Sizes      map[NoteType]float64 // Cache for Sprites' sizes.
	MaxSize    float64              // Either max width / height.
	Direction
}
type Direction int

const (
	Downward Direction = iota // e.g., Piano mode.
	Upward                    // e.g., Rhythm games using feet.
	Leftward                  // e.g., Drum mode.
	Rightward
)

func (d *NoteLaneDrawer) Update(speed float64) {
	// Update should use existing speed, not the new one.
	d.Cursor += speed * TimeStep
	switch d.Direction {
	case Downward:
		for d.Note.Position-d.MaxSize/2 >= screenSizeY {
			d.Note = d.Note.Next
		}
	case Leftward:
		for d.Note.Position+d.MaxSize/2 < 0 {
			d.Note = d.Note.Next
		}
	}
	d.Speed = speed
}
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

func (d NoteLaneDrawer) ScreenPosition(n *BaseNote, size float64) float64 {
	return d.HitPostion - (n.Position - d.Cursor) + size/2
}

// DrawLongBody finds sub-image of Body sprite corresponding to current exposed long body
// and scale the sub-image to (exposed length) / (sub-image length).
func (d NoteLaneDrawer) DrawLongBody(screen *ebiten.Image, tail *BaseNote) {
	head := tail.Prev
	var start, end float64
	switch d.Direction {
	case Downward:
		start := d.ScreenPosition(tail, +d.Sizes[Tail])
		end := d.ScreenPosition(head, -d.Sizes[Head])
		length := end - start
		if start < 0 {
			start = 0
		}
		if end > screenSizeY {
			end = screenSizeY
		}
	case Upward:
	case Leftward:
	case Rightward:
	}
}
