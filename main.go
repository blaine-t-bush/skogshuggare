package main

import (
	"fmt"
	"os"
	"time"

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

	quit := make(chan struct{})
	dir := make(chan tcell.Key)
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
					close(quit)
					return
				case tcell.KeyUp, tcell.KeyDown, tcell.KeyRight, tcell.KeyLeft:
					dir <- ev.Key()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()
	beep(s, quit, dir)
	s.Fini()
}

func beep(s tcell.Screen, quit <-chan struct{}, dir <-chan tcell.Key) {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-quit:
			return
		case <-dir:
			fmt.Println(<-dir)
		case <-t.C:
			s.Beep()
		}
	}
}
