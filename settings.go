package gosu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func SetKeySettings(props []ModeProp, KeyConfigs []KeyConfig) {
	for _, KeyConfig := range KeyConfigs {
		mode := KeyConfig.Mode
		keys := input.NamesToKeys(KeyConfig.Keys)
		if !input.IsKeysValid(keys) {
			fmt.Printf("mapping keys are duplicated: #{names}\n")
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
			subMode, err := strconv.Atoi(strings.Replace(mode, "Key", "", 1))
			if err != nil {
				fmt.Printf("error at loading key settings %s: %v", KeyConfig, err)
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
