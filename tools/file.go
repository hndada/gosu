package tools

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 폴더 두번 스캔
func LoadSongList(root, ext string) ([]string, error) {
	var songs []string
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		return songs, errors.New("invalid root dir")
	}
	sets, err := ioutil.ReadDir(root)
	if err != nil {
		return songs, errors.New("invalid root dir")
	}

	var absSet, absMap string
	for _, set := range sets {
		absSet = filepath.Join(root, set.Name())
		if info, err := os.Stat(absSet); err != nil || !info.IsDir() {
			continue
		}
		maps, err := ioutil.ReadDir(absSet)
		if err != nil {
			continue
		}
		for _, mapFile := range maps {
			absMap = filepath.Join(absSet, mapFile.Name())
			if !mapFile.IsDir() && filepath.Ext(mapFile.Name()) == ext {
				songs = append(songs, absMap)
			}
		}
	}
	return songs, nil
}
