package gosu

import "io/fs"

type RootPaths struct {
	ResourcesList []string
	Musics        []string
	Replays       []string
}

func NewRootPaths() RootPaths {
	return RootPaths{
		ResourcesList: []string{"resources"},
		Musics:        []string{"music"},
		Replays:       []string{"replay"},
	}
}

type Root struct {
	ResourcesList []fs.FS
	Musics        []fs.FS
	Replays       []fs.FS
}

func NewRoot(root fs.FS, rootPaths RootPaths) (r Root) {
	r.ResourcesList = make([]fs.FS, len(rootPaths.ResourcesList))
	for i, path := range rootPaths.ResourcesList {
		fsys, err := fs.Sub(root, path)
		if err != nil {
			continue
		}
		r.ResourcesList[i] = fsys
	}

	r.Musics = make([]fs.FS, len(rootPaths.Musics))
	for i, path := range rootPaths.Musics {
		fsys, err := fs.Sub(root, path)
		if err != nil {
			continue
		}
		r.Musics[i] = fsys
	}

	r.Replays = make([]fs.FS, len(rootPaths.Replays))
	for i, path := range rootPaths.Replays {
		fsys, err := fs.Sub(root, path)
		if err != nil {
			continue
		}
		r.Replays[i] = fsys
	}
	return r
}
