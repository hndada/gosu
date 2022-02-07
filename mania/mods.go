package mania

type Mods struct {
	TimeRate float64
	Mirror   bool
	Pitch    bool
}

func NewMods() Mods {
	m := Mods{
		TimeRate: 1,
	}
	return m
}
