체력바 및 내부
Sprite 애니메이션: 시간초, MaxTPS 를 param으로 하여 각 프레임마다 띄울 이미지 index 정하기 
판정,  노트/롱노트 히트 이펙트

input: https://github.com/eiannone/keyboard 
목표: 누른 거 떼는 것도 파악
시간, 키, down/up

8키 1+7
싱크 문제: 처음에 오디오 미리 load
select: box, bg preview

주석, docs, taiko 정리
game->engine, game(common) 으로 구분
ebiten v2로 업그레이드
NewHeaderFromOsu()에서 Path 삭제 시도
VNC server 구축
Sprite 연출: 스코어 스크롤, 콤보 튕기기, Fade-in.
    - 프레임마다 op다르게 주면 될듯
    - 전용 패키지 만들어보기