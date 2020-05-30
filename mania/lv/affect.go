package mania

import (
	"math"

	"github.com/hndada/gosu/game/tools"
)

const noFound = tools.NoFound
const cut = noFound - 1

func (beatmap *ManiaBeatmap) markAffect() {
	// scan the notes whether its affectable, chordable or not
	// further lanes which cause *miss* when be hit at the same time goes 'chuck cutter'
	prevNotesIdxs := tools.GetIntSlice(beatmap.Keymode, noFound)
	xValues := tools.SolveX(beatmap.Curves["TrillChord"], 0)
	if len(xValues) != 2 {
		panic(&tools.ValError{"Length of xValues of segments set TrillChord",
			tools.Itoa(len(xValues)), tools.ErrFlow})
	}
	beatmap.maxChunkDelta = int(math.Round(xValues[0]))
	beatmap.maxTrillDelta = int(math.Round(xValues[1]))
	beatmap.maxJackDelta = int(math.Round(tools.SolveX(beatmap.Curves["Jack"], 0)[0]))
	beatmap.markChordCut()
	for i, note := range beatmap.Notes {
		beatmap.markPrevAffect(i, prevNotesIdxs)
		beatmap.markNextAffect(i, prevNotesIdxs)
		beatmap.cutChord(i, prevNotesIdxs)
		prevNotesIdxs[note.Key] = i
	}
}

func (beatmap *ManiaBeatmap) markChordCut() {
	// a hold note acts as a chord cutter
	// except to notes which are chordable with the hold note head
	var holdStartTime, holdEndTime int
	for _, holdNoteHead := range beatmap.Notes {
		if holdNoteHead.NoteType != NtHoldHead {
			continue
		}
		holdStartTime = holdNoteHead.Time
		holdEndTime = holdNoteHead.OpponentTime
		for i, note := range beatmap.Notes {
			if note.Key == holdNoteHead.Key ||
				tools.AbsInt(holdStartTime-note.Time) <= beatmap.maxChunkDelta ||
				note.Time-holdEndTime > beatmap.maxTrillDelta {
				continue
			}
			beatmap.Notes[i].chord[holdNoteHead.Key] = cut
		}
	}
}

func (beatmap *ManiaBeatmap) markPrevAffect(i int, prevNotesIdxs []int) {
	note := beatmap.Notes[i]
	var prevNote ManiaNote
	var elapsedTime int
	for prevNoteKey, prevNoteIdx := range prevNotesIdxs {
		if prevNoteIdx == noFound {
			continue
		}
		prevNote = beatmap.Notes[prevNoteIdx]
		elapsedTime = note.Time - prevNote.Time
		switch prevNote.Key == note.Key {
		case true: // jack
			if elapsedTime <= beatmap.maxJackDelta {
				beatmap.Notes[i].trillJack[prevNoteKey] = prevNoteIdx
			}
		default: // trill or chord
			if elapsedTime <= beatmap.maxTrillDelta {
				if elapsedTime <= beatmap.maxChunkDelta {
					beatmap.Notes[i].chord[prevNoteKey] = prevNoteIdx
					beatmap.Notes[i].trillJack[prevNoteKey] = noFound
				} else {
					beatmap.Notes[i].trillJack[prevNoteKey] = prevNoteIdx
				}
			}
		}
	}
	beatmap.Notes[i].chord[note.Key] = i // putting note itself to chord
}
func (beatmap *ManiaBeatmap) markNextAffect(i int, prevNotesIdxs []int) {
	note := beatmap.Notes[i]
	var nextNote ManiaNote
	nextNoteIdx := i + 1
	var elapsedTime int
	for nextNoteIdx < len(beatmap.Notes) {
		nextNote = beatmap.Notes[nextNoteIdx]
		elapsedTime = nextNote.Time - note.Time
		if elapsedTime > beatmap.maxTrillDelta {
			break
		}

		if nextNote.NoteType != NtHoldTail &&
			nextNote.Key != note.Key && // jack is not relevant
			beatmap.Notes[i].chord[nextNote.Key] == noFound { // prev notes is prior to next notes
			switch {
			case elapsedTime <= beatmap.maxChunkDelta:
				beatmap.Notes[i].chord[nextNote.Key] = nextNoteIdx
			default:
				beatmap.Notes[i].chord[nextNote.Key] = cut
			}
		}
		nextNoteIdx++
	}
}

func (beatmap *ManiaBeatmap) cutChord(i int, prevNotesIdxs []int) {
	// notes over the chord cutter will be dropped from chord.
	note := beatmap.Notes[i]
	cutLeft, cutRight := false, false
	for key := note.Key - 1; key >= 0; key-- {
		if note.chord[key] == cut {
			cutLeft = true
		}
		if cutLeft {
			beatmap.Notes[i].chord[key] = noFound
		}
	}
	for key := note.Key + 1; key < beatmap.Keymode; key++ {
		if note.chord[key] == cut {
			cutRight = true
		}
		if cutRight {
			beatmap.Notes[i].chord[key] = noFound
		}
	}
}
