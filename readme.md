Sprite
1. 최소 노출 시간 구현: 100ms
2. 1개 원본 두고 매번 복사해서 쓰기
3. 애니메이션: 시간초, MaxTPS 를 param으로 하여 각 프레임마다 띄울 이미지 index 정하기 
4. SetFixed()
5. 연출: 스코어 스크롤, 콤보 튕기기 (프레임마다 op다르게 주면 될듯)
판정, 체력, 노트/롱노트 히트 이펙트, 판정

input 
1. https://github.com/eiannone/keyboard
2. goroutine

주석, docs, taiko 정리
game->engine, game(common) 으로 구분
NewHeaderFromOsu()에서 Path 삭제 시도
VNC server 구축