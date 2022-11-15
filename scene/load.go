package scene

import (
	"embed"
	"fmt"
	"io/fs"
)

// For loading default settings and skin.
func init() {
	// Settings{}.init()
	// Skin{}.init()

	//go:embed skin/*
	var fs embed.FS
	defaultSkin.Load(fs, nil)
}

// Todo: separate to LoadSettings, LoadSkin?
func Load(fsys fs.FS) {
	func() {
		data, err := fs.ReadFile(fsys, "settings.toml")
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		settings, err := Settings{}.Load(string(data))
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		Settings{}.Set(settings)
	}()
	func() {
		skin, err := Skin{}.Load(fsys)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		Skin{}.Set(skin)
	}()
}
