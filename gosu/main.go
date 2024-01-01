package main

import (
	"embed"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
)

//go:embed resources/* music/* replay/*
var defaultFS embed.FS

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root := os.DirFS(dir)
	rootPaths := gosu.NewRootPaths()

	g := gosu.NewGame(root)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
