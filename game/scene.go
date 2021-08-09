package game

type PlayScene struct {
	Tick        int64
	AudioPlayer *AudioPlayer
}

type TransSceneArgs struct {
	// 자기 자신은 .(type)으로 알 수 있음
	Next string // next scene name
	Args interface{}
}

// always follows audio's time
// func (s *Scene) Time() int64 {
// 	return s.audioPlayer.Time().Milliseconds()
// }

// 이 방법을 하려면 tps가 게임 중에 변하지 않아야 함
// CurrentTPS가 약간 딱 떨어지지 않는 게 마음에 걸리지만, 곧 보충되어 결과적으로 일정히 유지 된다고 상정하겠음
// -> Audio에서 Time 따오는 게 제일 정확. 그런데 지금 오디오가 내주는 시간이 버퍼에 의해 정확하지 않음
func (ps *PlayScene) Time() int64 {
	return Millisecond * ps.Tick / int64(MaxTPS()) // todo: MaxFPS로 해야하나? 240 설정했는데 모니터가 60일 수 있음. 한편 FPS는 float?
}
