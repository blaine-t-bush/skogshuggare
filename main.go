package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/gdamore/tcell"
)

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
	}

	// Set default style and clear terminal.
	s.SetStyle(tcell.StyleDefault)
	s.Clear()

	// Draw borders.
	w, h := s.Size()
	b := Border{1, 0, w - 1, 0, h - 1}
	b.Draw(s)

	// Draw trees.
	t1 := Tree{10, 10}
	t2 := Tree{15, 5}
	t3 := Tree{16, 8}
	t4 := Tree{20, 15}
	t5 := Tree{25, 10}
	t1.Draw(s)
	t2.Draw(s)
	t3.Draw(s)
	t4.Draw(s)
	t5.Draw(s)

	// Initialize player.

	// Wait for Loop() goroutine to finish before moving on.
	var wg sync.WaitGroup
	wg.Add(1)
	go Loop(&wg, s)
	wg.Wait()

	s.Fini()
}

func Loop(wg *sync.WaitGroup, s tcell.Screen) {
	defer wg.Done()
	player := Player{5, 5}
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyUp:
				player.Move(s, 1, 0)
			case tcell.KeyRight:
				player.Move(s, 1, 1)
			case tcell.KeyDown:
				player.Move(s, 1, 2)
			case tcell.KeyLeft:
				player.Move(s, 1, 3)
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}
