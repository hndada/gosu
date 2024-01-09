module github.com/hndada/gosu

go 1.21

toolchain go1.21.5

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/gopxl/beep v1.3.0
	github.com/hajimehoshi/ebiten/v2 v2.7.0-alpha.4
	github.com/ulikunitz/xz v0.5.11
	golang.org/x/exp v0.0.0-20231226003508-02704c960a9b
	golang.org/x/image v0.14.0
	golang.org/x/sys v0.14.0
)

require (
	github.com/ebitengine/oto/v3 v3.2.0-alpha.2.0.20231021101548-b794c0292b2b // indirect
	github.com/ebitengine/purego v0.6.0-alpha.1 // indirect
	github.com/go-text/typesetting v0.0.0-20231110223828-31a9559ebc00 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/jezek/xgb v1.1.0 // indirect
	github.com/jfreymuth/oggvorbis v1.0.5 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/exp/shiny v0.0.0-20230817173708-d852ddb80c63 // indirect
	golang.org/x/mobile v0.0.0-20231108233038-35478a0c49da // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

retract v1.0.1 // Put the version carelessly.

retract v1.0.2 // For retracting v1.0.1.

retract v1.0.3 // For adding suffix "+incompatible" (and failed)
