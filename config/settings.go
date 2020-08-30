package config

import (
	"bytes"
	"encoding/gob"
	"github.com/hajimehoshi/ebiten"
	"image"
	"io/ioutil"
	"os"
)

// 값 변경과 동시에 실행되어야 하는 함수가 있는 경우 private/method로 set하는 방법으로 변경
// 그 외에도 해당 value가 직접 쓰일 수 없는 것 (mode 값 등)인 경우도 private로

// Dimness 및 Speed, 어찌됐든 플레이 중에 바뀔 수 있음
// 그러나, Settings 쪽에선 Scene을 모르므로 Scene쪽에서 값 바꾸고 함수 호출하는 식으로.
type Settings struct {
	screenSize     image.Point
	maxTPS         int
	generalDimness uint8 // todo: 0 ~ 100으로만 되게.
	volumeMaster   uint8
	volumeBGM      uint8
	volumeSFX      uint8
	ManiaSettings
}

// 우선 image.Point로 다뤄보고 번거로운 부분이 발견되면 대체
// 아래 애들은 별도로 설정
func (s *Settings) SetScreenSize(p image.Point) {
	s.screenSize = p
	ebiten.SetWindowSize(p.X, p.Y)
}
func (s *Settings) ScreenSize() image.Point { return s.screenSize }
func (s *Settings) SetMaxTPS(tps int) {
	s.maxTPS = tps
	ebiten.SetMaxTPS(tps)
}
func (s *Settings) MaxTPS() int { return s.maxTPS }

// todo: Streamer 2개 만들기: BGM, SFX
func (s *Settings) SetDownVolumeMaster() {
	switch {
	case s.volumeMaster <= 0:
		return
	case s.volumeMaster <= 5:
		s.volumeMaster -= 1
	case s.volumeMaster <= 100:
		s.volumeMaster -= 5
	}
	// todo: Streamer에 vol 꽂기
	// percent(s.volumeBGM) * percent(s.volumeMaster)
}

// func percent(v uint8) float64 { return float64(v) / 100 }

func (s Settings) maniaStageCenter() int {
	return s.ManiaSettings.stageCenter(s.screenSize)
}

func (s *Settings) Load() {
	f, err := ioutil.ReadFile("settings")
	if err != nil {
		s.Reset()
		return
	}
	r := bytes.NewReader(f)
	dec := gob.NewDecoder(r)
	err = dec.Decode(s)
	if err != nil {
		panic(err)
	}
}

func (s *Settings) Save() {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(s)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("settings")
	if err != nil {
		panic(err)
	}
	_, err = f.Write(b.Bytes())
	if err != nil {
		panic(err) // todo: 파일 날아가는거 대비?
	}
}

func (s *Settings) Reset() {
	s.maxTPS = 240
	s.screenSize = image.Pt(1600, 900)
	s.generalDimness = 30
	s.volumeMaster = 100
	s.volumeBGM = 50
	s.volumeSFX = 50
	s.ManiaSettings.reset()
}
