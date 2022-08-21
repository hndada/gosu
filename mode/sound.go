package mode

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hndada/gosu/audios"
)

var Sounds audios.SoundMap

func LoadSounds(soundRoot string) error {
	Sounds = audios.NewSoundMap(&Volume)
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
