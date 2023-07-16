package mode

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
)

// := is called short assignment statement. When assigning multiple variables
// by :=, at least one of the variables on the left side must be newly declared.
// It will work just as = for already existing variables.
// https://go.dev/play/p/5SUt9uyrncD
func NewReplay(fsys fs.FS, name string, keyCount int) (input.KeyboardReader, error) {
	file, err := fsys.Open(name)
	if err != nil {
		return input.KeyboardReader{}, fmt.Errorf("failed to open replay file: %s", err)
	}

	switch strings.ToLower(filepath.Ext(name)) {
	case ".osr":
		f, err := osr.NewFormat(file)
		if err != nil {
			return input.KeyboardReader{}, fmt.Errorf("failed to parse replay file: %s", err)
		}
		return f.KeyboardReader(keyCount), nil
	default:
		return input.KeyboardReader{}, fmt.Errorf("unsupported replay file format")
	}
}
