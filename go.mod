module github.com/hndada/gosu

go 1.21

toolchain go1.21.5

require (
	github.com/gopxl/beep v1.4.1
	github.com/hajimehoshi/ebiten/v2 v2.7.6
	github.com/ulikunitz/xz v0.5.12
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8
	golang.org/x/image v0.18.0
	golang.org/x/sys v0.21.0
)

require (
	github.com/coder/websocket v1.8.12 // indirect
	github.com/ebitengine/gomobile v0.0.0-20240518074828-e86332849895 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/oto/v3 v3.2.0 // indirect
	github.com/ebitengine/purego v0.7.1 // indirect
	github.com/go-text/typesetting v0.1.1 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/jfreymuth/oggvorbis v1.0.5 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/exp/shiny v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/mobile v0.0.0-20240604190613-2782386b8afd // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

retract v1.0.1 // Put the version carelessly.

retract v1.0.2 // For retracting v1.0.1.

retract v1.0.3 // For adding suffix "+incompatible" (and failed)
