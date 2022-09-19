package draws

type Origin int

const (
	OriginLeft = iota
	OriginCenter
	OriginRight
)
const (
	OriginTop = iota
	OriginMiddle
	OriginBottom
)
const (
	OriginLeftTop      Origin = iota // Default Origin.
	OriginLeftMiddle                 // e.g., Notes in Piano mode.
	OriginLeftBottom                 // e.g., back button.
	OriginCenterTop                  // e.g., drawing field.
	OriginCenterMiddle               // Most of sprite's Origin.
	OriginCenterBottom               // e.g., Meter.
	OriginRightTop                   // e.g., score.
	OriginRightMiddle                // e.g., chart info boxes.
	OriginRightBottom                // e.g., Play button.
)

func (origin Origin) PositionX() int { return int(origin) / 3 }
func (origin Origin) PositionY() int { return int(origin) % 3 }
