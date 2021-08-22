Sprite
1. 최소 노출 시간 구현: 100ms
2. 1개 원본 두고 매번 복사해서 쓰기
3. 애니메이션: 시간초, MaxTPS 를 param으로 하여 각 프레임마다 띄울 이미지 index 정하기 
4. SetFixed()
5. 연출: 스코어 스크롤, 콤보 튕기기 (프레임마다 op다르게 주면 될듯)
판정, 체력, 노트/롱노트 히트 이펙트, 판정

input 
https://github.com/eiannone/keyboard
goroutine

주석, docs, taiko 정리
game->engine, game(common) 으로 구분
VNC server 구축

* TimingPoint는 각 모드 별로 알아서 불러오는 걸로
* NewBaseChartFromOsu() 에서 path 파라미터 제거, path 대입 삭제
    - ParseHeaderFromOsu() 로 변경
* Level은 Header에 있으면 복잡할듯. 
    - mods마다 바뀌는 level.
* Header에 path 필요 없음