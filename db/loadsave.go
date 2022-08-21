package db

import (
	"fmt"
	"os"
)

// This is not a database, just a data.

// dst should be pointer.
func LoadData(fpath string, dst any) error {
	b, err := os.ReadFile(fpath)
	if err != nil {
		_ = os.Rename(fpath, fpath+".crashed") // Try rename if exists.
		return err
	}
	err = Unmarshal(b, dst)
	if err != nil {
		_ = os.Rename(fpath, fpath+".crashed") // Try rename if exists.
		return err
	}
	return nil
}

// src should be pointer.
// Todo: error notification?
func SaveData(fpath string, src any) {
	b, err := Marshal(src)
	if err != nil {
		fmt.Printf("Failed to save data by an error: %s", err)
		return
	}
	err = os.WriteFile(fpath, b, 0644)
	if err != nil {
		fmt.Printf("Failed to save data by an error: %s", err)
		return
	}
}
