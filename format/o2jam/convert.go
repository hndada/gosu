// The following code is generated from ChatGPT,
// with minor manual adjustments for integration into gosu.
package o2jam

import (
	"fmt"
	"math"
	"sort"

	"github.com/hndada/gosu/format/osu"
)

// convertOJNToOsuFormat converts parsed OJN into an *osu.Format.
// This function tries to produce the fields used by gosu's game.NewDynamics
// and piano.newChartNotes: TimingPoints, HitObjects, CircleSize, Duration, Title, Artist, Creator, AudioFilename.
func ConvertOJNToOsuFormat(ojn *OJN, diff int) (*osu.Format, error) {
	if ojn == nil {
		return nil, fmt.Errorf("nil ojn")
	}
	if diff < 0 || diff > 2 {
		return nil, fmt.Errorf("invalid diff %d", diff)
	}

	// pick section
	var sec Section
	switch diff {
	case 0:
		sec = ojn.Easy
	case 1:
		sec = ojn.Normal
	case 2:
		sec = ojn.Hard
	}

	// Build a minimal osu.Format. The exact field names of your osu.Format
	// type must match; adapt if necessary. The following fields are commonly used.
	f := &osu.Format{}
	// Set basic metadata (adapt field names if needed).
	// Many osu.Format implementations call them Title/Artist/Creator/AudioFilename
	// If your type uses different names, change them.
	f.Title = ojn.Header.Title()
	f.Artist = ojn.Header.Artist()
	f.Creator = ojn.Header.Noter()
	f.AudioFilename = ojn.Header.OJMFile()

	// === Timing conversion ===
	// Strategy:
	//  - Use header BPM as default mainBPM.
	//  - Scan packages for BPM channel events and record a per-measure BPM override.
	//  - Compute measure start times by iterating measures from 0..maxMeasure.
	//  - For each package event, compute its absolute ms time and produce hitobjects.
	mainBPM := float64(ojn.Header.BPM)
	if mainBPM <= 0 {
		mainBPM = 120.0
	}
	// default meter
	const meter = 4

	// collect per-measure BPM overrides from CHAN_BPM packages
	measureBPM := map[uint32]float64{}
	for _, pkg := range sec.Packages {
		if pkg.Header.Channel == ChanBPM {
			for _, ev := range pkg.Events {
				if ev.BPM != nil {
					measureBPM[pkg.Header.Measure] = float64(*ev.BPM)
				}
			}
		}
	}
	// find max measure index
	var maxMeasure uint32
	for _, pkg := range sec.Packages {
		if pkg.Header.Measure > maxMeasure {
			maxMeasure = pkg.Header.Measure
		}
	}

	// build measureStart map (ms) by accumulating measure durations
	measureStart := make(map[uint32]int32, maxMeasure+1)
	cur := 0.0
	for m := uint32(0); m <= maxMeasure; m++ {
		measureStart[m] = int32(math.Round(cur))
		bpm := mainBPM
		if mb, ok := measureBPM[m]; ok && mb > 0 {
			bpm = mb
		}
		measureDur := float64(meter) * (60000.0 / bpm)
		cur += measureDur
	}

	// Build basic timing points: initial main BPM + explicit measure BPMs
	tps := make([]osu.TimingPoint, 0, len(measureBPM)+1)
	// initial timing point at t=0
	tps = append(tps, osu.TimingPoint{
		Time:        0,
		Uninherited: true,
		// adapt the field name used by your osu.TimingPoint (BPMValue/BPM/BeatLength etc.)
		// many osu structs use BPM or BeatLength; pick the correct one in your repo.
		BeatLength: 60000.0 / mainBPM,
		Meter:      meter,
		Volume:     100,
	})
	// per-measure timing points
	for m, bpm := range measureBPM {
		start := measureStart[m]
		if start <= 0 {
			// if it's time 0 (already represented) skip
			continue
		}
		tps = append(tps, osu.TimingPoint{
			Time:        int(start),
			Uninherited: true,
			BeatLength:  60000.0 / bpm,
			Meter:       meter,
			Volume:      100,
		})
	}
	// sort tps just in case
	sort.Slice(tps, func(i, j int) bool { return tps[i].Time < tps[j].Time })
	f.TimingPoints = tps

	// === Hit object conversion ===
	// playable lanes are ChanLane1 .. ChanLane7 (2..8)
	const laneMin = ChanLane1
	const laneMax = ChanLane7

	// temporary bucket per lane for pairing long note heads/tails
	type tmpNote struct {
		TimeMs int32
		Type   uint8 // NoteNormal/NoteLongHead/NoteLongTail
		Val    uint16
		Vol    uint8
		Pan    uint8
	}
	perKey := make([][]tmpNote, int(laneMax-laneMin+1))

	// helper to compute absolute time for event in package
	calcAbsTime := func(measure uint32, evPos int, eventsInPkg int) int32 {
		start := float64(measureStart[measure])
		var frac float64 = 0
		if eventsInPkg > 0 {
			frac = float64(evPos) / float64(eventsInPkg)
		}
		// BPM for the measure
		bpm := mainBPM
		if mb, ok := measureBPM[measure]; ok && mb > 0 {
			bpm = mb
		}
		measureDur := float64(meter) * (60000.0 / bpm)
		return int32(math.Round(start + frac*measureDur))
	}

	// gather tmp notes
	for _, pkg := range sec.Packages {
		ch := pkg.Header.Channel
		if ch < laneMin || ch > laneMax {
			continue
		}
		key := int(ch - laneMin) // 0..6
		for _, ev := range pkg.Events {
			if ev.Note == nil || ev.Note.Value == 0 {
				continue
			}
			abs := calcAbsTime(pkg.Header.Measure, ev.Position, int(pkg.Header.Events))
			perKey[key] = append(perKey[key], tmpNote{
				TimeMs: abs,
				Type:   ev.Note.Type,
				Val:    ev.Note.Value,
				Vol:    ev.Note.Volume,
				Pan:    ev.Note.Pan,
			})
		}
	}

	xs := []int{36, 109, 182, 256, 329, 402, 475} // 7-key x positions
	// Now convert perKey tmp notes into osu.HitObjects with head+end pairing
	hitObjects := make([]osu.HitObject, 0)

	for key, notes := range perKey {
		// ensure notes are sorted by time (they usually are)
		sort.Slice(notes, func(i, j int) bool { return notes[i].TimeMs < notes[j].TimeMs })

		var pendingHeadIndex int = -1
		for _, n := range notes {
			switch n.Type {
			case NoteNormal:
				ho := osu.HitObject{
					Time: int(n.TimeMs),
					X:    xs[key],
					Y:    192, // fixed y for now
					// sample fields left blank for now; we'll handle sample mapping separately
				}
				hitObjects = append(hitObjects, ho)
			case NoteLongHead:
				// create head and keep its index as pending
				ho := osu.HitObject{
					Time: int(n.TimeMs),
					X:    xs[key],
					Y:    192, // fixed y for now
					// mark as hold by EndTime==0 for now; we'll set EndTime later
				}
				hitObjects = append(hitObjects, ho)
				pendingHeadIndex = len(hitObjects) - 1
			case NoteLongTail:
				// close the pending head if present
				if pendingHeadIndex != -1 {
					// set end on pending head
					hitObjects[pendingHeadIndex].EndTime = int(n.TimeMs)
					// reset pending
					pendingHeadIndex = -1
				} else {
					// tail without head -> treat as normal note
					ho := osu.HitObject{
						Time: int(n.TimeMs),
						X:    xs[key],
						Y:    192, // fixed y for now
					}
					hitObjects = append(hitObjects, ho)
				}
			default:
				// unknown type -> normal
				ho := osu.HitObject{
					Time: int(n.TimeMs),
					X:    xs[key],
					Y:    192, // fixed y for now
				}
				hitObjects = append(hitObjects, ho)
			}
		}
		// leftover pending head without tail -> leave as normal (no EndTime)
	}

	// sort hitObjects by time (and column)
	sort.Slice(hitObjects, func(i, j int) bool {
		if hitObjects[i].Time == hitObjects[j].Time {
			return hitObjects[i].X < hitObjects[j].X
		}
		return hitObjects[i].Time < hitObjects[j].Time
	})

	// assign to format
	f.HitObjects = hitObjects

	// duration: last object's EndTime or Time
	// Excluded: No duration field in osu.Format
	// var dur int
	// for _, ho := range hitObjects {
	// 	if ho.EndTime > dur {
	// 		dur = ho.EndTime
	// 	}
	// 	if ho.Time > dur {
	// 		dur = ho.Time
	// 	}
	// }
	// f.DurationMs = dur
	f.Mode = osu.ModeMania
	f.CircleSize = 7 // O2Jam is 7 keys by default

	return f, nil
}
