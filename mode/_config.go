package mode

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// username
// token
// savePw

type Config struct {
	Skin         int
	GroupingMode int
	CurrentTab   int
	SfxVolume    float64
	BgmVolume    float64
	DimLv        float64
}

var singleton = NewConfig()
// todo: singleton 대신 game에 부여하기.

func NewConfig() *Config {
	config := &Config{
		Skin:         -1,
		GroupingMode: 0,
		CurrentTab:   0,
		SfxVolume:    0.3,
		BgmVolume:    0.3,
		DimLv:        0.8,
	}
	return config
}

// TODO: separate common json-relating part as general func
func load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, singleton); err != nil {
		return err
	}
	return nil
}

func save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := json.MarshalIndent(singleton, "", "    ")
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
