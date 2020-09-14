package game

// Game version:
// [Major build number].[Minor build number].[Revision].[Package]
// i.e. 1.0.15.2 (ip style)
// ChartMD5, ChartPlayMD5는 mods에 관계 없이 고정된 값을 가져야 함
const MaxScore = 1e6

type BaseScore struct {
	// GameMode     uint8
	GameVersion uint32
	// LevelVersion uint32 // for level
	ChartMD5     [16]byte
	ChartPlayMD5 [16]byte

	PlayerName string
	TimeStamp  int64
	Score      int64
	Combo      int32
}

// 인터페이스를 어케 쓰는지 뭔가 알았다
// todo: mods
type Score interface {
	JudgeCounts() []int64
	IsFullCombo() bool
	IsPerfect() bool
}
