# gosu

Classic rhythm games written in go with ebitengine

# How to play
Change the mode with `Ctrl`

Change the Speed with `Z/X`

Select the song with `Enter`

Press matching keys with notes!

Following is default key settings:
```
4 Key: S, D, J, K
7 Key: S, D, F, Space, J, K, L
Drum:  S, D, J, K
```

# Game play preview
## Video
[cillia - Ringo Uri no Utakata Shoujo [Ringo Oni]](https://youtu.be/8VgzAlc4SJ0)

[The Flashbulb - The Bridgeport Run [Escapism]](https://youtu.be/5VWaSAs7bbQ)

[Real-time replay of Taishi - bluefieldcreator [Etherealization]](https://www.youtube.com/watch?v=9kMUT8vQI24&list=PLQhd8A8gGbIBm_oJdW5K9Pwv9jZpmJzLW&index=2&ab_channel=MuangMuangE)

## Screenshots
![4 Key](https://i.imgur.com/6veaLI6.png)

![7 Key](https://i.imgur.com/MJTFmE3.png)

![Drum](https://i.imgur.com/VquWLWk.png)

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
  * The motivation of gosu dev.
  * Tried to make feel score and actual performance are related.
  * Level calculation is currently naive. 
    * Will be exquisite in a short time.

* Quick input supported (1ms)
  * *Hook* is used in `Windows`.
  * Others is currently depending on `ebiten.IsKeyPressed`.

* Codebase with high readability
  * Super-fast in loading files and playing
  * Fairly scalable for future work

# Build
Go to root directory of the repository first. 
```
cd cmd/gosu
go build .
```

# Flow of game logic
Work in Progress

Will also post details at [wiki](https://github.com/hndada/gosu/wiki).

# License
Codebase: MIT

Skin images and music tracks are from [osu-resources](https://github.com/ppy/osu-resources), licensed under [CC-BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/legalcode).
