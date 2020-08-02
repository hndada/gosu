package mania

import "math"

// od 0, sv2, nomod로 몇개 테스트
// -> No Mods except Rate one, thus no need to separate HitBonus (someday)


var Judgements [6]Judgement

// only NoMod for simple test plz
func init() {
	Values := [6]int{320, 300, 200, 100, 50, 0}
	BonusValues := [6]int{32, 32, 16, 8, 4, 0}
	Bonuses := [6]int{2, 1, 0, 0, 0, 0}
	Punishments := [6]int{0, 0, 8, 24, 44, 200}
	for i := range Judgements {
		Judgements[i] = Judgement{
			Values[i], BonusValues[i], Bonuses[i], Punishments[i],
		}
	}
}

// 굿나면 일단 굿 자체에서 절반 이상 까이고 시작
// 1굿 자체는 큰 영향 없음
// 단 미스의 경우, 미스 뒤 굿이라도 나면 25%.
// closure
func BonusScore(j Judgement) {
	var lastBonus float64
	bonus = lastBonus + j.Bonus - j.Punishment
	j.BonusValue * math.Sqrt(bonus) / 320
}

func UnitScore(c int) float64 { return MaxScore / float64(c) }