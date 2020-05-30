package mania

const defaultHand = right

func (beatmap *ManiaBeatmap) determineAlters() {
	// affect idx has already been calculated
	var leftCount, rightCount int
	var det int
	for i, note := range beatmap.Notes {
		if note.hand != alter {
			continue
		}
		// rule 1: use default hand if there is a note very next to alterable note
		if note.chord[note.Key+defaultHand] != noFound {
			det = defaultHand
		} else {
			leftCount, rightCount = 0, 0
			for key := note.Key - 1; key >= 0; key-- {
				if note.chord[key] <= noFound {
					break
				}
				leftCount++
			}
			for key := note.Key + 1; key < len(note.chord); key++ {
				if note.chord[key] <= noFound {
					break
				}
				rightCount++
			}
			// rule 2: alter follows the hand which has more notes in the chord
			// rule 3: alter follows default hand if each hand has same number of notes
			switch {
			case leftCount > rightCount:
				det = left
			case leftCount < rightCount:
				det = right
			default: // if two counts are same
				det = defaultHand
			}
		}
		beatmap.Notes[i].hand = det
	}
}
