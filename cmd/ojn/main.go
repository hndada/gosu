package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hndada/gosu/format/o2jam"
)

// Example usage (not built by default):
func main() {
	f, _ := os.Open("o2ma111.ojn")
	defer f.Close()
	ojn, err := o2jam.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ojn.Header.Title(), ojn.Header.Artist())
	notes := ojn.Hard.ScanChannelEvents(o2jam.ChanLane4)
	fmt.Printf("Hard middle-lane events: %d\n", len(notes))

	// Convert to *osu.Format for further processing (e.g., loading as piano.Chart).
	of, err := o2jam.ConvertOJNToOsuFormat(ojn, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Converted to osu.Format: %d hit objects, duration %.1f sec\n",
		len(of.HitObjects), float64(of.HitObjects[len(of.HitObjects)-1].EndTime)/1000.0)
	fmt.Println("First 5 hit objects:")
	for i, ho := range of.HitObjects {
		if i >= 5 {
			break
		}
		fmt.Printf("  Time %d ms, X %d, Y %d\n", ho.Time, ho.X, ho.Y)
	}
}
