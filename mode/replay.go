package mode

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

type ReplayPlayer struct {
	states []input.KeyboardState
	index  int // states[index] is last latest state
}

// The type of return value should be pointer because to implement
// interface, all methods should be either value or pointer receiver.
func NewReplayPlayer(f *osr.Format, keyCount int) *ReplayPlayer {
	return &ReplayPlayer{states: f.KeyboardStates(keyCount)}
}
func NewReplayPlayerFromFile(fsys fs.FS, name string, keyCount int) (*ReplayPlayer, error) {
	file, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}

	f, err := osr.NewFormat(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse replay file: %s", err)
	}

	return NewReplayPlayer(f, keyCount), nil
}

func (rp *ReplayPlayer) Fetch(now int32) (kas []input.KeyboardAction) {
	count := 0
	for _, s := range rp.states[rp.index:] {
		if s.Time > now {
			break
		}
		count++
	}

	states := rp.states[rp.index : rp.index+count]
	if len(states) == 0 {
		ps := make([]bool, len(rp.states[0].Pressed))
		dummy := input.KeyboardState{Time: now, Pressed: ps}
		states = append(states, dummy, dummy)
	} else if len(states) == 1 {
		dummy := input.KeyboardState{Time: now, Pressed: states[0].Pressed}
		states = append(states, dummy)
	}

	lps := states[0].Pressed
	for _, s := range states[1:] {
		ps := s.Pressed
		as := input.KeyActions(lps, ps)
		ka := input.KeyboardAction{Time: s.Time, Action: as}
		kas = append(kas, ka)
		lps = s.Pressed
	}

	rp.index += count
	return
}
func (rp *ReplayPlayer) Output() []input.KeyboardState { return rp.states }

// ReplayPlayer does nothing at these methods.
func (rp *ReplayPlayer) Pause()  {}
func (rp *ReplayPlayer) Resume() {}
func (rp *ReplayPlayer) Close()  {}
