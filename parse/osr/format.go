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
