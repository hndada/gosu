package config

type SoundSettings struct {
	volumeMaster uint8
	volumeBGM    uint8
	volumeSFX    uint8
}

// todo: Streamer 2개 만들기: BGM, SFX
func (s *SoundSettings) SetDownVolumeMaster() {
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