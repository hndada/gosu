package mania

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestProcess(t *testing.T) {
	files, err := ioutil.ReadDir("test_full")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		cwd, _ := os.Getwd()
		path := filepath.Join(cwd, "test_full", f.Name())
		fmt.Println(f.Name())
		process(path)
	}
}
