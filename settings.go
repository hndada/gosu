package gosu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/hndada/gosu/input"
)

var (
	MusicRoot   = "music"
	WindowSizeX = 1600
	WindowSizeY = 900
)
var (
	// TPS supposed to be multiple of 1000, since only one speed value
	// goes passed per Update, while unit of TransPoint's time is 1ms.
	// TPS affects only on Update(), not on Draw().
	TPS int = 1000 // TPS should be 1000 or greater.

	CursorScale        float64 = 0.1
	ChartInfoBoxWidth  float64 = 450
	ChartInfoBoxHeight float64 = 50
	ChartInfoBoxShrink float64 = 0.15
	chartInfoBoxshrink float64 = ChartInfoBoxWidth * ChartInfoBoxShrink
	chartItemBoxCount  int     = int(screenSizeY/ChartInfoBoxHeight) + 2 // Gives some margin.

	ScoreScale    float64 = 0.65
	ScoreDigitGap float64 = 0
	MeterWidth    float64 = 4 // The number of pixels per 1ms.
	MeterHeight   float64 = 50
)

type Size2D struct {
	Width  int `json:"Width"`
	Height int `json:"Height"`
}

type KeyConfig struct {
	Mode string   `json:"Mode"`
	Keys []string `json:"Keys"`
}

type Config struct {
	MusicRoot  string      `json:"MusicRoot"`
	WindowSize Size2D      `json:"WindowSize"`
	KeyConfigs []KeyConfig `json:"KeyConfig"`
}

// Todo: reset all tick-dependent variables.
// They are mostly at drawer.go or play.go, settings.go
// Keyword: TimeToTick
func SetTPS() {}

func LoadSettings() (config Config) {
	// TODO: Overwrite user's custom settings
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Printf("error: #{err}")
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}
	return
}

// Temporary function.
func SetKeySettings(props []ModeProp) {
	data, err := os.ReadFile("keys.txt")
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	for _, line := range strings.Split(string(data), "\n") {
		if len(line) == 0 {
			continue
		}
		if len(line) >= 1 && line[0] == '#' {
			continue
		}
		if len(line) >= 2 && line[0] == '/' && line[1] == '/' {
			continue
		}
		kv := strings.Split(line, ": ")
		mode := kv[0]
		names := strings.Split(kv[1], ", ")
		for i, name := range names {
			names[i] = strings.TrimSpace(name)
		}
		keys := input.NamesToKeys(names)
		if !input.IsKeysValid(keys) {
			fmt.Printf("mapping keys are duplicated: %v\n", names)
			continue
		}
		switch mode {
		case "Drum", "drum":
			for _, prop := range props {
				if strings.Contains(strings.ToLower(prop.Name), "drum") {
					prop.KeySettings[4] = keys
					break
				}
			}
		default:
			subMode, err := strconv.Atoi(mode)
			if err != nil {
				fmt.Printf("error at loading key settings %s: %v", line, err)
				continue
			}
			for _, prop := range props {
				if strings.Contains(strings.ToLower(prop.Name), "piano") {
					prop.KeySettings[subMode] = keys
					break
				}
			}
		}
	}
}
