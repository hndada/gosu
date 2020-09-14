package taiko

const (
	left  = -1
	right = 1
)

// kddk
func finger(hand, color int) int {
	switch hand {
	case left:
		return 2 - color
	case right:
		return 1 + color
	}
	panic("ErrValue")
}

// hand==color in kkdd
