# gosu

Classic rhythm games written in go with Ebitengine

Latest version: 0.6 (July 23rd, 2023)

# Game play preview
Click thumbnails to watch at YouTube.

[![Taishi - bluefieldcreator [Etherealization]](https://i.imgur.com/DN8JTzQ.png)](https://youtu.be/9kMUT8vQI24)

[![The Flashbulb - The Bridgeport Run [Escapism]](https://i.imgur.com/tIVTiXo.png)](https://youtu.be/5VWaSAs7bbQ)

[![cillia - Ringo Uri no Utakata Shoujo [Ringo Oni]](https://i.imgur.com/0Ven6Oa.png)](https://youtu.be/8VgzAlc4SJ0)


# How to play
1. Select the song with `Enter`.
2. Press matching keys with notes.
3. Change the Speed with `PageUp / PageDown`

```
4 Key: S, D, J, K
7 Key: S, D, F, Space, J, K, L
```

# Feature
* osu! files supported
  * .osu (osu! beatmap file)
  * .osr (osu! replay file)

* Practical score and level system
  * The motivation of gosu dev.
  * WIP: Level calculation

* Customize in-game sprites
  * Put your favorite skin in `asset/` with matching name.

* Quick input listener
  * WINAPI is used in `Windows`.
  * Others is currently depending on `ebiten.IsKeyPressed` .

* Codebase with high readability

# Build
1. For MacOS and Linux users, install Ebitengine dependencies first by referring to the 
documentation([Ebitengine/Install](https://ebitengine.org/en/documents/install.html)).

2. Go to root directory of the repository and type: 
```zsh
go build .
```


# Web version
Version: 0.4.1

**[https://gosu-web-orcin.vercel.app](https://gosu-web-orcin.vercel.app)**

# Game structure
### Package flow
![Game structure](https://i.imgur.com/gwFA6es.png)

### [Introduction of gosu development](https://www.slideshare.net/MuangMuangE/gosupresentpptx-253675145)
[![gosu-present](https://i.imgur.com/rtq5n9p.png)](https://www.slideshare.net/MuangMuangE/gosupresentpptx-253675145)

Will also post details at [wiki](https://github.com/hndada/gosu/wiki).

# License
Codebase: Apache License 2.0

Most skin images and music tracks are from [osu-resources](https://github.com/ppy/osu-resources), 
licensed under [CC-BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/legalcode).
