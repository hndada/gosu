package draws

type Unit int

const (
	Pixel Unit = iota
	Percent
	RootPercent
	// Em // A unit of length in typography.
	// RootEm
)

type Length struct {
	Value float64
	Unit  Unit
}
