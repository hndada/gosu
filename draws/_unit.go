package draws

type (
	Px = float64
	RW float64 // Relative Width
	RH float64 // Relative Height
	RX = RW    // Relative X
	RY = RH    // Relative Y
	// Em    float64 // A unit of length in typography.
)

func (opts ScreenOptions) Px(rv any) draws.Px {
	switch rv := rv.(type) {
	case draws.Px:
		return rv
	case draws.RW:
		return float64(rv) * opts.W
	case draws.RH:
		return float64(rv) * opts.H
	}
	return rv.(draws.Px)
}
