package beatmap

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/hndada/gosu/game/tools"
)

//const (
//	ModeOsu = iota
//	ModeTaiko
//	ModeCatch
//	ModeMania
//)

// space 는 못잡고 tab 은 잡음
// todo: Events, Colours
type Beatmap struct {
	Md5                                   [md5.Size]byte
	General, Editor, Metadata, Difficulty Info
	TimingPoints                          []TimingPoint
	HitObjects                            []HitObject

	//HitWindows    map[string]float64 // all maps will have same value
	//Curves        map[string][]tools.Segment
	StarRating    float64
	OldStarRating float64
}

// section을 유지해야 editor로 변경하기 쉬움
func (beatmap *Beatmap) Parse(path string) error {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	beatmap.Md5 = md5.Sum(dat)
	lines, sectionInfo, err := trimLines(dat)
	if err != nil {
		return err
	}

	general, editor, metadata, difficulty := make(Info), make(Info), make(Info), make(Info)
	//colours
	//event

	var section string
	var kv []string
	for _, line := range lines {
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
						return errors.New("invalid mode")
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
				continue
				//beatmap.TimingPoints = append(beatmap.TimingPoints, parseTimingPoint(line))
			case "HitObjects":
				continue
				//beatmap.HitObjects = append(beatmap.HitObjects, parseNote(line))
			}
		}
	}
	lineIdx := sectionInfo["TimingPoints"][0]
	timingPoints := make([]TimingPoint, sectionInfo["TimingPoints"][1])
	for i := 0; i < sectionInfo["TimingPoints"][1]; i++ {
		timingPoints[i].Parse(lines[lineIdx])
		lineIdx++
	}
	sort.Slice(beatmap.TimingPoints, func(i, j int) bool {
		return beatmap.TimingPoints[i].Time < beatmap.TimingPoints[j].Time
	})

	lineIdx = sectionInfo["HitObjects"][0]
	hitObjects := make([]HitObject, sectionInfo["HitObjects"][1])
	for i := 0; i < sectionInfo["HitObjects"][1]; i++ {
		hitObjects[i].Parse(lines[lineIdx])
		lineIdx++
	}
	sort.Slice(beatmap.HitObjects, func(i, j int) bool {
		return beatmap.HitObjects[i].StartTime < beatmap.HitObjects[j].StartTime
	})
	beatmap.calcSliderEndTime()
	return nil
}

func trimLines(dat []byte) ([]string, map[string][2]int, error) {
	dat = bytes.ReplaceAll(dat, []byte("\r\n"), []byte("\n"))
	dat = bytes.ReplaceAll(dat, []byte("\n\n"), []byte("\n"))
	rawLines := bytes.Split(dat, []byte("\n"))

	var line, section string
	var i, c int
	lines := make([]string, 0, len(rawLines))
	sectionInfo := make(map[string][2]int)
	for _, byteLine := range rawLines {
		if string(byteLine[:2]) == "//" {
			continue
		}
		byteLine = bytes.TrimSpace(byteLine)
		if len(byteLine) == 0 {
			continue
		}
		line = string(byteLine)
		switch {
		case isSection(line):
			if section != "" {
				if _, ok := sectionInfo[section]; ok {
					return lines, sectionInfo, errors.New("duplicated sections")
				}
				sectionInfo[section] = [2]int{i, c - i} // idx, len
			}
			i = c
			section = strings.Trim(line, "[]")
		default:
			c++
		}
		lines = append(lines, line)
	}
	sectionInfo[section] = [2]int{i, c - i}
	return lines, sectionInfo, nil
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
