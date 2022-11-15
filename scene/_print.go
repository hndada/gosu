package scene

import (
	"fmt"
	"strings"
)

// These codes are for generating boilerplate.
func PrintCurrent(settings string) {
	fmt.Println("func (Settings) Current() Settings {")
	fmt.Println("	return Settings{")
	for _, line := range strings.Split(settings, "\n") {
		line = strings.TrimSpace(line)
		words := strings.Split(line, " ")
		if len(words) < 2 {
			fmt.Println()
			continue
		}
		field := words[0]
		field = strings.ToLower(string(field[0])) + field[1:]
		fmt.Printf("\t\t%s: %s,\n", words[0], field)
	}
	fmt.Println("	}")
	fmt.Print("}")
}
func PrintSet(settings string) {
	fmt.Println("func (Settings) Set(s Settings) {")
	for _, line := range strings.Split(settings, "\n") {
		line = strings.TrimSpace(line)
		words := strings.Split(line, " ")
		if len(words) < 2 {
			fmt.Println()
			continue
		}
		field := words[0]
		field = strings.ToLower(string(field[0])) + field[1:]
		fmt.Printf("\t%s = s.%s\n", field, words[0])
	}
	fmt.Print("}")
}
