package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

// type Settings struct {
// 	Age int
// }

var Age int = 35

func main() {
	fmt.Printf("Age before: %d\n", Age)
	var data = `Age = 25`
	var settings map[string]interface{}
	_, err := toml.Decode(data, &settings)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range settings {
		switch k {
		case "Age":
			age, ok := v.(int64)
			if ok {
				Age = int(age)
			}
		}
	}
	fmt.Printf("Age after: %d\n", Age)
}
