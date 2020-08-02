package lv

import (
	"fmt"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/mania"
	"github.com/hndada/gosu/tools"
	"log"
)

// ppy 방식처럼, 구간 내 최고 strain을 잡아야 할까?
// 우선 chord 알고리즘 먼저 고쳐보자
const (
	diffWeightDecay = 0.90
	sectionLength   = 800
)

func process(path string){
	base, err := beatmap.ParseBeatmap(path)
	if err != nil {
		log.Fatal(err)
	}
	b := &mania.Beatmap{}
	b.Beatmap = base

	keymode := int(b.Difficulty["CircleSize"].(float64))
	mods := game.Mods{1, false, false}
	ns, err:= mania.NewNotes(b.HitObjects, keymode, mods)
	if err!=nil {log.Fatal(err) }

	mania.SortNotes(ns)
	CalcStrain(ns, keymode)
	CalcStamina(ns, keymode)

	// b.Notes=ns
	// b.Keymode=keymode
	CalcLV(ns)
}

func CalcLV(ns []mania.Note) float64 {
	if len(ns) == 0 {
		return 0
	}
	sectionCounts := (ns[len(ns)-1].Time-ns[0].Time) / sectionLength
	sectionEndTime := sectionLength + ns[0].Time

	var aggregate float64
	aggregates := make([]float64, 0, sectionCounts)
	for _, n := range ns {
		for n.Time > sectionEndTime {
			aggregates = append(aggregates, aggregate)
			aggregate = 0.0
			sectionEndTime += sectionLength
		}
		aggregate += n.Strain+n.Stamina
	}
	switch {
	case len(aggregates)==sectionCounts:
	case len(aggregates)==sectionCounts-1:
		aggregates = append(aggregates, aggregate)
	default:
		fmt.Println(len(aggregates), sectionCounts)
		panic("section count mismatch")
	}

	lv1 := tools.WeightedSum(aggregates, diffWeightDecay)
	// newSectionCounts:=sectionCounts
	// for aggregates[newSectionCounts-1]<3 {newSectionCounts--}

	// q1:=newSectionCounts/4
	// fmt.Println(q1)
	// lv2:=0.9*tools.WeightedSum(aggregates[:q1], diffWeightDecay)+
		// 0.1*tools.WeightedSum(aggregates[q1:newSectionCounts], diffWeightDecay)
	fmt.Printf("%.2f\\n", lv1*0.031)
	return lv1
}