package scene

import "github.com/hndada/gosu/ui"

// Resources, Options, States.
// Resources are loaded from file system.
// Options are set by user and saved to file system.
// States are generated when runtime, and not saved.
type States struct {
	KeyboardStatus *ui.KeyboardStatus
}
