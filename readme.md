# gosu

Images, music are from [https://github.com/ppy/osu-resources]

## Todo
* input + score를 scene (-> ebiten)의 Update()에 의존하지 않고 자체적으로 Update() 돌기
* replay를 real-time/score-calc 두 가지로 나누기
* skin.go
    ScenePlay의 Sprite들 global로 빼기
    Sprite Draw할 때 Scale 자동 적용?

* music, replay, skin 폴더를 gosu/gosu/에 옮기기
* 문서 작업

* level
    * Unbalanced chart will be penalized with strict decay factor 