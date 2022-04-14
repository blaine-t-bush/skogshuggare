package main

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	// github.com/JoelOtter/termloop
)

func main() {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				char, _, err := keyboard.GetSingleKey()
				if err != nil {
					panic(err)
				}
				fmt.Printf("You pressed: %q\r\n", char)
			}
		}
	}()

	for {
		// Loop infinitely to listen for keyboard events
	}
}
