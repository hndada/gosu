package scene

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/resource"
	"github.com/hndada/gosu/ui"
	"github.com/hndada/gosu/util"
)

// Resources are loaded from file system.
// Options are set by user and saved to file system.
type Scene interface {
	Update() any
	Draw(screen draws.Image)

	WindowTitle() string
	DebugString() string
}

// Shared values should be passed by reference in general.
// However, fields of Resources, Handlers, Database are
// either pointer or read-only.
type Context struct {
	*ui.KeyboardState
	Resources Resources
	// Option should not be passed by value,
	// as its value should be shared and modified at runtime.
	Options  *Options
	Handlers Handlers
	Database Database
}

func NewContext(fsys fs.FS) (Context, error) {
	c := Context{
		KeyboardState: &ui.KeyboardState{},
	}

	c.Resources = NewResources(resource.DefaultFS)
	if ok, _ := util.DirectoryExists(fsys, "resource"); ok {
		resFS, err := fs.Sub(fsys, "resource")
		if err == nil {
			c.Resources = NewResources(resFS)
		}
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

	c.Handlers = NewHandlers(c.Options, c.KeyboardState)

	dbs, err := NewDatabase(fsys)
	if err != nil {
		return Context{}, err
	}
	c.Database = dbs

	return c, nil
}
