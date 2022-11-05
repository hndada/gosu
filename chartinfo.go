package gosu

import (
	"fmt"
)

// ChartInfo is used at SceneSelect.
type ChartInfo struct {
	Path string
	// Mods    Mods
	// Header  ChartHeader
	ChartHeader
	Mode    int
	SubMode int
	Level   float64

	Duration   int64
	NoteCounts []int
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
	// Tags       []string // Auto-generated or User-defined
}

func (c ChartInfo) Text() string {
	switch c.Mode {
	case ModePiano4, ModePiano7:
		return fmt.Sprintf("(%dK Level %3.1f) %s [%s]", c.SubMode, c.Level, c.MusicName, c.ChartName)
	case ModeDrum:
		return fmt.Sprintf("(Level %3.1f) %s [%s]", c.Level, c.MusicName, c.ChartName)
	}
	return ""
}
func (c ChartInfo) BackgroundPath() string {
	return c.ChartHeader.BackgroundPath(c.Path)
}
func (c ChartInfo) TimeString() string {
	c.Duration /= 1000
	return fmt.Sprintf("%02d:%02d", c.Duration/60, c.Duration%60)
}
func (c ChartInfo) BPMString() string {
	return fmt.Sprintf("%.0f BPM (%.0f ~ %.0f)", c.MainBPM, c.MinBPM, c.MaxBPM)
}
func (c ChartInfo) NoteCountString() string {
	return fmt.Sprintf("◎ %d", c.NoteCounts[0])
}

// func (c ChartInfo) Download() error {
// 	const noVideo = 1
// 	u := fmt.Sprintf("%s%d?n=%d", chimuURLDownload, r.SetId, noVideo)
// 	fmt.Printf("download URL: %s\n", u)
// 	resp, err := http.Get(u)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	f, err := os.Create(filepath.Join(dir, r.Filename()))
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	_, err = io.Copy(f, resp.Body)
// }

// Background brightness at Song select: 60% (153 / 255), confirmed.
// Score box color: Gray128 with 50% transparent
// Hovered Score box color: Gray96 with 50% transparent
