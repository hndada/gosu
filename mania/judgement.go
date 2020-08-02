package mania

type Judgement struct {
	Name    string
	Value   float64
	Penalty float64
	Window  int64
}

var Judgements = [5]Judgement{
	{"KOOL", 16 / 16, 0, 16},
	{"COOL", 15 / 16, 0, 40},
	{"GOOD", 10 / 16, 4, 70},
	{"BAD", 4 / 16, 10, 100},
	{"MISS", 0, 25, 150},
}