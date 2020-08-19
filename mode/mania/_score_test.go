package mania

// // closure
// func bonusScore(j Judgement) {
// 	var lastBonus float64
// 	bonus = lastBonus + j.Bonus - j.Punishment
// 	j.BonusValue * math.Sqrt(bonus) / 320
// }
//
// // od 0, sv2
// func TestProcessLegacyScore(t *testing.T) {
// 	var o *Chart
// 	var Judgements [6]Judgement
// 	Values := [6]int{320, 300, 200, 100, 50, 0}
// 	BonusValues := [6]int{32, 32, 16, 8, 4, 0}
// 	Bonuses := [6]int{2, 1, 0, 0, 0, 0}
// 	Punishments := [6]int{0, 0, 8, 24, 44, 200}
// 	for i := range Judgements {
// 		Judgements[i] = Judgement{
// 			Values[i], BonusValues[i], Bonuses[i], Punishments[i],
// 		}
// 	}
// 	scoreUnit := MaxScore / float64(len(o.Notes))
// }
