package gosu

import "time"

// Todo: SceneResult
// Todo: implement Replay
type Result struct {
	MD5        [16]byte  // MD5 for raw chart file. md5.Size = 16
	PlayedTime time.Time // Finish time of playing.

	ScoreFactors   [3]float64 // Retrieved from the chart.
	Scores         [4]float64
	JudgmentCounts []int
	MaxCombo       int
	// FlowMarks      []float64 // Length is around 100 ~ 200.
	// KeyLogs []KeyLog // Entire timed-log key strokes.
}
