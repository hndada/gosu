package lv

import (
	"github.com/hndada/gosu/internal/tools"
	"github.com/hndada/gosu/mode/mania"
)

const noFound = tools.NoFound

func markAffect(ns []mania.Note) {
	// scan the notes whether its affectable, chordable or not
	// further lanes which cause *miss* when be hit at the same time goes 'chuck cutter'
	if len(ns)==0 {return}
	pnIdxs := tools.GetIntSlice(len(ns[0].chord), noFound)
	for i, n := range ns {
		markPrevAffect(ns, i, pnIdxs)
		markNextAffect(ns, i)
		pnIdxs[n.Key] = i
	}
}

func markPrevAffect(ns []mania.Note, i int, pnIdxs []int) {
	n := ns[i]
	var pn mania.Note // previous note
	var elapsedTime int
	for pnKey, pnIdx := range pnIdxs {
		if pnIdx == noFound {
			continue
		}
		pn = ns[pnIdx]
		elapsedTime = n.Time - pn.Time
		switch pn.Key == n.Key {
		case true: // jack
			if elapsedTime <= maxDeltaJack {
				ns[i].trillJack[pnKey] = pnIdx
			}
		default:
			if elapsedTime <= maxDeltaTrill {
				if elapsedTime <= maxDeltaChord { // chord
					ns[i].chord[pnKey] = pnIdx
					ns[i].trillJack[pnKey] = noFound
				} else { // trill
					ns[i].trillJack[pnKey] = pnIdx
				}
			}
		}
	}
	ns[i].chord[n.Key] = i // putting note itself to chord
}

func markNextAffect(ns []mania.Note, i int) {
	n := ns[i]
	var nn mania.Note
	nnIdx := i + 1
	var elapsedTime int
	for nnIdx < len(ns) {
		nn = ns[nnIdx]
		elapsedTime = nn.Time - n.Time
		if elapsedTime > maxDeltaTrill {
			break
		}

		if nn.NoteType != mania.NtHoldTail &&
			nn.Key != n.Key && // jack is not relevant
			ns[i].chord[nn.Key] == noFound { // prev notes is prior to next notes
			switch {
			case elapsedTime <= maxDeltaChord:
				ns[i].chord[nn.Key] = nnIdx
			// default:
			// 	ns[i].chord[nn.Key] = cut
			}
		}
		nnIdx++
	}
}
