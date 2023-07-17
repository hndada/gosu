package input

type KeyboardReader struct {
	index  int // states[index] is last element of fetched state.
	states []KeyboardState
}

func NewKeyboardReader(states []KeyboardState) KeyboardReader {
	return KeyboardReader{states: states}
}

// No worry of accessing with nil pointer.
// https://go.dev/play/p/B4Z1LwQC_jP
func (kb KeyboardReader) IsEmpty() bool { return len(kb.states) == 0 }

// It is great to wrap a slice with a struct and name it in the singular form,
// as it ensures the slice is always handled in a packed form.
type KeyboardState struct {
	Time    int32
	Presses []bool
}

// KeyboardAction is for handling keyboard states conveniently.
type KeyboardAction struct {
	Time       int32
	KeyActions []KeyActionType
}

func (kb *KeyboardReader) Read(now int32) []KeyboardAction {
	kas := []KeyboardAction{}

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

// Tidy removes redundant states.
func (kb *KeyboardReader) Tidy() {
	if len(kb.states) == 0 {
		return
	}

	news := []KeyboardState{}
	last := kb.states[0]
	for _, s := range kb.states[1:] {
		if areStatesEqual(last, s) {
			continue
		}
		news = append(news, s)
		last = s
	}
	kb.states = news
}

func areStatesEqual(old, new KeyboardState) bool {
	for k, p := range new.Presses {
		lp := old.Presses[k]
		if lp != p {
			return false
		}
	}
	return true
}

func (kb KeyboardReader) Output() []KeyboardState { return kb.states }
