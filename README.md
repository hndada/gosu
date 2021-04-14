# gosu
New client of mania players, by a mania player, for mania players
=================================================================
우선 변수를 줄이기 위해 stamina 요소는 비활성화
// strain.go
    // 우선 chord 알고리즘 먼저 고쳐보자
// difficulty.go
    // ppy 방식처럼, 구간 내 최고 strain을 잡아야 할까?
// internal/tools
    // score, level 다 정리되고 나서 정리하겠음
// hand.go: 
    func lnLocation() int {}: hold outer, inner, adj 한번에 작성하기
0. rg-parser
    - [ ] Parse sound samples 
1. gosu-calc
2. gosu-nity
    - [ ] Essential features
        - [ ] chart synced with music
        - [ ] Quick input system
        - [ ] Level, score, hp, pp system
    - [ ] Local/online Leaderboard
        * [ ] Real-time local/online score competition

3. Mwang (gosu server)
    - [ ] User struct
    - [ ] Collect user info, scores and replay
        - [ ] Collceting user info: Key stroke, play time count
    - [ ] Chatting and multiplaying

4. Website
    - [ ] Ranking page
    - [ ] Beatmap download pages (for collecting all separated files)
    - [ ] Beatmap discussion
    - [ ] Beatmap git system
    - [ ] Userpage