package main

// PlayNote is for in-game. Handled by pointers to modify its fields easily.
type PlayNote struct {
	Note
	Prev   *PlayNote
	Next   *PlayNote
	Scored bool

	// Sprite
	// LongSprite
}

func NewPlayNotes(c *Chart) []*PlayNote {
	pns := make([]*PlayNote, 0, len(c.Notes))
	prevs := make([]*PlayNote, c.Parameter.KeyCount)
	for _, n := range c.Notes {
		prev := prevs[n.Key]
		next := &PlayNote{
			Note: n,
			Prev: prev,
		}
		if prev != nil { // Set Next value later
			prev.Next = next
		}
		prevs[n.Key] = next
	}
	return pns
}
func (n PlayNote) PlaySE() {}

func (n *PlayNote) UpdateSprite() {}

// type TimeStamp struct {
// 	Time     int64
// 	NextTime int64
// 	Position float64
// 	Factor   float64
// }
