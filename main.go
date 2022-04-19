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

	s.SetStyle(tcell.StyleDefault)
	s.Clear()

	// Wait for Loop() goroutine to finish before moving on.
	var wg sync.WaitGroup
	wg.Add(1)
	go Loop(&wg, s)
	wg.Wait()

	s.Fini()
}

func Loop(wg *sync.WaitGroup, s tcell.Screen) {
	defer wg.Done()
	loc := Point{5, 5}
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyEnter:
				loc = DrawPlayer(s, loc, false, 0)
			case tcell.KeyUp:
				loc = DrawPlayer(s, loc, true, 0)
			case tcell.KeyRight:
				loc = DrawPlayer(s, loc, true, 1)
			case tcell.KeyDown:
				loc = DrawPlayer(s, loc, true, 2)
			case tcell.KeyLeft:
				loc = DrawPlayer(s, loc, true, 3)
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func DrawPlayer(s tcell.Screen, loc Point, isMoving bool, dir int) Point {
	st := tcell.StyleDefault
	gl := '@'
	var newLoc Point
	s.Clear()
	if isMoving {
		switch dir {
		case 0: // up
			newLoc = Point{loc.x, loc.y - 1}
		case 1: // right
			newLoc = Point{loc.x + 1, loc.y}
		case 2: // down
			newLoc = Point{loc.x, loc.y + 1}
		case 3: // left
			newLoc = Point{loc.x - 1, loc.y}
		}
	} else {
		newLoc = loc
	}
	s.SetCell(newLoc.x, newLoc.y, st, gl)
	s.Show()
	return newLoc
}
