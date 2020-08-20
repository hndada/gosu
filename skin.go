package gosu

// todo: backward-compatibility하면 게임 스킨 파싱 버전 안 적어도 되나?
// case1: 클라가 최신, 스킨이 구식 버전 (보통)
// case2: 스킨은 최신, 클라가 구식

// skin은 only for asset/resources; 설정은 전적으로 user에게 맡긴다
// author나 license는 skin 폴더에 별도 텍스트 파일로 저장
type Skin struct {
	Name string
	// NoteImage#: // LNHead와 동일
	// NoteImage#L:
	// (NoteImage#H)
	// (NoteImage#T)
	// LightingN, L

	// HitResults
	// ScorePrefix (스코어 이미지 접두어)
	// ComboPrefix

	// KeyImage (버튼)
	// KeyImage#D (눌린 버튼)
	// StageLeft
	// StageRight
	// StageBottom
	// StageJudgeLine (Hint)
	// HPBar
}
