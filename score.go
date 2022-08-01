package main


if done {
    if 롱노트 꼬리 {
        (시간-miss 범위) 지났으면 flush/fetch (ff)
    } else {
        panic // 다른 노트는 done이면 이미 ff 되었어야 함
    }
} else {
    result := miss
    // 시간차 = note.Time - 현재 시간
    // 시간차 < - miss 범위 이면 더 처리할 것 없음
    if 시간차 >= - miss 범위 && 시간차 < miss 범위 && judgable(type, action) {
        if 시간차 < 0 {
            시간차 *= -1
        }
        for j, w := range ws {
            if 시간차 < w.Window {
                result = j
                break
            }
        }
    }
    mark(notes, i, j)
    // 키음 재생
    staged[key] = notes[i].next
}

func mark(notes, i, j) {
    score() // karma, w를 저장하고 있어야 함 // Scene에서 관리해야할듯 
    notes[i].Done = true
    if notes[i].Type == 머리 && j == Miss {
        mark(notes, notes[i].next, Miss)
    }
}

func judgable(t NoteType, a KeyAction) bool {
	if t == Tail {
		return a == Release
	}
	return a == Press
}

func inRange(time int64, j Judgment) {
    return time < j.Window && time > -j.Window 
}