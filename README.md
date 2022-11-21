# gosu

Classic rhythm games written in go with Ebitengine

# How to play
Change the mode with `F1`

Change the Speed with `PageUp / PageDown`

Select the song with `Enter`

Press matching keys with notes!

You can change key settings by modifying `keys.txt`. Default Key settings are below:
```
4 Key: S, D, J, K
7 Key: S, D, F, Space, J, K, L
Drum:  S, D, J, K
```

# Game play preview
Click thumbnails to watch at YouTube.

[![Taishi - bluefieldcreator [Etherealization]](https://i.imgur.com/DN8JTzQ.png)](https://youtu.be/9kMUT8vQI24)

[![The Flashbulb - The Bridgeport Run [Escapism]](https://i.imgur.com/tIVTiXo.png)](https://youtu.be/5VWaSAs7bbQ)

[![cillia - Ringo Uri no Utakata Shoujo [Ringo Oni]](https://i.imgur.com/0Ven6Oa.png)](https://youtu.be/8VgzAlc4SJ0)


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
  * Others is currently depending on `ebiten.IsKeyPressed` .

* Codebase with high readability
  * Super-fast in loading files and playing.
  * Fairly scalable for future work.

# Build

1. For MacOS and Linux users, install Ebitengine dependencies first by referring to the documentation([Ebitengine/Install](https://ebitengine.org/en/documents/install.html)).

2. Go to root directory of the repository and build as below. 

```zsh
cd cmd/gosu
go build .
```

3. Run gosu

```zsh
./gosu
```

# Flow of game logic
[Powerpoint and SlideShare.](https://www.slideshare.net/MuangMuangE/gosupresentpptx-253675145)

Will also post details at [wiki](https://github.com/hndada/gosu/wiki).

# Community
~~[Discord server](https://discord.gg/4TztgpaC)~~ Will be open after stable version has been release.

# License
Codebase: Apache License 2.0

Skin images and music tracks are from [osu-resources](https://github.com/ppy/osu-resources), licensed under [CC-BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/legalcode).
