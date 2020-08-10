package taiko

const modeTaiko = element.ModeTaiko

type TaikoBeatmap struct {
	diff.Beatmap
	Notes    []TaikoNote
	OldNotes []TaikoOldNote
}

func (beatmap *TaikoBeatmap) SetBase(path string, modsBits int) {
	beatmap.Beatmap = element.ParseBeatmap(path)
	diff.CheckMode(beatmap.Mode, modeTaiko)
	beatmap.Mods = diff.GetMods(modeTaiko, modsBits)
}
