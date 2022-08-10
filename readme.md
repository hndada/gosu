# gosu

A classic rhythm game written in go

# How to play
Select the song with `Enter`, and press keys corresponding to falling notes.
Following is default key settings:
```
var KeySettings = map[int][]Key{
	4:               {KeyD, KeyF, KeyJ, KeyK},
	5:               {KeyD, KeyF, KeySpace, KeyJ, KeyK},
	6:               {KeyS, KeyD, KeyF, KeyJ, KeyK, KeyL},
	7:               {KeyS, KeyD, KeyF, KeySpace, KeyJ, KeyK, KeyL},
	8 + LeftScratch: {KeyA, KeyS, KeyD, KeyF, KeySpace, KeyJ, KeyK, KeyL},
	8:               {KeyA, KeyS, KeyD, KeyF, KeyJ, KeyK, KeyL, KeySemicolon},
	9:               {KeyA, KeyS, KeyD, KeyF, KeySpace, KeyJ, KeyK, KeyL, KeySemicolon},
	10:              {KeyA, KeyS, KeyD, KeyF, KeyV, KeyN, KeyJ, KeyK, KeyL, KeySemicolon},
}
```

# Game play preview
## Video
[Here is the YouTube link of live-playing! (replay)](https://youtu.be/YMRgGQZHpQo)

## Screenshots
![4 Key](https://i.imgur.com/6veaLI6.png)

![7 Key](https://i.imgur.com/MJTFmE3.png)
# Feature
* osu! files supported
  * .osu (osu! beatmap file)
    * Speed-change effects work (called `SV`).
  * .osr (osu! replay file)
    * Put replay files at `replay/` with `ReplayMode` at select scene.

* Skinnable in-game images
  * Put your favorite skin in `skin/` (should match the file name though).
  * Image size in game are settable by user (WIP).
    * You can try it right now with changing value at `settings.go` and build. 

* Effective score and level system (originally designed)
  * The motivation of gosu dev
  * Tried to make feel score and actual performance are related.
  * Level calculation is currently primitive 
    * Will be exquisite in a short time 

* Quick input supported (1ms)
  * *Hook* is used in `Windows`.
  * Others is currently depending on `ebiten.IsKeyPressed` .

* Lightweight codebase
  * Rebase version of previous *gosu* (see `v0` branch)
    * `v0` was verbose and has complex structure. 
  * Super-fast in loading files and playing
  * Fairly scalable for future work

# Flow of game logic
Work in Progress

Will also post details at [wiki](https://github.com/hndada/gosu/wiki).

# Build
`cd cmd/gosu`

`go build .`

# License
Skin images and music tracks are from [osu-resources](https://github.com/ppy/osu-resources), licensed under [CC-BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/legalcode).