package choose

type sortBy int

const (
	sortByMusicName sortBy = iota
	sortByLevel
	sortByTime
	sortByAddAtTime
)

func sortCharts(src []Chart, sortBy sortBy) []Chart {
	cs := make([]Chart, len(src))
	copy(cs, src)

	var less func(i, j int) bool
	switch sortBy {
	case sortByMusicName:
		less = func(i, j int) bool {
			if cs[i].MusicName < cs[j].MusicName {
				return true
			} else if cs[i].MusicName > cs[j].MusicName {
				return false
			}
			if cs[i].Artist < cs[j].Artist {
				return true
			} else if cs[i].Artist > cs[j].Artist {
				return false
			}
			return cs[i].Level < cs[j].Level
		}
	case sortByLevel:
		less = func(i, j int) bool {
			return cs[i].Level < cs[j].Level
		}
	case sortByTime:
		less = func(i, j int) bool {
			if cs[i].Duration < cs[j].Duration {
				return true
			} else if cs[i].Duration > cs[j].Duration {
				return false
			}
			return cs[i].Level < cs[j].Level
		}
	case sortByAddAtTime:
		less = func(i, j int) bool {
			return cs[i].AddAtTime.Before(cs[j].AddAtTime)
		}
	}

	sort.Slice(cs, less)
	return cs
}

// Group1, Group2, Sort, Filter int
const (
	listDepthSearch = -1
	listDepthMusic  = iota
	listDepthChart
	// FocusKeySettings
)

// var NewSliceIndexKeyHandler func(index *int, len int) ctrl.KeyHandler
func (s Scene) NewListMoveKeyHandler(handler ctrl.Handler) ctrl.KeyHandler {
	return ctrl.KeyHandler{
		Handler:  handler,
		Modifier: input.KeyNone,
		Keys:     scene.UpDownKeys,
		Sounds:   [2]audios.SoundPlayer{s.SwipeSoundPod, s.SwipeSoundPod},
		Volume:   &s.SoundVolume,
	}
}
