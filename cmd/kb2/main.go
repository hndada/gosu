package main

import (
	"fmt"
	"time"

	"github.com/hndada/gosu/engine/kb"
)

func main() {
	{
		startTime := time.Now()
		go kb.Listen()
		for time.Since(startTime).Seconds() < 2 {
			time.Sleep(1 * time.Microsecond) // prevents 100% CPU usage
		}
		events := kb.Fetch()
		fmt.Println("1st")
		fmt.Printf("%v\n", events)
		for time.Since(startTime).Seconds() < 4 {
			time.Sleep(1 * time.Microsecond) // prevents 100% CPU usage
		}
		fmt.Println("2nd")
		events = kb.Fetch()
		fmt.Printf("%v\n", events)
		kb.Exit()
	}
	{
		startTime := time.Now()
		go kb.Listen()
		for time.Since(startTime).Seconds() < 2 {
			time.Sleep(1 * time.Microsecond) // prevents 100% CPU usage
		}
		events := kb.Fetch()
		fmt.Println("3rd")
		fmt.Printf("%v\n", events)
		for time.Since(startTime).Seconds() < 4 {
			time.Sleep(1 * time.Microsecond) // prevents 100% CPU usage
		}
		fmt.Println("4th")
		events = kb.Fetch()
		fmt.Printf("%v\n", events)
		kb.Exit()
	}
}
