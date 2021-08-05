스코어 고치기



// strain.go
    변인을 통제해야함: stamina, Hold
    우선 chord 알고리즘 먼저 고쳐보자
// difficulty.go
    // ppy 방식처럼, 구간 내 최고 strain을 잡아야 할까?
// internal/tools
    // score, level 다 정리되고 나서 정리하겠음
// hand.go: 
    func lnLocation() int {}: hold outer, inner, adj 한번에 작성하기


* Song -> Music
* Title -> MusicName
* BaseChart -> ChartHeader 
    - TimingPoint는 각 모드 별로 알아서 불러오는 걸로
* NewBaseChartFromOsu() 에서 path 파라미터 제거, path 대입 삭제
    - ParseHeaderFromOsu() 로 변경
* Level은 Header에 있으면 복잡할듯. 
    - mods마다 바뀌는 level.
* Header에 path 필요 없음

기획 -> 로직 -> 코드 

곡 선택 Input Control

// todo: timing points, (decending/ascending) order로 sort -> rg-parser에서