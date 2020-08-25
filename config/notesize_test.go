package config

import (
	"fmt"
	"testing"
)

func TestNoteSize(t *testing.T) {
	s:=newSettings()
	s.setNoteSizes()
	fmt.Printf("%+v, %+v\n", NoteSizes[4], NoteSizes[7])
}
