package mania

import (
	"fmt"
	"github.com/hndada/gosu/mode"
)

// todo: score, level 다 정리되고 나서 internal/tools 정리하겠음
// ppy 방식처럼, 구간 내 최고 strain을 잡아야 할까?
// 우선 chord 알고리즘 먼저 고쳐보자
const (
	diffWeightDecay = 0.90
	sectionLength   = 800
)

// Difficulty: 난이도 관련 전반적인 값들
// Level: difficulty의 대푯값
func (c *Chart) CalcDifficulty() {
	if len(c.Notes) == 0 {
		return
	}
	sectionCounts := int(c.EndTime()-c.Notes[0].Time) / sectionLength // independent of note offset
	sectionEndTime := sectionLength + c.Notes[0].Time

	var d float64
	ds := make([]float64, 0, sectionCounts)
	for _, n := range c.Notes {
		for n.Time >= sectionEndTime {
			ds = append(ds, d)
			d = 0.0 // todo: stamina를 여기서 적용해보는건? 5단계로 감소되는 stamina들
			sectionEndTime += sectionLength
		}
		d += n.strain + n.stamina // n.read someday
	}

	if len(ds) != sectionCounts {
		fmt.Println(len(ds), sectionCounts)
		panic("section count mismatch")
	}

	c.Level = mode.WeightedSum(ds, diffWeightDecay)
	// newSectionCounts:=sectionCounts
	// for ds[newSectionCounts-1]<3 {newSectionCounts--}

	// q1:=newSectionCounts/4
	// fmt.Println(q1)
	// lv2:=0.9*tools.WeightedSum(ds[:q1], diffWeightDecay)+
	// 0.1*tools.WeightedSum(ds[q1:newSectionCounts], diffWeightDecay)
}
