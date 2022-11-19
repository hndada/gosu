package ctrl

type Handler interface {
	Decrease()
	Increase()
}

const (
	none = iota - 1
	decrease
	increase
)
