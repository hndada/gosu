UI 꾸미기
스코어/콤보 위치
스코어/콤보 다이나믹 (프레임마다 op다르게 주면 될듯)
판정, 체력, 오차, 키버튼 누름/떼기, HitPosition 밑으로는 안 보이게.
히트 이펙트 롱노트 애니메이션

키보드 hook
주석, docs 정리

* Song -> Music
* Title -> MusicName
* BaseChart -> ChartHeader 
    - TimingPoint는 각 모드 별로 알아서 불러오는 걸로
* NewBaseChartFromOsu() 에서 path 파라미터 제거, path 대입 삭제
    - ParseHeaderFromOsu() 로 변경
* Level은 Header에 있으면 복잡할듯. 
    - mods마다 바뀌는 level.
* Header에 path 필요 없음