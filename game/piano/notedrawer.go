package piano

import (
	"image/color"
	"sort"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

// NoteDrawer itself has no effect on the chart and score calculation.
type NoteDrawer struct {
	notes            []Note
	keysAnims        [][4]draws.Animation
	keysLowest       []int // indexes of lowest notes
	cursor           float64
	scaledScreenSize float64 // added in 241006
	keysColor        []color.NRGBA
	keysHolding      []bool
	// h           float64 // used for drawLongNoteBody
}

func NewNoteDrawer(res *Resources, opts *Options, c *Chart) NoteDrawer {
	nd := NoteDrawer{}
	nd.notes = c.notes
	nd.keysAnims = make([][4]draws.Animation, c.keyCount)
	for k := range nd.keysAnims {
		for nk, frames := range res.NotesFramesList {
			a := draws.NewAnimation(frames, 400)
			w := opts.keyWidthsMap[c.keyCount][k]
			h := opts.NoteHeight
			a.SetSize(w, h)

			x := opts.keyPositionXsMap[c.keyCount][k]
			y := opts.KeyPositionY
			if nk == int(Body) {
				a.Locate(x, y, draws.CenterTop)
			} else {
				a.Locate(x, y, draws.CenterBottom)
			}
			nd.keysAnims[k][nk] = a
		}
	}

	// Apply default sample values.
	cd := c.FuncCurrentDynamic()
	// dys.Reset()
	for i, n := range c.notes {
		d := cd(n.Time)
		// d := c.Dynamics.UpdateIndex(n.Time)
		// if n.Sample.Filename == "" {
		// 	c.notes.data[i].Sample.Filename = res.defaultSampleName
		// }
		if n.Sample.Volume == 0 {
			c.notes[i].Sample.Volume = d.Volume
		}
	}

	// Apply TailOffset to Tail's Position.
	cd = c.FuncCurrentDynamic()
	// dys.Reset()
	for i, n := range c.notes {
		if n.Kind != Tail {
			continue
		}
		d := cd(n.Time)
		// d := dys.UpdateIndex(n.Time)
		c.notes[i].position += float64(opts.TailNoteOffset) * d.Speed

		// Tail's Position should be always equal or larger than Head's.
		if n.Prev == -1 {
			continue
		}
		if head := c.notes[n.Prev]; n.position < head.position {
			c.notes[i].position = head.position
		}
	}

	nd.keysLowest = make([]int, c.keyCount)
	copy(nd.keysLowest, c.keysFocusedNote)
	nd.scaledScreenSize = opts.screenSizeY * opts.SpeedScale

	nd.keysColor = make([]color.NRGBA, c.keyCount)
	order := opts.KeyOrders[c.keyCount]
	for k := range nd.keysColor {
		nd.keysColor[k] = opts.NoteColors[order[k]]
	}
	nd.keysHolding = make([]bool, c.keyCount)
	// nd.h = opts.H

	return nd
}

func (nd *NoteDrawer) Update(ka game.KeyboardAction, cursor float64) {
	lowermost := cursor - nd.scaledScreenSize
	for k, lowest := range nd.keysLowest {
		for ni := lowest; ni < len(nd.notes); ni = nd.notes[ni].Next {
			if ni < 0 {
				break
			}
			n := nd.notes[ni]
			if n.position > lowermost {
				break
			}
			// index should be updated outside of if block.
			nd.keysLowest[k] = ni
		}
		// When Head is off the screen but Tail is on,
		// update Tail to Head since drawLongNote uses Head.
		ni := nd.keysLowest[k]
		if ni < 0 {
			continue
		}
		if n := nd.notes[ni]; n.Kind == Tail {
			nd.keysLowest[k] = n.Prev
		}
	}
	nd.cursor = cursor
	nd.keysHolding = ka.KeysHolding()
}

// Notes are fixed. Lane itself moves, all notes move as same amount.
func (nd NoteDrawer) Draw(dst draws.Image) {
	uppermost := nd.cursor + nd.scaledScreenSize
	for k, lowest := range nd.keysLowest {
		var nis []int
		for ni := lowest; ni < len(nd.notes); ni = nd.notes[ni].Next {
			if ni < 0 {
				break
			}
			n := nd.notes[ni]
			if n.position > uppermost {
				break
			}
			nis = append(nis, ni)
		}

		// Make farther notes overlapped by nearer notes.
		sort.Sort(sort.Reverse(sort.IntSlice(nis)))

		for _, ni := range nis {
			if ni < 0 {
				break
			}
			n := nd.notes[ni]
			// Make long note's body overlapped by its Head and Tail.
			if n.Kind == Head {
				nd.drawLongNoteBody(dst, n)
			}

			a := nd.keysAnims[k][n.Kind]
			pos := n.position - nd.cursor
			a.Move(0, -pos)
			a.ColorScale.ScaleWithColor(nd.keysColor[k])
			if n.scored {
				a.ColorScale.ScaleWithColor(color.Gray{128})
				// op.ColorM.ChangeHSV(0, 0.3, 0.3)
			}
			a.Draw(dst)
		}
	}
}

// drawLongNoteBody draws stretched long note body sprite.
func (nd NoteDrawer) drawLongNoteBody(dst draws.Image, head Note) {
	tail := nd.notes[head.Next]
	if head.Kind != Head || tail.Kind != Tail {
		return
	}

	a := nd.keysAnims[head.Key][Body]
	if !nd.keysHolding[head.Key] {
		a.Reset()
	}

	length := tail.position - head.position
	// length += nd.h
	if length < 0 {
		length = 0
	}
	a.SetSize(a.W(), length)

	pos := tail.position - nd.cursor
	a.Move(0, -pos)
	a.ColorScale.ScaleWithColor(nd.keysColor[tail.Key])
	if tail.scored {
		a.ColorScale.ScaleWithColor(color.Gray{128})
	}
	a.Draw(dst)
}
