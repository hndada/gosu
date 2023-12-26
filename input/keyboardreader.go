package input

import "time"

func (kb *Keyboard) Read(now time.Time) (ks []KeyboardState) {
	add := 0
	for _, next := range kb.states[kb.index+1:] {
		if next.Time > now {
			break
		}
		add++
	}

	count := add + 1
	// Beware: states can manipulate kb.states.
	states := make([]KeyboardState, count)
	copy(states, kb.states[kb.index:kb.index+count])

	// states should be at least two elements to get KeyboardAction.
	// If states is empty, add dummy state.
	// if len(states) == 0 {
	// 	blank := make([]bool, len(kb.states[0].Presses))
	// 	dummy := KeyboardState{Time: now, Presses: blank}
	// 	states = append(states, dummy)
	// }

	// Time of the last state is always 'now'.
	if states[len(states)-1].Time != now { // len(states) <= 1 ||
		currentState := KeyboardState{
			Time:    now,
			Presses: states[len(states)-1].Presses,
		}
		states = append(states, currentState)
	}

	lps := states[0].Presses
	for _, s := range states[1:] {
		ps := s.Presses
		as := KeyActions(lps, ps)
		ka := KeyboardAction{Time: s.Time, KeyActions: as}
		kas = append(kas, ka)
		lps = s.Presses
	}
	kb.index += add

	return kas
}
