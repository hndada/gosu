package config

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

// 값 변경과 동시에 실행되어야 하는 함수가 있는 경우 private/method로 set하는 방법으로 변경
// 그 외에도 해당 value가 직접 쓰일 수 없는 것 (mode 값 등)인 경우도 private로

// Dimness 및 Speed, 어찌됐든 플레이 중에 바뀔 수 있음
// 그러나, Settings 쪽에선 Scene을 모르므로 Scene쪽에서 값 바꾸고 함수 호출하는 식으로.

// todo: 여기다가 load 코드 옮겨오기
// 여기다가는 toml 대신 gob로.
type Settings struct {
	skin           *Skin
	screenSize     image.Point
	maxTPS         int
	GeneralDimness uint8 // todo: 0 ~ 100으로만 되게.

	Sound SoundSettings
	Mania ManiaSettings
}

// 우선 image.Point로 다뤄보고 번거로운 부분이 발견되면 대체
func (s *Settings) SetScreenSize(p image.Point) {
	s.screenSize = p
	s.Mania.Display.refresh()
	ebiten.SetWindowSize(p.X, p.Y)
}
func (s *Settings) ScreenSize() image.Point { return s.screenSize }
func (s *Settings) SetMaxTPS(tps int) {
	s.maxTPS = tps
	ebiten.SetMaxTPS(tps)
}
func (s *Settings) MaxTPS() int { return s.maxTPS }
