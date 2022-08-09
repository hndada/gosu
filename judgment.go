package gosu

type Judgment struct {
	Karma  float64
	Acc    float64
	Window int64
}

// A frame is 16 ~ 17ms in 60 FPS
var (
	Kool = Judgment{Karma: 0.01, Acc: 1, Window: 20}    // 1 frame
	Cool = Judgment{Karma: 0.01, Acc: 1, Window: 40}    // 2 frames
	Good = Judgment{Karma: 0.01, Acc: 0.25, Window: 70} // 4 frames
	Bad  = Judgment{Karma: 0.01, Acc: 0, Window: 100}   // 6 frames // Todo: Karma 0.01 -> 0?
	Miss = Judgment{Karma: -1, Acc: 0, Window: 150}     // 9 frames
)

var Judgments = []Judgment{Kool, Cool, Good, Bad, Miss}

func Judge(td int64) Judgment {
	if td < 0 { // Absolute value
		td *= -1
	}
	for _, j := range Judgments {
		if td <= j.Window {
			return j
		}
	}
	return Judgment{} // Returns None when the input is out of widest range
}

// func inRange(td int64, j Judgment) bool { return td < j.Window && td > -j.Window }
