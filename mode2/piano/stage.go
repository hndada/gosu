package piano

type StageConfig struct {
	Widths         map[int]float64
	KeyKindWidths  [4]float64
	FieldPositionX float64
	FieldOpacity   float64
	HintHeight     float64
	HintPositionY  float64
}

func NewStageConfig(screen ScreenConfig) StageConfig {
	return StageConfig{
		Widths: map[int]float64{
			4: 0.50 * screen.Size.X,
			5: 0.55 * screen.Size.X,
			6: 0.60 * screen.Size.X,
			7: 0.65 * screen.Size.X,
			8: 0.70 * screen.Size.X,
			9: 0.75 * screen.Size.X,
		},
		KeyKindWidths: [4]float64{
			0.055 * screen.Size.X, // One
			0.055 * screen.Size.X, // Two
			0.055 * screen.Size.X, // Mid
			0.055 * screen.Size.X, // Tip
		},
		FieldPositionX: 0.50 * screen.Size.X,
		FieldOpacity:   0.8,
		HintHeight:     0.05 * screen.Size.Y,
		HintPositionY:  0.95 * screen.Size.Y,
	}
}
