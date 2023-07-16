package input

type Keyboard struct {
	states []KeyboardState
	index  int // states[index] is last latest state
	*KeyboardListener
}

type KeyboardState struct {
	Time    int32
	Pressed []bool
}

type KeyboardAction struct {
	Time   int32
	Action []KeyActionType
}

func NewKeyboard(keys []Key) Keyboard {
	return Keyboard{KeyboardListener: NewKeyboardListener(keys)}
}

// NewKeyboardFromStates is for replay.
func NewKeyboardFromStates(states []KeyboardState) Keyboard {
	return Keyboard{states: states}
}

func (kb *Keyboard) Fetch(now int32) []KeyboardAction {
	kas := []KeyboardAction{}

	add := 0
	for _, next := range kb.states[kb.index+1:] {
		if next.Time > now {
			break
		}
		add++
	}

	// Beware: states can manipulate kb.states.
	states := make([]KeyboardState, add+1)
	copy(states, kb.states[kb.index:kb.index+add+1])

	// states should be at least two elements to get KeyboardAction.
	// If states is empty, add dummy state.
	if len(states) == 0 {
		blank := make([]bool, len(kb.states[0].Pressed))
		dummy := KeyboardState{Time: now, Pressed: blank}
		states = append(states, dummy)
	}

	// Time of the last state is always 'now'.
	currentState := KeyboardState{
		Time:    now,
		Pressed: states[len(states)-1].Pressed,
	}
	if len(states) <= 1 || states[len(states)-1].Time != now {
		states = append(states, currentState)
	}

	lps := states[0].Pressed
	for _, s := range states[1:] {
		ps := s.Pressed
		as := KeyActions(lps, ps)
		ka := KeyboardAction{Time: s.Time, Action: as}
		kas = append(kas, ka)
		lps = s.Pressed
	}
	kb.index += add

	return kas
}

func (kb Keyboard) Output() []KeyboardState { return kb.states }
