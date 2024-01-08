//go:build js || wasm

// Code for JavaScript and WebAssembly

package vfs

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed host.txt
var host string

// vfs (Virtual File System): If your implementation abstracts
// multiple file systems or creates a virtual layer over existing ones.
func DirFS(name string) fs.FS {
	url := fmt.Sprintf("%s/%s", host, name)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "prefix")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer os.Remove(tmpFile.Name()) // clean up

	// Write the response body to the temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Now you can use os.DirFS to get an fs.FS that represents
	// the directory containing the temporary file
	dir := filepath.Dir(tmpFile.Name())
	fs := os.DirFS(dir)
}
