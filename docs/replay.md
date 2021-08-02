// time series; acending order
// 새거 추가된 뒤 매번 Sort()
type KeyLogs struct {
	logs []keyLog
}

func (l KeyLogs) Search(time int64) int {
	idx := sort.Search(len(l.logs), func(i int) bool { 
			return l.logs[i].time >= time 
		})
	if idx < len(l.logs) && l.logs[idx].time == time { // 이미 해당 시간에 기록 남아있음
		return idx
	}
	return -1
} 

func (l KeyLogs) Sort() {
	sort.Slice(l.logs, func(i, j int) bool { 
		return l.logs[i].time < l.logs[j].time 
	})
}
