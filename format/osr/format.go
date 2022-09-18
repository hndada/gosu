package osr

type Format struct {
	GameMode    int8
	GameVersion int32
	BeatmapMD5  string
	PlayerName  string
	ReplayMD5   string
	Num300      int16
	Num100      int16
	Num50       int16
	NumGeki     int16
	NumKatu     int16
	NumMiss     int16
	Score       int32
	Combo       int16
	FullCombo   bool
	ModsBits    int32
	LifeBar     string
	TimeStamp   int64
	ReplayData  []Action
	OnlineID    int64
	// AddMods     float64 // indirect data of accuracy at Target Practice. It exists only when the mod is on.
}

type Action struct {
	W int64   // elapsed time since last action
	X float64 // mouse cursor; pressed keys at mania. The least bit refers to state of the leftmost column and so on.
	Y float64 // mouse cursor
	Z int64   // pressed keys at standard
}

// In normal replay, first 2 data are dummy with x = 256 and y = -500
// I assume it is for setting time offset: -1.
// In auto replay, first data is blank action.
// func (f Format) IsAuto() bool {
// 	const (
// 		x = 256
// 		y = -500
// 	)
// 	if len(f.ReplayData) < 2 {
// 		return true
// 	}
// 	a0, a1 := f.ReplayData[0], f.ReplayData[1]
// 	if a0.X == 0 && a0.Y == 0 {
// 		return true
// 	} else if a0.X == x && a0.Y == y && a1.X == x && a1.Y == y {
// 		return false
// 	}
// 	panic("no reach")
// }

// // Last action data is dummy which is for random seed.
// func (f Format) TrimmedActions() []Action {
// 	if f.IsAuto() {
// 		return f.ReplayData
// 	}
// 	return f.ReplayData[2 : len(f.ReplayData)-1]
// }

// BufferTime returns the amount of time of waiting before music start when playing a chart.
// func (f Format) BufferTime() int64 {
// 	if f.IsAuto() {
// 		return 0
// 	}
// 	a0, a1, a2 := f.ReplayData[0], f.ReplayData[1], f.ReplayData[2]
// 	return a0.W + a1.W + a2.W // Must be 0 - 1 - (actual buffer time)
// }
