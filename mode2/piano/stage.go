package piano

// Ground, Hint, KeyButton
// Todo: generate hint image?
type StageConfig struct {
	KeyboardMapping map[int][]string
	GroundOpacity   float64
	HintHeight      float64
}

func NewStageConfig() StageConfig {
	return StageConfig{
		KeyboardMapping: map[int][]string{
			4:  {"D", "F", "J", "K"},
			5:  {"D", "F", "Space", "J", "K"},
			6:  {"S", "D", "F", "J", "K", "L"},
			7:  {"S", "D", "F", "Space", "J", "K", "L"},
			8:  {"A", "S", "D", "F", "Space", "J", "K", "L"},
			9:  {"A", "S", "D", "F", "Space", "J", "K", "L", "Semicolon"},
			10: {"A", "S", "D", "F", "V", "N", "J", "K", "L", "Semicolon"},
		},
		GroundOpacity: 0.8,
		HintHeight:    0.05,
	}
}

type StageArgs struct {
	KeyboardMapping map[int][]string
	GroundWidth     float64
	GroundHeight    float64
	GroundOpacity   float64
	HintHeight      float64
	HintPositionY   float64
}

func (cfg Config) StageArgs() StageArgs {
	return StageArgs{
		KeyboardMapping: cfg.KeyboardMapping,
		GroundWidth:     cfg.StageWidth(),
		GroundHeight:    cfg.ScreenSize.Y,
		GroundOpacity:   cfg.GroundOpacity,
		HintHeight:      cfg.HintHeight * cfg.ScreenSize.Y,
		HintPositionY:   cfg.HitPositionY,
	}
}
