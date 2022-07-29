package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var sizeMap = make(map[int64]int)
var nameMap = make(map[int]string)

func main() {
	sizeMap[19606442] = 1222543
	nameMap[1222543] = "Lagtrain"
	var wr = bufio.NewWriter(os.Stdout)
	defer wr.Flush()

	b, err := ioutil.ReadFile("dir.txt")
	if err != nil {
		panic(err)
	}
	root := string(b)
	dirs, err := ioutil.ReadDir(root) // music dirs
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			ss := strings.SplitN(dir.Name(), " ", 2)
			setId, err := strconv.Atoi(ss[0])
			if err != nil {
				// fmt.Printf("%s: %s\n", err, dir.Name())
				continue
			}
			if len(ss) < 2 {
				fmt.Printf("invalid dir name format: %s\n", dir.Name())
				continue
			}
			ss = strings.Split(ss[1], " - ")
			var name string
			if len(ss) < 2 {
				name = ss[0]
			} else {
				name = ss[1]
			}
			nameMap[setId] = name
			size, err := dirSize(filepath.Join(root, dir.Name()))
			if err != nil {
				fmt.Printf("%s: %s\n", err, dir.Name())
				continue
			}

			if id, ok := sizeMap[size]; ok {
				t := fmt.Sprintf("%dbyte: %d(%s), %d(%s)", size, setId, name, id, nameMap[id])
				fmt.Fprintln(wr, t)
			} else {
				sizeMap[size] = setId
			}
		}
	}
}
func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
