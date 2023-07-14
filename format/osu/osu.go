package osu

import (
	"bufio"
	"bytes"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func Parse(dat []byte) (*Format, error) {
	o := &Format{
		General: General{
			PreviewTime:      -1,
			Countdown:        1,
			SampleSet:        "Normal",
			StackLeniency:    0.7,
			StoryFireInFront: true,
			OverlayPosition:  "NoChange",
		},
		Events:       make([]Event, 0),
		TimingPoints: make([]TimingPoint, 0),
		HitObjects:   make([]HitObject, 0),
	}
	dat = bytes.ReplaceAll(dat, []byte("\r\n"), []byte("\n"))

	var section string
	for _, l := range bytes.Split(dat, []byte("\n")) {
		l = bytes.TrimLeftFunc(l, unicode.IsSpace) // prevent trimming delimiter
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
			kv := strings.SplitN(line, `: `, 2)
			if len(kv) < 2 {
				continue
			}
			kv[1] = strings.TrimRightFunc(kv[1], unicode.IsSpace)
			switch kv[0] {
			case "AudioFilename":
				o.General.AudioFilename = kv[1]
			case "AudioLeadIn":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.AudioLeadIn = int(f)
			case "AudioHash":
				o.General.AudioHash = kv[1]
			case "PreviewTime":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.PreviewTime = int(f)
			case "Countdown":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.Countdown = i
			case "SampleSet":
				o.General.SampleSet = kv[1]
			case "StackLeniency":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.StackLeniency = f
			case "Mode":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.Mode = i
			case "LetterboxInBreaks":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.LetterboxInBreaks = b
			case "StoryFireInFront":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.StoryFireInFront = b
			case "UseSkinSprites":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.UseSkinSprites = b
			case "AlwaysShowPlayfield":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.AlwaysShowPlayfield = b
			case "OverlayPosition":
				o.General.OverlayPosition = kv[1]
			case "SkinPreference":
				o.General.SkinPreference = kv[1]
			case "EpilepsyWarning":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.EpilepsyWarning = b
			case "CountdownOffset":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.CountdownOffset = int(f)
			case "SpecialStyle":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.SpecialStyle = b
			case "WidescreenStoryboard":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.WidescreenStoryboard = b
			case "SamplesMatchPlaybackRate":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.General.SamplesMatchPlaybackRate = b
			}
		case "Editor":
			kv := strings.SplitN(line, `: `, 2)
			if len(kv) < 2 {
				continue
			}
			// Number-only sections may be trimmed both space.
			// kv[1] = strings.TrimRightFunc(kv[1], unicode.IsSpace)
			kv[1] = strings.TrimSpace(kv[1])
			switch kv[0] {
			case "Bookmarks":
				slice := make([]int, 0)
				for _, s := range strings.Split(kv[1], ",") {
					i, err := strconv.Atoi(s)
					if err != nil {
						// return o, fmt.Errorf("error at %s: %s", line, err)
						continue
					}
					slice = append(slice, i)
				}
				o.Editor.Bookmarks = slice
			case "DistanceSpacing":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Editor.DistanceSpacing = f
			case "BeatDivisor":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Editor.BeatDivisor = f
			case "GridSize":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Editor.GridSize = int(f)
			case "TimelineZoom":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Editor.TimelineZoom = f
			}
		case "Metadata":
			kv := strings.SplitN(line, `:`, 2)
			if len(kv) < 2 {
				continue
			}
			kv[1] = strings.TrimRightFunc(kv[1], unicode.IsSpace)
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
				o.Metadata.Tags = strings.Split(kv[1], " ")
			case "BeatmapID":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Metadata.BeatmapID = i
			case "BeatmapSetID":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Metadata.BeatmapSetID = i
			}
		case "Difficulty":
			kv := strings.SplitN(line, `:`, 2)
			if len(kv) < 2 {
				continue
			}
			// Number-only sections may be trimmed both space.
			// kv[1] = strings.TrimRightFunc(kv[1], unicode.IsSpace)
			kv[1] = strings.TrimSpace(kv[1])
			switch kv[0] {
			case "HPDrainRate":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Difficulty.HPDrainRate = f
			case "CircleSize":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Difficulty.CircleSize = f
			case "OverallDifficulty":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Difficulty.OverallDifficulty = f
			case "ApproachRate":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Difficulty.ApproachRate = f
			case "SliderMultiplier":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Difficulty.SliderMultiplier = f
			case "SliderTickRate":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return o, fmt.Errorf("error at %s: %s", line, err)
				}
				o.Difficulty.SliderTickRate = f
			}
		case "Events":
			e, err := newEvent(line)
			if err != nil {
				// return o, fmt.Errorf("error at %s: %s", line, err)
				continue
			}
			o.Events = append(o.Events, e)
		case "TimingPoints":
			tp, err := newTimingPoint(line)
			if err != nil {
				return o, fmt.Errorf("error at %s: %s", line, err)
			}
			o.TimingPoints = append(o.TimingPoints, tp)
		case "Colours":
			kv := strings.Split(line, ` : `)
			rgb := newRGB(kv[1])
			switch kv[0] {
			case "Combo1":
				o.Colours.Combos[0] = rgb
			case "Combo2":
				o.Colours.Combos[1] = rgb
			case "Combo3":
				o.Colours.Combos[2] = rgb
			case "Combo4":
				o.Colours.Combos[3] = rgb
			case "Combo5":
				o.Colours.Combos[4] = rgb
			case "Combo6":
				o.Colours.Combos[5] = rgb
			case "Combo7":
				o.Colours.Combos[6] = rgb
			case "Combo8":
				o.Colours.Combos[7] = rgb
			case "SliderTrackOverride":
				o.Colours.SliderTrackOverride = rgb
			case "SliderBorder":
				o.Colours.SliderBorder = rgb
			}
		case "HitObjects":
			ho, err := newHitObject(line)
			if err != nil {
				return o, fmt.Errorf("error at %s: %s", line, err)
			}
			o.HitObjects = append(o.HitObjects, ho)
		}
	}
	return o, nil
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

