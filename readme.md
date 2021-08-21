체력, 스테이지 커버 (지나간 노트 가리기)
애니메이션: 노트/롱노트 히트 이펙트 
- 애니메이션 시간초, MaxTPS 를 param으로 하여 각 프레임마다 띄울 이미지 index 정하기 
위치 조정: 스코어, 콤보 
연출: 스코어, 콤보 튕기기 (프레임마다 op다르게 주면 될듯)

키보드 hook
주석, docs 정리
VNC server 구축


* Song -> Music
* Title -> MusicName
* TimingPoint는 각 모드 별로 알아서 불러오는 걸로
* NewBaseChartFromOsu() 에서 path 파라미터 제거, path 대입 삭제
    - ParseHeaderFromOsu() 로 변경
* Level은 Header에 있으면 복잡할듯. 
    - mods마다 바뀌는 level.
* Header에 path 필요 없음