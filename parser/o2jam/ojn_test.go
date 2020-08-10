package o2jam

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	paths := []string{"test/o2ma204.ojn"}
	for _, path := range paths {
		ojn, _ := Parse(path)
		// tools.PrettyPrint(ojn.Charts[2])
		for _, s := range ojn.Charts[2].Notes {
			fmt.Println(s.Measure, s.Channel, s.EventCount)
		}
		// fmt.Printf("%+v\n", ojn)
	}
}
