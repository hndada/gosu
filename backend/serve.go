package vfs

import (
	"net/http"
	"os"
)

func Serve() {
	const root = "."
	fsys := http.FS(os.DirFS(root))
	fs := http.FileServer(fsys)
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
