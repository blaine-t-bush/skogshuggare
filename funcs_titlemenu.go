package main

import (
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell"
)

func (titleMenu *TitleMenu) Draw(screen tcell.Screen) {
	screen.Clear()
	time.Sleep(100 * time.Millisecond)
	titleMenu.DrawPage(screen, titleMenu.pageState)
	screen.Show()
}

func (titleMenu *TitleMenu) DrawPage(screen tcell.Screen, pageState int) {
	currentPage := titleMenu.titleMenuPages[pageState]
	pageContent := currentPage.content
	//pageItems := currentPage.titleMenuItems

	maxY := strings.Count(pageContent, "\n")
	maxX := strings.Index(pageContent, "\n")
	runeIndex := 0
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			screen.SetContent(x, y, rune(pageContent[runeIndex]), nil, tcell.StyleDefault)
			runeIndex++
		}
	}
}

func (titleMenu *TitleMenu) Update(screen tcell.Screen) {

}

func (titleMenu *TitleMenu) InputHandler(screen tcell.Screen) {
	ev := screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			os.Exit(0)
		case tcell.KeyUp:
			titleMenu.MoveMenuCursor(DirUp)
		case tcell.KeyDown:
			titleMenu.MoveMenuCursor(DirDown)
		case tcell.KeyEnter:
			titleMenu.HandleEnterEvent()
		}
	case *tcell.EventResize:
		screen.Sync()
	}
}

func (titleMenu *TitleMenu) MoveMenuCursor(dir int) {
	/*switch dir {
	case DirUp:
		if titleMenu.cursorState > 0 {
			titleMenu.cursorState--
		}
	case DirDown:
		if titleMenu.cursorState < len(titleMenu.titleMenuItems) {
			titleMenu.cursorState++
		}
	}*/
}

func (titleMenu *TitleMenu) HandleEnterEvent() {
	// Change to next place based on current cursor position and current view
}

func GenerateTitleMenu() TitleMenu {
	newGamePageItem := TitleMenuItem{0, "New game", nil}
	loadGamePageItem := TitleMenuItem{1, "Load game", nil}
	exitGameItem := TitleMenuItem{2, "Exit", nil}

	mainMenu := TitleMenuPage{
		MainMenuPageOrder,
		TitleMenuHeader,
		0,
		map[int]TitleMenuItem{
			0: newGamePageItem,
			1: loadGamePageItem,
			2: exitGameItem,
		},
	}

	newGamePage := TitleMenuPage{
		NewGamePageOrder,
		TitleMenuHeader + "Select map from list:\n",
		0,
		GenerateNewGameMapList(),
	}

	tm := TitleMenu{0, 0, map[int]TitleMenuPage{MainMenuPageOrder: mainMenu, NewGamePageOrder: newGamePage}, false}
	return tm
}

func GenerateNewGameMapList() map[int]TitleMenuItem {

	titleMenuItems := make(map[int]TitleMenuItem)
	files, err := os.ReadDir("kartor/")

	if err != nil {
		return nil
	}

	for i, file := range files {
		titleMenuItems[i] = TitleMenuItem{
			i,
			file.Name(),
			file.Name(),
		}

	}

	return titleMenuItems
}
