package parser

import (
	"fmt"
	"github.com/hndada/gosu/tools"
	"testing"
)

func TestParseOJN(t *testing.T) {
	paths := []string{"test/o2ma204.ojn"}
	for _, path := range paths {
		ojn := *ParseOJN(path)
		fmt.Printf("%x\n", ojn.Artist[:])
		tools.PrettyPrint(ojn)
		fmt.Println(string(ojn.Title[:]), (ojn.Artist[:]))
		// fmt.Printf("%+v\n", ojn)
	}
}
