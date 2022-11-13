package format

import (
	"io"
	"path/filepath"
	"strings"
)

// File: a collection of data forming a single unit.
// There is no word Filetype. File Type is.
type Type int
type File struct {
	io.Reader
	io.Seeker
	io.Closer
	io.Writer
	Type
	// Data []byte
}

func Ext(path string) string { return strings.ToLower(filepath.Ext(path)) }
