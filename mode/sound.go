package mode

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hndada/gosu/audioutil"
)

var Sounds audioutil.SoundMap

func LoadSounds(soundRoot string) error {
	Sounds = audioutil.NewSoundMap(&Volume)
	fs, err := os.ReadDir(soundRoot)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		spath := filepath.Join(soundRoot, f.Name())
		err = Sounds.Register(spath)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
