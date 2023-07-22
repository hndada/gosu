func main() {
	musicsFS, err := fs.Sub(fsys, musicRoot)
	if err != nil {
		panic(err)
	}
	g.loadTestPiano(cfg, asset, musicsFS)
}

func (g *game) loadTestPiano(cfg *scene.Config, asset *scene.Asset, musicsFS fs.FS) {
	musicName := "nekodex - circles!"
	musicFS, err := fs.Sub(musicsFS, musicName)
	// musicFS := ZipFS(filepath.Join(dir, musicName+".osz"))
	if err != nil {
		panic(err)
	}
	name := "nekodex - circles! (MuangMuangE) [Hard].osu"

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fsys := os.DirFS(dir)

	replay, err := mode.NewReplay(fsys, "format/osr/testdata/circles(7k).osr", 7)
	if err != nil {
		panic(err)
	}

	scenePlay, err := play.NewScene(cfg, asset, musicFS, name, replay)
	if err != nil {
		panic(err)
	}
	g.scene = scenePlay
}