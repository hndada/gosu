package scene

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/resources"
	"github.com/hndada/gosu/ui"
)

// Resources are loaded from file system.
// Options are set by user and saved to file system.
type Scene interface {
	// New(*Game, Args) (Scene, error)
	Update() any
	Draw(screen draws.Image)
	WindowTitle() string
	DebugString() string
}

type Context struct {
	ui.KeyboardState

	Resources Resources
	// Option should not be passed by value,
	// as its value should be shared and modified.
	Options  *Options
	Handlers Handlers
	Database Database
}

func NewContext(fsys fs.FS) (*Context, error) {
	c := &Context{
		KeyboardState: ui.KeyboardState{},
	}

	if resFS, err := fs.Sub(fsys, "resources"); err == nil {
		c.Resources = NewResources(resFS)
	} else {
		c.Resources = NewResources(resources.DefaultFS)
	}

	// NewOptions is always called, as there
	// might be omitted fields on a local option file.
	c.Options = NewOptions()
	if data, err := fs.ReadFile(fsys, "options.json"); err == nil {
		if err := json.Unmarshal(data, c.Options); err != nil {
			fmt.Printf("Failed to unmarshal options.json: %v\n", err)
		}
	}
	// It is always necessary to set derived values.
	c.Options.Normalize()
	c.Options.Piano.SetDerived()

	c.Handlers = NewHandlers(c.Options, &c.KeyboardState)

	dbs, err := NewDatabase(fsys)
	if err != nil {
		return nil, err
	}
	c.Database = dbs

	return c, nil
}
