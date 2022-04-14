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
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				fmt.Println(ev.Key())
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
					close(quit)
					return
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()
	beep(s, quit)
	s.Fini()
}

func beep(s tcell.Screen, quit <-chan struct{}) {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-quit:
			return
		case <-t.C:
			s.Beep()
		}
	}
}
