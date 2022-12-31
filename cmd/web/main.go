package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene/play"
)

var ls = js.Global().Get("localStorage")

func loadString(k string) string {
	v := ls.Call("getItem", k).String()
	// fmt.Printf("%s: %s\n", k, v)
	return v
}
func loadFloat(k string) float64 {
	f, err := strconv.ParseFloat(loadString(k), 64)
	if err != nil {
		fmt.Println(err)
	}
	return f
}
func loadInt(k string) int { return int(loadFloat(k)) }
func loadStrings(k string) []string {
	return strings.Split(loadString(k), ",")
}
func loadSettings() {
	if v := loadFloat("volumeMusic"); v != 0 {
		mode.S.VolumeMusic = v
	}
	if v := loadFloat("volumeSound"); v != 0 {
		mode.S.VolumeSound = v
	}
	if v := loadFloat("brightness"); v != 0 {
		mode.S.BackgroundBrightness = v
	}
	if v := loadFloat("offset"); v != 0 {
		mode.S.Offset = int64(v)
	}
	if v := loadFloat("speedPiano"); v != 0 {
		piano.S.SpeedScale = v
	}
	if v := loadFloat("speedDrum"); v != 0 {
		drum.S.SpeedScale = v
	}

	if v := loadStrings("keyPiano4"); len(v) != 0 {
		piano.S.KeySettings[4] = mode.NormalizeKeys(v)
		// for i, k := range piano.S.KeySettings[4] {
		// 	fmt.Printf("index: %d key: %s\n", i, k)
		// }
	}
	if v := loadStrings("keyPiano7"); len(v) != 0 {
		piano.S.KeySettings[7] = mode.NormalizeKeys(v)
	}
	if v := loadStrings("keyDrum4"); len(v) != 0 {
		drum.S.KeySettings[4] = mode.NormalizeKeys(v)
	}
}
func saveSettings() {
	ls.Call("setItem", "volumeMusic", mode.S.VolumeMusic)
	ls.Call("setItem", "volumeSound", mode.S.VolumeSound)
	ls.Call("setItem", "brightness", mode.S.BackgroundBrightness)
	ls.Call("setItem", "offset", mode.S.Offset)
	ls.Call("setItem", "speedPiano", piano.S.SpeedScale)
	ls.Call("setItem", "speedDrum", drum.S.SpeedScale)
}
func main() {
	loadSettings()
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case <-ticker.C:
				saveSettings()
			}
		}
	}()
	c := play.Chart{
		ParentSetId:  loadInt("setID"),
		OsuMode:      loadInt("osuMode"),
		CS:           loadInt("cs"),
		OsuFile:      loadString("osuFile"),
		ChartName:    loadString("chartName"),
		DownloadPath: loadString("downloadPath"),
	}
	fmt.Printf("%+v\n", c)

	s, err := play.NewScene(c)
	if err != nil {
		js.Global().Call("alert", err.Error())
	}
	g := &gosu.Game{
		IsWeb: true,
		Scene: s,
	}
	if err := ebiten.RunGame(g); err != nil {
		js.Global().Call("alert", err.Error())
	}
}

// c = play.Chart{
// 	ParentSetId:  1699001,
// 	OsuMode:      3,
// 	CS:           4,
// 	OsuFile:      "Pecorine (CV M.A.O), Kokkoro (CV Ito Miku), Karyl (CV Tachibana Rika) - Yes! Precious Harmony! (ML-ysg) [Easy].osu",
// 	DownloadPath: "/d/1699001",
// }
// fmt.Printf("%+v\n", c)

// c := play.Chart{
// 	ParentSetId:  320687,
// 	Mode:         3,
// 	CS:           4,
// 	OsuFile:      "Jason Hayes - War (Feerum) [Easy].osu",
// 	DownloadPath: "/d/677891",
// }
