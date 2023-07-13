module github.com/hndada/gosu

go 1.20

require (
	github.com/faiface/beep v1.1.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hajimehoshi/ebiten/v2 v2.5.5
	github.com/ulikunitz/xz v0.5.11
	golang.org/x/image v0.9.0
	golang.org/x/sys v0.10.0
)

require (
	github.com/ebitengine/purego v0.3.0 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20221017161538-93cebf72946b // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/hajimehoshi/oto v0.7.1 // indirect
	github.com/jezek/xgb v1.1.0 // indirect
	github.com/jfreymuth/oggvorbis v1.0.5 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/exp v0.0.0-20190731235908-ec7cb31e5a56 // indirect
	golang.org/x/mobile v0.0.0-20230301163155-e0f57694e12c // indirect
	golang.org/x/sync v0.1.0 // indirect
)

retract v1.0.1 // Put the version carelessly.

retract v1.0.2 // For retracting v1.0.1.

retract v1.0.3 // For adding suffix "+incompatible" (and failed)
