// +build ignore

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// 필요할 때만 ``풀기
// src1, 삽입, src2

func main() {
	// Events, TimingPoints, HitObjects는 어떤 시점에서 자동 스킵됨
	sections := []string{"General", "Editor", "Metadata", "Difficulty"} // [4]string
	fieldNameTypes := make(map[string][][2]string)                      // [4][][2]string
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "format.go", nil, 0)
	if err != nil {
		panic(err)
	}
	for _, node := range f.Decls {
		if genDecl, ok := node.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					structName := typeSpec.Name.Name
					switch structName {
					case "General", "Editor", "Metadata", "Difficulty":
					default:
						continue
					}
					fieldNameTypes[structName] = make([][2]string, 0)
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						for _, field := range structType.Fields.List {
							for _, name := range field.Names {
								fieldName := name.Name
								var fieldType string
								switch field.Type.(type) {
								case *ast.Ident:
									fieldType = field.Type.(*ast.Ident).Name
								case *ast.ArrayType:
									a := field.Type.(*ast.ArrayType)
									fieldType = fmt.Sprintf("[]%s", a.Elt.(*ast.Ident).Name)
								}
								fieldNameTypes[structName] = append(fieldNameTypes[structName],
									[2]string{fieldName, fieldType})
								// fmt.Printf("%s: %s %s\n", structName, fieldName, fieldType)
							}
						}
					}
				}
			}
		}
	}
	// template:
	// case "General":
	// kv := strings.Split(line, `: `)
	// switch kv[0] {
	// case "AudioFilename":
	// 	o.General.AudioFilename = kv[1]
	// }

	for _, section := range sections {
		fmt.Printf("case \"%s\":\n", section)
		var delimiter string
		switch section {
		case "General", "Editor":
			delimiter = `: `
		case "Metadata", "Difficulty":
			delimiter = `:`
		}
		fmt.Printf("kv := strings.Split(line, `%s`)\n", delimiter)
		fmt.Printf("switch kv[0] {\n")
		for _, field := range fieldNameTypes[section] {
			fieldName, fieldType := field[0], field[1]
			fmt.Printf("case \"%s\":", fieldName)
			switch fieldType {
			case "string":
				PrintSetString(section, fieldName)
			case "int":
				PrintSetInt(section, fieldName)
			case "float64":
				PrintSetFloat64(section, fieldName)
			case "bool":
				PrintSetBool(section, fieldName)
			case "[]int":
				switch fieldName {
				case "Bookmarks":
					PrintSetIntSlice(section, fieldName, ",")
				}
			case "[]string":
				switch fieldName {
				case "Tags":
					PrintSetStringSlice(section, fieldName, " ")
				}
			}
		}
		fmt.Printf("}\n")
	}
}

func PrintSetString(section, fieldName string) {
	fmt.Printf(`
	o.%s.%s = kv[1]
`, section, fieldName)
}

func PrintSetInt(section, fieldName string) {
	fmt.Printf(`
	i, err := strconv.Atoi(kv[1])
	if err != nil {
			return &o, err
		}
	o.%s.%s = i
`, section, fieldName)
}

func PrintSetFloat64(section, fieldName string) {
	fmt.Printf(`
	f, err := strconv.ParseFloat(kv[1], 64)
	if err != nil {
			return &o, err
		}
	o.%s.%s = f
`, section, fieldName)
}

func PrintSetBool(section, fieldName string) {
	fmt.Printf(`
	b, err := strconv.ParseBool(kv[1])
	if err != nil {
			return &o, err
		}
	o.%s.%s = b
`, section, fieldName)
}

func PrintSetIntSlice(section, fieldName, delimiter string) {
	fmt.Printf(`
	slice := make([]int, 0)
	for _, s := range strings.Split(kv[1], "%s") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return &o, err
		}
		slice = append(slice, i)
	}
	o.%s.%s = slice
`, delimiter, section, fieldName)
}

func PrintSetStringSlice(section, fieldName, delimiter string) {
	fmt.Printf(`
	slice := make([]string, 0)
	for _, s := range strings.Split(kv[1], "%s") {
		slice = append(slice, s)
	}
	o.%s.%s = slice
`, delimiter, section, fieldName)
}
