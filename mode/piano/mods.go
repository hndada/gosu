package piano

type ScratchMode int

const (
	NoScratch ScratchMode = iota
	LeftScratch
	RightScratch
)

type Mods struct{}
