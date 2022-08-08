# gosu

Images, music are from [https://github.com/ppy/osu-resources]

## Todo
* Key input
    1. ebiten.IsKeyPressed: 얘는 Update에 물려있으니, Scene의 Tick 에 의존해야 함
    2. hook
    3. replay - real time
        3-1. replay - result only: 얘는 ebiten.Game에 의존하지 않고 for loop 돌면서 Update 호출하면 됨.
* skin.go
    ScenePlay의 Sprite들 global로 빼기
    Sprite Draw할 때 Scale 자동 적용?
* Songs
    * osu resources로 짧게 형식만
* 문서 작업
* level
    * Unbalanced chart will be penalized with strict decay factor 

## Issue
* Precise input time 
    * Written code, inspired by "go-hook"
    * [https://github.com/petercunha/GoLANG-Windows-KeyLogger/blob/master/w32Keylogger.go]
    * [https://gist.github.com/obonyojimmy/61abcbc6022cb7399813db8ac4d1de4d]
    * [https://golangscript.com/g/simple-proof-of-concept-windows-go-keylogger-via-conventional-windows-apis-setwindowshookex-low-level-keyboard-hook]
* Syncing notes with music position