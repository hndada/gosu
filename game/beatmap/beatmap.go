package beatmap

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io/ioutil"
	"sort"
	"strings"
)

// todo: colours, event
type Beatmap struct {
	Md5          [md5.Size]byte
	General      Info
	Editor       Info
	Metadata     Info
	Difficulty   Info
	TimingPoints []TimingPoint
	HitObjects   []HitObject

	// HitWindows    map[string]float64 // all maps will have same value
	// Curves        map[string][]tools.Segment
	Level         float64
	OldStarRating float64
}

func ParseBeatmap(path string) (Beatmap, error) {
	var beatmap Beatmap
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return beatmap, err
	}
	beatmap.Md5 = md5.Sum(dat)

	dat = bytes.ReplaceAll(dat, []byte("\r\n"), []byte("\n"))
	lines := strings.Split(string(dat), "\n")
	sectionLens, err := getSectionLength(lines)
	if err != nil {
		return beatmap, err
	}

	general, editor, metadata, difficulty := make(Info), make(Info), make(Info), make(Info)
	timingPoints := make([]TimingPoint, sectionLens["TimingPoints"])
	hitObjects := make([]HitObject, sectionLens["HitObjects"])

	var l, section string
	var kv []string
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
					if err = general.PutInt(kv); err != nil {
						return beatmap, errors.New("invalid mode")
					}
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
			case "TimingPoints":
				timingPoint, err := parseTimingPoint(line)
				if err != nil {
					return beatmap, err
				}
				timingPoints = append(timingPoints, timingPoint)
			case "HitObjects":
				hitObject, err := parseHitObject(line)
				if err != nil {
					return beatmap, err
				}
				hitObjects = append(hitObjects, hitObject)
			}
		}
	}
	sort.Slice(timingPoints, func(i, j int) bool {
		return timingPoints[i].Time < timingPoints[j].Time
	})
	beatmap.TimingPoints = timingPoints

	sort.Slice(hitObjects, func(i, j int) bool {
		return hitObjects[i].StartTime < hitObjects[j].StartTime
	})
	beatmap.HitObjects = hitObjects

	beatmap.calcSliderEndTime()
	return beatmap, nil
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

func isSection(line string) bool {
	if len(line) == 0 {
		return false
	}
	return string(line[0]) == "[" && string(line[len(line)-1]) == "]"
}

// todo: correctness check yet
func (beatmap *Beatmap) calcSliderEndTime() {
	var duration float64
	var tPoint TimingPoint
	for i, note := range beatmap.HitObjects {
		if note.NoteType != NtSlider {
			continue
		}
		for j := len(beatmap.TimingPoints) - 1; j >= 0; j-- {
			tPoint = beatmap.TimingPoints[j]
			if tPoint.Time > note.StartTime || tPoint.Uninherited {
				continue
			}
			duration = note.SliderParams.Length / tPoint.SpeedScale
			duration /= beatmap.Difficulty["SliderMultiplier"].(float64) * 100
			beatmap.HitObjects[i].EndTime = note.StartTime + int(duration)
		}
	}
}
