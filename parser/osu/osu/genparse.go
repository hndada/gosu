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

// ScanStructs supposes gofmt at given file was already proceeded
func ScanStructs(path string) ([]string, map[string][]fieldInfo) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	structs := make([]string, 0)
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
			infos = make([]fieldInfo, 0)
			structs = append(structs, structName)
		case structName != "" && len(vs) >= 2:
			inf := fieldInfo{name: vs[0], fieldType: vs[1]}
			inf.delimiter = make([]string, 0)
			for i := 0; i < strings.Count(vs[1], "["); i++ {
				delimiter := strings.TrimLeft(vs[3+2*i], `delimiter`)
				if delimiter == "(space)" {
					delimiter = " "
				}
				inf.delimiter = append(inf.delimiter, delimiter)
			}
			infos = append(infos, inf)
		case vs[0] == "}":
			m[structName] = infos
			structName = ""
		}
	}
	return structs, m
}

func PrintSetValue(fields []fieldInfo, lname string, valName string, isPtr bool) {
	var ptrmark string
	if isPtr {
		ptrmark = "&"
	}
	for _, f := range fields {
		switch f.fieldType {
		case "string":
			fmt.Printf(`
	%s.%s = %s
`, lname, f.name, valName)
		case "int":
			fmt.Printf(`
	i, err := strconv.Atoi(v)
	if err != nil {
			return %s%s, err
		}
	%s.%s = %s
`, ptrmark, lname, lname, f.name, valName)
		case "float64":
			fmt.Printf(`
	f, err := strconv.ParseFloat(%s, 64)
	if err != nil {
			return %s%s, err
		}
	%s.%s = f
`, valName, ptrmark, lname, lname, f.name)
		case "bool":
			fmt.Printf(`
	b, err := strconv.ParseBool(%s)
	if err != nil {
		return %s%s, err
		}
	%s.%s = b
`, valName, ptrmark, lname, lname, f.name)
		case "[]string":
			fmt.Printf(`
	slice := make([]string, 0)
	for _, s := range strings.Split(%s, "%s") {
		slice = append(slice, s)
	}
	%s.%s = slice
`, valName, f.delimiter[0], lname, f.name)
		case "[]int":
			fmt.Printf(`
	slice := make([]int, 0)
	for _, s := range strings.Split(%s, "%s") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return %s%s, err
		}
		slice = append(slice, i)
	}
	%s.%s = slice
`, valName, f.delimiter[0], ptrmark, lname, lname, f.name)
		}
	}
}

// generate local name
// example: TimingPoint -> tp
func localName(structName string) string {
	var name string
	for i, s := range strings.ToLower(structName) {
		if structName[i] != byte(s) {
			name += string(s)
		}
	}
	return name
}

func main() {
	structs, m:=ScanStructs("format.go")
	for _, s:=range structs{
		switch s {
		case "General", "Editor", "Metadata", "Difficulty":
			PrintSetValue(m[s], localName(s), "kv[1]", false)
		case "TimingPoint", "HitObject", "SliderParams", "HitSample":
			PrintSetValue(m[s], localName(s), "v", false)
		}
		// fmt.Printf("%s: %q\n", s, m[s])
	}
}

// type typeCode int

// const (
// 	typeString = iota
// 	typeInt
// 	typeFloat64
// 	typeBool
//
// 	typeStringSlice
// 	typeIntSlice
// )
