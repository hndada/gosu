aka Clavier
우선 돌아가는 게임 무언가를 만들고 싶음

// Todo: 오디오가 buffered 되었을 때, 버퍼가 끝나는 값인지 시작하는 값인지. 끝나는 값이면 그만큼 싱크 어긋나는 셈이니 곤란
type SceneMania struct {
	bufferedTime int64
	time int64
	notes []Note
	
	chart
	audio player
	endTime

	Input system
	last pressed
	staged
	hitposition	
}
func (s *SceneMania) start() {
	s.notes[i].position = n.Time
	배경 및 스테이지 그리기
	오디오, 인풋 준비
}

func (s *SceneMania) update() {
	s.updateTime() // 오디오의 현재 시간과 동일하게
	s.updateNotes()
	if 시간 >= endTime {
		select music scene으로 이동
	}
}

// todo: if, 오디오 플레이어의 시간으로부터 약 5ms 추가로 진행됐다면?
func (s *SceneMania) updateTime() {
	if s.bufferedTime == 오디오플레이어.시간 {
		s.time += 마지막 프레임으로부터 지난 시간
	} else {
		s.time = 오디오플레이어.시간
		s.bufferedTime = 오디오플레이어.시간
	} 
}

func (s *SceneMania) updateNotes() {
	s.Notes[i].Position.Y를 시간에 맞춰 변경
}
