package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
	"github.com/eiannone/keyboard"
)

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()

	width := float32(cw)
	height := float32(ch)

	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	fmt.Println("hello")

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

	go func() {
		w := app.NewWindow(app.Title("Skogshuggare"), app.Size(unit.Px(width), unit.Px(height)))
		if err := CreateCanvas(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func CreateCanvas(w *app.Window, width, height float32) error {
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
			CreateBlock(canvas)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func CreateBlock(canvas *giocanvas.Canvas) {
	canvas.CenterRect(50, 50, 10, 10, color.NRGBA{100, 0, 0, 255})
}
