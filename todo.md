* Song -> Music
* Title -> MusicName
* BaseChart -> ChartHeader 
    - TimingPoint는 각 모드 별로 알아서 불러오는 걸로
* NewBaseChartFromOsu() 에서 path 파라미터 제거, path 대입 삭제
    - ParseHeaderFromOsu() 로 변경
* Level은 Header에 있으면 복잡할듯. 
    - mods마다 바뀌는 level.
* Header에 path 필요 없음