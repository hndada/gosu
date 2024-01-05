package draws

import "image/color"

const (
	Top = iota
	Right
	Bottom
	Left
)

const (
	TopLeft = iota
	TopRight
	BottomRight
	BottomLeft
)

// type Radius [2]Length

type Border struct {
	// Margins [4]Length
	// Style    BorderStyle
	// Widths   [4]float64
	Colors   [4]color.Color
	Collapse bool
	Spacings [4]Length
	// Radiuses [4]Radius
	// Paddings [4]Length
}

// type BorderStyle int

// const (
// 	BorderNone BorderStyle = iota
// 	BorderSolid
// 	BorderDashed
// 	BorderDotted
// 	BorderDouble
// 	BorderGroove
// 	BorderRidge
// 	BorderInset
// 	BorderOutset
// )
