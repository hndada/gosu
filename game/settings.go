package game

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// 값 변경과 동시에 실행되어야 하는 함수가 있는 경우 private/method로 set하는 방법으로 변경
// 그 외에도 해당 value가 직접 쓰일 수 없는 것 (mode 값 등)인 경우도 private로

// Dimness 및 Speed, 어찌됐든 플레이 중에 바뀔 수 있음
// 그러나, Settings 쪽에선 Scene을 모르므로 Scene쪽에서 값 바꾸고 함수 호출하는 식으로.
type SettingsTemplate struct { // temporary struct for load/save settings
	screenSize     image.Point
	maxTPS         int
	generalDimness uint8 // todo: 0 ~ 100으로만 되게.
	volumeMaster   uint8
	// volumeBGM      uint8
	// volumeSFX      uint8
}

var Settings SettingsTemplate

// 우선 image.Point로 다뤄보고 번거로운 부분이 발견되면 대체
// 아래 애들은 별도로 설정
func SetScreenSize(p image.Point) {
	Settings.screenSize = p
	ebiten.SetWindowSize(p.X, p.Y)
}
func ScreenSize() image.Point { return Settings.screenSize }
func ScaleY() float64         { return float64(Settings.screenSize.Y) / 100 }
func SetMaxTPS(tps int) {
	Settings.maxTPS = tps
	ebiten.SetMaxTPS(tps)
}
func MaxTPS() int { return Settings.maxTPS }

func GeneralDimness() uint8 { return Settings.generalDimness }

// todo: Streamer 2개 만들기: BGM, SFX
// todo: Streamer에 vol 꽂기
func SetVolumeMasterDown() {
	switch {
	case Settings.volumeMaster <= 0:
		return
	case Settings.volumeMaster <= 5:
		Settings.volumeMaster -= 1
	case Settings.volumeMaster <= 100:
		Settings.volumeMaster -= 5
	}
	// percent(Settings.volumeBGM) * percent(Settings.volumeMaster)
}

// func percent(v uint8) float64 { return float64(v) / 100 }

// const settingsFilename = "common.settings"

func LoadSettings() {
	// if f, err := ioutil.ReadFile(settingsFilename); err != nil {
	// 	ResetSettings()
	// } else {
	// 	r := bytes.NewReader(f)
	// 	dec := gob.NewDecoder(r)
	// 	if err := dec.Decode(&Settings); err != nil {
	// 		panic(err)
	// 	}
	// }
	ResetSettings()
	ebiten.SetMaxTPS(MaxTPS())
	ebiten.SetWindowSize(ScreenSize().X, ScreenSize().Y)
}

// func SaveSettings() {
// 	var b bytes.Buffer
// 	enc := gob.NewEncoder(&b)
// 	err := enc.Encode(&Settings)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	f, err := os.Create(settingsFilename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	_, err = f.Write(b.Bytes())
// 	if err != nil {
// 		panic(err) // todo: 파일 날아가는거 대비?
// 	}
// }

func ResetSettings() {
	Settings.maxTPS = 60
	Settings.screenSize = image.Pt(800, 600)
	Settings.generalDimness = 30
	Settings.volumeMaster = 30
	// Settings.volumeBGM = 50
	// Settings.volumeSFX = 50
}
