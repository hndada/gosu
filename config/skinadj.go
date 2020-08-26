package config

import "github.com/hajimehoshi/ebiten"

// screen에 처음부터 fixed 된 상태로 그려질 애들
// option 건들 때마다 갱신

type ScaledManiaSkin struct {
	Note             [4]ebiten.Image
	LNHead           [4]ebiten.Image
	LNBody           [4][]ebiten.Image
	LNTail           [4]ebiten.Image
	KeyButton        [4]ebiten.Image
	KeyButtonPressed [4]ebiten.Image

	HitResults [5]ebiten.Image
	// NoteLighting
	// LNLighting
	// StageLeft   ebiten.Image
	StageRight  ebiten.Image
	StageBottom ebiten.Image
	StageHint   ebiten.Image
}

type ManiaStage struct {
	Image *ebiten.Image
	Op    *ebiten.DrawImageOptions
}

// 먼저 100 스케일로 그리고 확대하면 깨지니까
// todo: 판정선 가운데에 노트 가운데가 맞을 때 Max가 뜨게
func NewStageMania() ManiaStage {
	var sm ManiaStage
	sm.Op = &ebiten.DrawImageOptions{}
	// skin에서 불러오기
	var (
		main        *ebiten.Image // fieldWidth, screenHeight (generated)
		stageBottom *ebiten.Image // fieldWidth, 폭맞춤 y
		stageHint   *ebiten.Image // fieldWidth, 설정값 ('노트와 동일한 높이로' 옵션 추가)
		stageRight  *ebiten.Image // 폭맞춤x, screenHeigth
		hpBarFrame  *ebiten.Image // 폭맞춤x, screenHeigth
	)
	stageCenter := float64(s.g.ScreenSize().X) * s.g.StagePosition / 100
	var fieldWidth float64
	sm.Op.GeoM.Translate(stageCenter, 0)
	sm.Op.GeoM.Translate(-fieldWidth/2, 0)
}

func (s *ManiaStage) render() {
	// 필드의 중앙이 스크린의 중앙에 오게 op.GeoM.Translate(dx, 0)
}
