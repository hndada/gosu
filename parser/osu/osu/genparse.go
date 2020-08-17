// +build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// too many inconsistent pattern in .osu file format, very hard to write fully-generating code
type fieldInfo struct {
	name      string
	fieldType string
	delimiter []string
}

// ScanStructs supposes gofmt was already proceeded at given file
func ScanStructs(path string) ([]string, map[string]string, map[string][]fieldInfo) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	structs := make([]string, 0)
	delimiters := make(map[string]string)
	m := make(map[string][]fieldInfo)
	scanner := bufio.NewScanner(f)
	var structName string
	var infos []fieldInfo
	for scanner.Scan() {
		vs := strings.Fields(scanner.Text())
		switch {
		case len(vs) == 0 || vs[0] == "//" || vs[len(vs)-1] == "manual": // maybe panic won't happen
			continue

		case vs[0] == "type" && len(vs) > 2 && vs[2] == "struct":
			structName = vs[1]
			structs = append(structs, structName)
			if strings.HasPrefix(vs[len(vs)-1], `delimiter`) {
				delimiter := strings.TrimLeft(vs[len(vs)-1], `delimiter`)
				delimiters[structName] = delimiter
			}
			infos = make([]fieldInfo, 0)

		case structName != "" && len(vs) >= 2:
			info := fieldInfo{name: vs[0], fieldType: vs[1]}
			if len(vs) >= 4 && vs[3] == "nofloat" {
				info.fieldType += " (nofloat)"
			}
			info.delimiter = make([]string, 0)
			for i := 0; i < strings.Count(vs[1], "["); i++ {
				delimiter := strings.TrimLeft(vs[3+2*i], `delimiter`)
				if delimiter == "(space)" {
					delimiter = " "
				}
				info.delimiter = append(info.delimiter, delimiter)
			}
			infos = append(infos, info)

		case vs[0] == "}":
			m[structName] = infos
			structName = ""

		}
	}
	for k := range delimiters {
		delimiters[k] = strings.Replace(delimiters[k], "(space)", " ", -1)
	}
	return structs, delimiters, m
}

func PrintSetValue(valName, returnName, localName string, f fieldInfo) {
	switch f.fieldType {
	case "string":
		fmt.Printf(`
	%s.%s = %s
`, localName, f.name, valName)
	case "int":
		fmt.Printf(`
	f, err := strconv.ParseFloat(%s, 64)
	if err != nil {
			return %s, err
		}
	%s.%s = int(f)
`, valName, returnName, localName, f.name)
	case "int (nofloat)":
		fmt.Printf(`
	i, err := strconv.Atoi(%s)
	if err != nil {
			return %s, err
		}
	%s.%s = i
`, valName, returnName, localName, f.name)
	case "float64":
		fmt.Printf(`
	f, err := strconv.ParseFloat(%s, 64)
	if err != nil {
			return %s, err
		}
	%s.%s = f
`, valName, returnName, localName, f.name)
	case "bool":
		fmt.Printf(`
	b, err := strconv.ParseBool(%s)
	if err != nil {
		return %s, err
		}
	%s.%s = b
`, valName, returnName, localName, f.name)
	case "[]string":
		fmt.Printf(`
	slice := make([]string, 0)
	for _, s := range strings.Split(%s, "%s") {
		slice = append(slice, s)
	}
	%s.%s = slice
`, valName, f.delimiter[0], localName, f.name)
	case "[]int":
		fmt.Printf(`
	slice := make([]int, 0)
	for _, s := range strings.Split(%s, "%s") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return %s, err
		}
		slice = append(slice, i)
	}
	%s.%s = slice
`, valName, f.delimiter[0], returnName, localName, f.name)
	}
}

// generate compressed name as local name; ex) TimingPoint -> tp
func genLocalName(structName string) string {
	var name string
	for i, s := range strings.ToLower(structName) {
		if structName[i] != byte(s) {
			name += string(s)
		}
	}
	return name
}

func main() {
	structs, delimiters, m := ScanStructs("format.go")
	for _, structName := range structs {
		switch structName {
		case "General", "Editor", "Metadata", "Difficulty":
			// PrintSetValue(m[s], s, delimiters[s], "section")
			localName := "o." + structName
			valName := "kv[1]"
			returnName := "&o"

			fmt.Printf("case \"%s\":\n", structName)
			fmt.Printf("kv := strings.Split(line, `%s`)\n", delimiters[structName])
			fmt.Printf("switch kv[0] {\n")
			for _, f := range m[structName] {
				fmt.Printf("case \"%s\":", f.name)
				PrintSetValue(valName, returnName, localName, f)
			}
			fmt.Printf("}\n")

		case "TimingPoint", "HitObject", "SliderParams", "HitSample":
			// PrintSetValue(m[s], s, delimiters[s], "substruct")
			var valName string
			localName := genLocalName(structName)
			returnName := localName

			fmt.Printf("\nfunc new%s(line string) (%s, error) {\n", structName, structName)
			fmt.Printf("var %s %s\n", localName, structName)
			fmt.Printf("vs := strings.Split(line, `%s`)\n", delimiters[structName])
			for i, f := range m[structName] {
				fmt.Printf("{")
				valName = fmt.Sprintf("vs[%d]", i)
				PrintSetValue(valName, returnName, localName, f)
				fmt.Printf("}\n")
			}
			fmt.Printf("return %s, nil\n}", localName)
			fmt.Println()
		}
	}
}
