input: 시간, 키, down/up
    https://github.com/eiannone/keyboard 
    목표: 누른 거 떼는 것도 파악

8키 1+7
판정오차 막대기: 문제 없이 그려지고 있는것 같은데 어째선지 랙 발생
select scene에서 bg preview: 어째선지 랙 발생
chart-bg.go 정리

주석, docs, taiko 정리
ebiten v2로 업그레이드
NewHeaderFromOsu()에서 Path 삭제 시도
VNC server 구축

game->engine, game(common) 으로 구분. sedyh's ebiten ecs model 도입 이후 한번 시도
https://github.com/andygeiss/ecs
Sprite 연출: 스코어 스크롤, 콤보 튕기기, Fade-in.
    - 프레임마다 op다르게 주면 될듯
    - 전용 패키지 만들어보기