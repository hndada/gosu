// +build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type fieldInfo struct {
	name      string
	fieldType string
}

// ScanStructs supposes gofmt was already proceeded at given file
// copy&pasted function
func ScanStructs(path string) ([]string, map[string][]fieldInfo) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var structName string
	var infos []fieldInfo
	structs := make([]string, 0)
	m := make(map[string][]fieldInfo)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		vs := strings.Fields(scanner.Text())
		switch {
		case len(vs) == 0 || vs[0] == "//":
			continue
		case vs[0] == "type" && len(vs) > 2 && vs[2] == "struct":
			structName = vs[1]
			structs = append(structs, structName)
			infos = make([]fieldInfo, 0)
		case structName != "" && len(vs) >= 2:
			info := fieldInfo{name: vs[0], fieldType: vs[1]}
			infos = append(infos, info)
		case vs[0] == "}":
			m[structName] = infos
			structName = ""
		}
	}
	return structs, m
}

func printSetSkinValue(f fieldInfo) {
	switch f.fieldType {
	case "[10]*ebiten.Image":
		fmt.Printf(`for i := 0; i < 10; i++ {
	filename = fmt.Sprintf(, i)
	path = filepath.Join(skinPath, filename)
	if s.%s[i], ok = LoadImage(path); !ok {
		s.%s[i] = defaultSkin.%s[i]
	}
}
`, f.name, f.name, f.name)
	case "*ebiten.Image":
		fmt.Printf(`	filename = 
	path = filepath.Join(skinPath, filename)
	if s.%s, ok = LoadImage(path); !ok {
		s.%s = defaultSkin.%s
}
`, f.name, f.name, f.name)
	}
}

func printSetManiaSkinValue(f fieldInfo) {
	switch f.fieldType {
	case "[4]*ebiten.Image":
		for i := 0; i < 4; i++ {
			fmt.Printf(`filename = 
path = filepath.Join(skinPath, filename)
if s.%s[%d], ok = LoadImage(path); !ok {
s.%s[%d] = defaultSkin.mania.%s[%d]
}
`, f.name, i, f.name, i, f.name, i)
		}
	case "[5]*ebiten.Image":
		for i := 0; i < 5; i++ {
			fmt.Printf(`filename = 
path = filepath.Join(skinPath, filename)
if s.%s[%d], ok = LoadImage(path); !ok {
s.%s[%d] = defaultSkin.mania.%s[%d]
}
`, f.name, i, f.name, i, f.name, i)
		}
	case "[4][]*ebiten.Image":
		for i := 0; i < 4; i++ {
			fmt.Printf(`	for i := range s.%s {
		s.%s[i] = make([]*ebiten.Image, 1, 1)
	}
filename = 
path = filepath.Join(skinPath, filename)
if s.%s[%d][0], ok = LoadImage(path); !ok {
s.%s[%d][0] = defaultSkin.mania.%s[%d][0]
}
`, f.name, f.name, f.name, i, f.name, i, f.name, i)
		}
	case "*ebiten.Image":
		fmt.Printf(`filename = 
path = filepath.Join(skinPath, filename)
if s.%s, ok = LoadImage(path); !ok {
	s.%s = defaultSkin.mania.%s
}
`, f.name, f.name, f.name)
	case "[]*ebiten.Image":
		fmt.Printf(`s.%s = make([]*ebiten.Image, 1)
filename = 
path = filepath.Join(skinPath, filename)
if s.%s[0], ok = LoadImage(path); !ok {
s.%s[0] = defaultSkin.mania.%s[0]'
}
`, f.name, f.name, f.name, f.name)
	}
}

func main() {
	_, m := ScanStructs("skin.go")
	fmt.Printf(`func (s *skin) LoadSkin(skinPath string) {
var filename, path string
	var ok bool
`)
	for _, f := range m["skin"] {
		printSetSkinValue(f)
	}
	fmt.Printf(`s.mania.load(skinPath)
}
`)
	fmt.Printf(`func (s *maniaSkin) load(skinPath string) {
var filename, path string
	var ok bool
`)
	for _, f := range m["maniaSkin"] {
		printSetManiaSkinValue(f)
	}
	fmt.Printf("}\n")
}
