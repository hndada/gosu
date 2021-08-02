type SceneSelect struct {
    ChartSets [][]ChartInfo
    View []ChartInfo
    Cursor int
    마지막 music list 업데이트 시간
    
    // 실용적인 move가 되도록 하는데 필요한 field
    // MusicFilter MusicFilter
}

func (s *SceneSelect) start() {
    forced := false
    s.LoadCharts(forced)
    s.updateView()
}
func (s *SceneSelect) update() {
    한번에 차트 10개씩 표시
    좌우키 누르면 커서 10개씩 이동
    상하키 누르면 커서 1개씩 이동
    엔터 누르면 곡 재생
}

단일의 ChartSet마다 단일의 Music이 원칙 (즉, pack은 권장 X)
다른 사람이 contribute할 수 있음 (github-like system)

한 set에 chart가 너무 많을 경우 (+10 등), 
osu!의 경우 싹다 표시 -> 실용성 별로
gosu에서 "이중 스크롤"을 도입하자: set 내에 스크롤 추가 -> set별 이동 시 이동 거리 상한의 효과
(그러나 현재, set별로 뭉치는 건 실제로 쓸지 모르겠어서 일단 보류)

tree/group 대신에 *filter* 방식으로 변경
filter: ranked, 4k/7k/5~6k/8k+/~3k, 난이도 5.13~6.99, set별로 뭉치기

전체 차트 slice, 필터 적용 slice.
필터 바뀔 때마다 전체 기반으로 slice 재생성. 어차피 한번에 다룰 맵들은 많이 잡아도 1천개 정도일 테니.

type MusicFilter struct {
    ranked bool
    keys    int // 0번째 bit는 비워둠
    LV1, LV2 float64
}

// TODO: ChartHeader는 공통이라 Mania의 Keys 정보가 없음
type (mf MusicFilter) ok(h ChartHeader) bool {
    b1 := mf.keys & (1 << h.Keys) != 0
    return b1
}

func (s *SceneSelect) updateView() {
    view := make([]ChartInfo, 0, 1000)
    for i, cs := s.ChartSets {
        for j, ci := range cs {
            view = append(view, ci) // 필터 코드 추가
        } 
    }
    s.View = view
}