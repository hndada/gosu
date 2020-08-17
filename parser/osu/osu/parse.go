package osu

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
)

// todo: 자동 생성 코드 검사, 수동 코드 추가
// todo: 미구현 파싱 추가
// todo: put flag: force parse (ignore error)
func Parse(path string) (*FormatOsu, error) {
	var o FormatOsu
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return &o, err
	}
	dat = bytes.ReplaceAll(dat, []byte("\r\n"), []byte("\n"))

	o.Events = make([]Event, 0)
	o.TimingPoints = make([]TimingPoint, 0)
	o.HitObjects = make([]HitObject, 0)
	var section string
	for _, l := range bytes.Split(dat, []byte("\n")) {
		l = bytes.TrimSpace(l)
		line := string(l)
		if isPass(line) {
			continue
		}
		if isSection(line) {
			section = strings.Trim(line, "[]")
			continue
		}
		switch section {
		case "General":
			kv := strings.Split(line, `:(space)`)
			switch kv[0] {
			case "AudioFilename":
				o.General.AudioFilename = kv[1]
			case "AudioLeadIn":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.AudioLeadIn = i
			case "AudioHash":
				o.General.AudioHash = kv[1]
			case "PreviewTime":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.PreviewTime = i
			case "Countdown":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.Countdown = i
			case "SampleSet":
				o.General.SampleSet = kv[1]
			case "StackLeniency":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.General.StackLeniency = f
			case "Mode":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.Mode = i
			case "LetterboxInBreaks":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.LetterboxInBreaks = b
			case "StoryFireInFront":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.StoryFireInFront = b
			case "UseSkinSprites":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.UseSkinSprites = b
			case "AlwaysShowPlayfield":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.AlwaysShowPlayfield = b
			case "OverlayPosition":
				o.General.OverlayPosition = kv[1]
			case "SkinPreference":
				o.General.SkinPreference = kv[1]
			case "EpilepsyWarning":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.EpilepsyWarning = b
			case "CountdownOffset":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.CountdownOffset = i
			case "SpecialStyle":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.SpecialStyle = b
			case "WidescreenStoryboard":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.WidescreenStoryboard = b
			case "SamplesMatchPlaybackRate":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.SamplesMatchPlaybackRate = b
			}
		case "Editor":
			kv := strings.Split(line, `:(space)`)
			switch kv[0] {
			case "Bookmarks":
				slice := make([]int, 0)
				for _, s := range strings.Split(kv[1], ",") {
					i, err := strconv.Atoi(s)
					if err != nil {
						return &o, err
					}
					slice = append(slice, i)
				}
				o.Editor.Bookmarks = slice
			case "DistanceSpacing":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Editor.DistanceSpacing = f
			case "BeatDivisor":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Editor.BeatDivisor = f
			case "GridSize":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.Editor.GridSize = i
			case "TimelineZoom":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Editor.TimelineZoom = f
			}
		case "Metadata":
			kv := strings.Split(line, `:`)
			switch kv[0] {
			case "Title":
				o.Metadata.Title = kv[1]
			case "TitleUnicode":
				o.Metadata.TitleUnicode = kv[1]
			case "Artist":
				o.Metadata.Artist = kv[1]
			case "ArtistUnicode":
				o.Metadata.ArtistUnicode = kv[1]
			case "Creator":
				o.Metadata.Creator = kv[1]
			case "Version":
				o.Metadata.Version = kv[1]
			case "Source":
				o.Metadata.Source = kv[1]
			case "Tags":
				slice := make([]string, 0)
				for _, s := range strings.Split(kv[1], " ") {
					slice = append(slice, s)
				}
				o.Metadata.Tags = slice
			case "BeatmapID":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.Metadata.BeatmapID = i
			case "BeatmapSetID":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.Metadata.BeatmapSetID = i
			}
		case "Difficulty":
			kv := strings.Split(line, `:`)
			switch kv[0] {
			case "HPDrainRate":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.HPDrainRate = f
			case "CircleSize":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.CircleSize = f
			case "OverallDifficulty":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.OverallDifficulty = f
			case "ApproachRate":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.ApproachRate = f
			case "SliderMultiplier":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.SliderMultiplier = f
			case "SliderTickRate":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.SliderTickRate = f
			}
		
		
		
		
		case "Events":
			continue
		case "TimingPoints":
			// tp, err := newTimingPoint(line)
			// if err != nil {
			// 	return &o, err
			// }
			// o.TimingPoints = append(o.TimingPoints, tp)
		}
	}
	return &o, nil
}
func isPass(line string) bool {
	return len(line) == 0 || len(line) >= 2 && line[:2] == "//"
}
func isSection(line string) bool {
	if len(line) == 0 {
		return false
	}
	return string(line[0]) == "[" && string(line[len(line)-1]) == "]"
}
