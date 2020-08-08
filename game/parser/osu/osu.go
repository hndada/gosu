package osu

import (
	"bytes"
	"errors"
	"github.com/hndada/gosu/tools"
	"io/ioutil"
	"sort"
	"strings"
)

type OSU struct {
	General    tools.Info
	Editor     tools.Info
	Metadata   tools.Info
	Difficulty tools.Info
	Image      Event
	Video      Event
	Events     []string
	Colours    tools.Info

	TimingPoints []TimingPoint
	HitObjects   []HitObject
}

type Event struct {
	StartTime int64
	Filename  string
	XOffset   int
	YOffset   int
}

func NewOSU(path string) (OSU, error) { // todo: return pointer
	var o OSU
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return o, err
	}
	dat = bytes.ReplaceAll(dat, []byte("\r\n"), []byte("\n"))
	lines := make([]string, 0)
	for _, l := range bytes.Split(dat, []byte("\n")) {
		lines = append(lines, string(l))
	}

	sectionLens, err := getSectionLength(lines)
	if err != nil {
		return o, err
	}

	general, editor, metadata, difficulty := make(tools.Info), make(tools.Info), make(tools.Info), make(tools.Info)
	events := make([]string, 0, sectionLens["Events"])
	timingPoints := make([]TimingPoint, 0, sectionLens["TimingPoints"])
	colours := make(tools.Info)
	hitObjects := make([]HitObject, 0, sectionLens["HitObjects"])

	var l, section string
	var kv, vs []string
	for _, line := range lines {
		l = strings.TrimSpace(line)
		if len(l) == 0 || len(l) >= 2 && l[:2] == "//" {
			continue
		}
		switch {
		case isSection(line):
			section = strings.Trim(line, "[]")
		default:
			switch section {
			case "General":
				kv = strings.Split(line, `: `)
				switch kv[0] {
				case "AudioFilename", "AudioHash", "SampleSet", "OverlayPosition", "SkinPreference":
					general.PutStr(kv)
				case "AudioLeadIn", "PreviewTime", "Countdown", "CountdownOffset":
					general.PutInt(kv)
				case "Mode":
					general.PutInt(kv)
					// if err = general.PutInt(kv); err != nil {
					// 	return o, errors.New("invalid mode")
					// }
				case "StackLeniency":
					general.PutF64(kv)
				default:
					general.PutBool(kv)
				}
			case "Editor":
				kv = strings.Split(line, `: `)
				switch kv[0] {
				case "Bookmarks":
					editor.PutIntSlice(kv)
				case "GridSize":
					editor.PutInt(kv)
				default:
					editor.PutF64(kv)
				}
			case "Metadata":
				kv = strings.Split(line, `:`)
				switch kv[0] {
				case "Tags":
					metadata.PutStrSlice(kv)
				case "BeatmapID", "BeatmapSetID":
					metadata.PutInt(kv)
				default:
					metadata.PutStr(kv)
				}
			case "Difficulty":
				kv = strings.Split(line, `:`)
				difficulty.PutF64(kv)
			case "Events":
				// 0,0,filename,xOffset,yOffset
				vs = strings.Split(line, `,`)
				var xOffset, yOffset int
				startTime, _ := tools.Atoi(vs[1])
				filename := strings.Trim(vs[2], `"`)
				if len(vs) > 3 {
					xOffset, _ = tools.Atoi(vs[3])
					yOffset, _ = tools.Atoi(vs[4])
				}
				event := Event{int64(startTime), filename, xOffset, yOffset}
				switch vs[0] {
				case "0":
					o.Image = event
				case "1", "Video":
					o.Video = event
				default:
					events = append(events, line)
				}
			case "TimingPoints":
				timingPoint, err := parseTimingPoint(line)
				if err != nil {
					return o, err
				}
				timingPoints = append(timingPoints, timingPoint)
			case "Colours":
				kv = strings.Split(line, ` : `)
				colours.PutIntSlice(kv)
			case "HitObjects":
				hitObject, err := parseHitObject(line)
				if err != nil {
					return o, err
				}
				hitObjects = append(hitObjects, hitObject)
			}
		}
	}
	o.General = general
	o.Editor = editor
	o.Metadata = metadata
	o.Difficulty = difficulty
	o.Events = events
	o.Colours = colours

	sort.Slice(timingPoints, func(i, j int) bool {
		return timingPoints[i].Time < timingPoints[j].Time
	})
	o.TimingPoints = timingPoints

	sort.Slice(hitObjects, func(i, j int) bool {
		return hitObjects[i].StartTime < hitObjects[j].StartTime
	})
	o.HitObjects = hitObjects

	o.calcSliderEndTime()
	return o, nil
}

func isSection(line string) bool {
	if len(line) == 0 {
		return false
	}
	return string(line[0]) == "[" && string(line[len(line)-1]) == "]"
}

func getSectionLength(lines []string) (map[string]int, error) {
	// todo: yet too early to convert to []string? should i treat this with [][]byte?
	// [][]byte로 하면 아예 수정되니까 이게 나은 거 같음
	var l, section string
	var c int
	lens := make(map[string]int)

	for _, line := range lines {
		l = strings.TrimSpace(line)
		if len(l) == 0 || len(l) >= 2 && l[:2] == "//" {
			continue
		}
		switch {
		case isSection(line):
			if section != "" {
				if _, ok := lens[section]; ok {
					return lens, errors.New("duplicated sections")
				}
				lens[section] = c
			}
			c = 0
			section = strings.Trim(line, "[]")
		default:
			c++
		}
	}
	lens[section] = c // for last section
	return lens, nil
}

// todo: correctness check yet
func (o *OSU) calcSliderEndTime() {
	var duration float64
	var tPoint TimingPoint
	for i, note := range o.HitObjects {
		if note.NoteType != NtSlider {
			continue
		}
		for j := len(o.TimingPoints) - 1; j >= 0; j-- {
			tPoint = o.TimingPoints[j]
			if tPoint.Time > note.StartTime || tPoint.Uninherited {
				continue
			}
			duration = note.SliderParams.Length / tPoint.SpeedScale
			duration /= o.Difficulty["SliderMultiplier"].(float64) * 100
			o.HitObjects[i].EndTime = note.StartTime + int(duration)
		}
	}
}
