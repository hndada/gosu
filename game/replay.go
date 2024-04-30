package game

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

//	type Replay interface {
//		KeyboardStates() []input.KeyboardState
//	}
type Replay = *input.KeyboardStateBuffer

// := is called short assignment statement. When assigning multiple variables
// by :=, at least one of the variables on the left side must be newly declared.
// It will work just as = for already existing variables.
// https://go.dev/play/p/5SUt9uyrncD

// If directory path is passed in fs.ReadFile, it will return an error.
func NewReplay(fsys fs.FS, name string, keyCount int) (Replay, error) {
	dat, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(filepath.Ext(name)) {
	case ".osr":
		f, err := osr.NewFormat(dat)
		if err != nil {
			err = fmt.Errorf("failed to parse replay file: %s", err)
			return nil, err
		}
		states := f.KeyboardStates(keyCount)
		r := input.NewKeyboardStateBuffer(states)
		r.Trim()
		return r, nil
	}
	return nil, fmt.Errorf("unsupported replay file format")
}
