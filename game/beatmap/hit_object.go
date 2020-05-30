package beatmap

import (
	"errors"
	"strings"

	"github.com/hndada/gosu/game/tools"
)

const (
	NtNote = 1 << iota
	NtSlider
	NewCombo
	NtSpinner
	ColorSkip1
	ColorSkip2
	ColorSkip3
	NtHoldNote
	LastNoteType
)
const ComboMask = ^(NewCombo + ColorSkip1 + ColorSkip2 + ColorSkip3)

var NtArray = [...]int{NtNote, NtSlider, NtSpinner, NtHoldNote}

const (
	Normal = iota
	Whistle
	Finish
	Clap
)

type HitObject struct {
	X         int
	Y         int
	StartTime int
	NoteType  int
	HitSound  int
	EndTime   int
	*SliderParams
}

type SliderParams struct {
	CurveType   string   // format: one letter
	CurvePoints [][2]int // format: slice of paired integers
	Slides      int
	Length      float64
}

// hit circle, slider, spinner, hold in a order
// x,y,time,type,hitSound,hitSample
// x,y,time,type,hitSound,curveType|curvePoints,slides,length,edgeSounds,edgeSets,hitSample
// x,y,time,type,hitSound,endTime,hitSample
// x,y,time,type,hitSound,endTime:hitSample
func (hitObject *HitObject) parseNote(line string) error {
	vs := strings.Split(line, `,`)
	x, err := tools.Atoi(vs[0])
	if err != nil {
		return err
	}
	hitObject.X = x
	y, err := tools.Atoi(vs[1])
	if err != nil {
		return err
	}
	hitObject.Y = y
	startTime, err := tools.Atoi(vs[2])
	if err != nil {
		return err
	}
	hitObject.StartTime = startTime
	noteType, err := tools.Atoi(vs[3])
	if err != nil {
		return err
	}
	noteType, err = extractNoteType(noteType)
	if err != nil {
		return err
	}
	hitObject.NoteType = noteType
	hitSound, err := tools.Atoi(vs[4])
	if err != nil {
		return err
	}
	hitObject.HitSound = hitSound

	switch hitObject.NoteType {
	case NtSlider:
		// curveType|curvePoints, slides, length
		var params SliderParams
		curveValues := strings.Split(vs[5], `|`)
		params.CurveType = curveValues[0]
		params.CurvePoints = make([][2]int, len(curveValues)-1)
		for i, v := range curveValues[1:] {
			params.CurvePoints[i] = parsePoints(v)
		}
		slides, err := tools.Atoi(vs[6])
		if err != nil {
			return err
		}
		params.Slides = slides

		params.Length = tools.Atof(vs[7])
		hitObject.SliderParams = &params
	case NtSpinner:
		hitObject.EndTime = tools.Atoi(vs[5])
	case NtHoldNote:
		hitObject.EndTime = tools.Atoi(strings.Split(vs[5], `:`)[0])
	}
	return HitObjects
}

func extractNoteType(v int) (int, error) {
	nt := v & ComboMask
	for _, v := range NtArray {
		if nt == v {
			return nt, nil
		}
	}
	return -1, errors.New("invalid note type")
}

func parsePoints(v string) [2]int {
	var points [2]int
	s := strings.Split(v, `:`)
	for i := range points {
		points[i] = tools.Atoi(s[i])
	}
	return points
}
