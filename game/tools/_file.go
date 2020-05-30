package tools

import "os"

const CacheDir = "./cached/"

var ModePrefix = map[int]string{0: "o", 1: "t", 2: "c", 3: "m"}

func SetCacheDir() {
	SetDir(CacheDir)
}
func SetDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0644)
	}
}
