package draws

type Unit int

const (
	Pixel Unit = iota
	Percent
	RootPercent

	// Em // A unit of length in typography.
	// RootEm

	// Extra stands for additional size to parent's size.
	Extra
)

type Length struct {
	Value float64
	Unit  Unit
}
