package scene

import (
	"io/fs"
)

// Load is called at game.go.
func Load(fsys fs.FS) {
	data, _ := fs.ReadFile(fsys, "settings.toml")
	LoadSettings(string(data), defaultSettings)
	LoadSkin(fsys, defaultSkin)
}