func newRGB(s string) color.RGBA {
	var rgb color.RGBA
	for i, c := range strings.Split(s, `,`) {
		f, err := strconv.ParseFloat(c, 64)
		if err != nil {
			f = 0
		}
		switch i {
		case 0:
			rgb.R = uint8(f)
		case 1:
			rgb.G = uint8(f)
		case 2:
			rgb.B = uint8(f)
		}
	}
	rgb.A = 255
	return rgb
}

const (
	ModeStandard = iota
	ModeTaiko
	ModeCatch
	ModeMania
)
const ModeOsu = ModeStandard
const ModeDefault = ModeStandard

func Mode(path string) (int, int) {
	const modeError = -1

	f, err := os.Open(path)
	if err != nil {
		return modeError, 0
	}
	defer f.Close()
	var (
		mode     int
		keyCount int
	)
	scanner := bufio.NewScanner(f) // Splits on newlines by default.
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Mode: ") {
			vs := strings.Split(scanner.Text(), "Mode: ")
			if len(vs) != 2 {
				return ModeDefault, 0 // Blank goes default mode.
			}
			v, err := strconv.Atoi(vs[1])
			if err != nil {
				return modeError, 0
			}
			mode = v
		}
		if strings.HasPrefix(scanner.Text(), "CircleSize:") {
			vs := strings.Split(scanner.Text(), "CircleSize:")
			if len(vs) != 2 {
				return mode, 0
			}
			v, err := strconv.Atoi(vs[1])
			if err != nil {
				return mode, 0
			}
			keyCount = v
			return mode, keyCount
		}
	}
	if err := scanner.Err(); err != nil {
		return modeError, 0
	}
	return ModeDefault, 0
}

func parseInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return int(f), nil
	}
	return i, nil
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
