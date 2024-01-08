//go:build !js && !wasm

package vfs

import "os"

var DirFS = os.DirFS
