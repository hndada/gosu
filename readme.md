싱크 문제: 처음에 오디오 미리 load
판정: 새 판정 나올 때마다 새로 애니메이션 재생 
8키 1+7
select: box, bg preview

input: https://github.com/eiannone/keyboard 
목표: 누른 거 떼는 것도 파악
시간, 키, down/up

주석, docs, taiko 정리
ebiten v2로 업그레이드
NewHeaderFromOsu()에서 Path 삭제 시도
판정오차 막대기
VNC server 구축

game->engine, game(common) 으로 구분. sedyh's ebiten ecs model 도입 이후 한번 시도
https://github.com/andygeiss/ecs
Sprite 연출: 스코어 스크롤, 콤보 튕기기, Fade-in.
    - 프레임마다 op다르게 주면 될듯
    - 전용 패키지 만들어보기