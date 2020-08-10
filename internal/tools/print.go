package tools

import (
	"encoding/json"
	"fmt"
)

// print the contents of the obj
func PrettyPrint(data interface{}) {
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
