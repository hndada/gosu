module github.com/hndada/gosu

go 1.18

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/hajimehoshi/ebiten/v2 v2.4.12
	github.com/ulikunitz/xz v0.5.10
	golang.org/x/image v0.1.0
	golang.org/x/sys v0.2.0
)

require (
	github.com/ebitengine/purego v0.2.0-alpha.0.20221107040104-382f4c69b854 // indirect
	github.com/hajimehoshi/file2byteslice v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20221017161538-93cebf72946b // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/hajimehoshi/oto/v2 v2.4.0-alpha.6 // indirect
	github.com/jezek/xgb v1.1.0 // indirect
	github.com/jfreymuth/oggvorbis v1.0.4 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/vmihailenco/msgpack/v5 v5.3.5
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/exp/shiny v0.0.0-20221106115401-f9659909a136 // indirect
	golang.org/x/mobile v0.0.0-20221110043201-43a038452099 // indirect
)

retract [v1.0.1, v1.0.3]
retract v1.0.1 // Put the version carelessly.
retract v1.0.2 // For retracting v1.0.1.
retract v1.0.3 // For adding suffix "+incompatible" (and failed)