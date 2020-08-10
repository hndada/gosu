package graphic

type MouseHandler func(event)
type MouseMoveHandler func()

type widget interface {
	render(target surface) error
	advance(elapsed float64) error

	onMouseMove(event ) bool
	onMouseEnter(event ) bool
	onMouseLeave(event ) bool
	onMouseOver(event ) bool
	onMouseButtonDown(event ) bool
	onMouseButtonUp(event ) bool
	onMouseButtonClick(event ) bool

	getPosition() (int, int)
	getSize() (int, int)
	getLayer() int
	isVisible() bool
	// isExpanding() bool
}

type widgetBase struct {
	x int
	y int
	layer int
	visible bool
	// expanding bool

	mouseEnterHandler MouseMoveHandler
	mouseLeaveHandler MouseMoveHandler
	mouseClickHandler MouseHandler
}

func (w *widgetBase) SetPosition(x, y int) {
	w.x=x
	w.y=y
}
