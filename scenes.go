package gosu

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/play"
)

type Scenes struct {
	scenes []scene.Scene
	idx    int
}

func (scs Scenes) Scene() scene.Scene { return scs.scenes[scs.idx] }

const (
	SceneChoose = iota
	ScenePlay
)

func (scs *Scenes) Update() error {
	if scs.Scene == nil {
		scs.scenes = scs.scenes["choose"]
	}

	sc := scs.Scene()
	switch args := sc.Update().(type) {
	case error:
		fmt.Println("play scene error:", args)
		scs.scene = scs.scenes["choose"]
	case piano.Scorer:
		ebiten.SetWindowTitle("gosu")
		// debug.SetGCPercent(100)
		scs.scene = scs.scenes["choose"]
	case scene.PlayArgs:
		fsys := args.MusicFS
		name := args.ChartFilename
		replay := args.Replay
		scene, err := play.NewScene(g.Config, g.Asset, fsys, name, replay)
		if err != nil {
			fmt.Println("play scene error:", args)
			scs.scene = scs.scenes["choose"]
		} else {
			// debug.SetGCPercent(0)
			scs.scene = scene
		}
	}
	return nil
}

func (scs Scenes) Draw(screen draws.Image) {
	scs.Scene().Draw(screen)
}
