package main

import (
	"os"
	"strings"

	"github.com/gdamore/tcell"
)

func (titleMenu *TitleMenu) Draw(screen tcell.Screen) {
	screen.Clear()
	titleMenu.DrawPage(screen, titleMenu.pageState)
	screen.Show()
}

func (titleMenu *TitleMenu) DrawPage(screen tcell.Screen, pageState int) {
	currentPage, found := titleMenu.titleMenuPages[pageState]
	if !found {
		currentPage = titleMenu.titleMenuPages[MainMenuPageOrder]
	}
	pageContent := currentPage.content
	pageItems := currentPage.titleMenuItems

	currentY := strings.Count(pageContent[0], "\n")
	currentAnimation := pageContent[currentPage.animationState]

	for x := 0; x < len(currentAnimation); x++ {
		screen.SetContent(x, 0, rune(currentAnimation[x]), nil, tcell.StyleDefault)
	}

	if currentPage.animationState >= len(currentPage.content)-1 {
		currentPage.animationState = 0
	} else {
		currentPage.animationState++ // Increment animation state after drawing it
	}

	startItemPositionX := 5
	for i := 0; i < len(pageItems); i++ { // This ensures order
		pageItem := pageItems[i]

		if i == titleMenu.cursorState {
			screen.SetContent(startItemPositionX-1, currentY, '>', nil, tcell.StyleDefault)
		}
		for i, c := range pageItem.text {
			posX := startItemPositionX + i
			screen.SetContent(posX, currentY, c, nil, tcell.StyleDefault)
		}
		currentY++
	}
}

func (titleMenu *TitleMenu) Update(screen tcell.Screen) {

}

func (titleMenu *TitleMenu) InputHandler(screen tcell.Screen) {
	for {
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
}

func (titleMenu *TitleMenu) MoveMenuCursor(dir int) {
	switch dir {
	case DirUp:
		if titleMenu.cursorState > 0 {
			titleMenu.cursorState--
		}
	case DirDown:
		if titleMenu.cursorState < 3 {
			titleMenu.cursorState++
		}
	}
}

func (titleMenu *TitleMenu) HandleEnterEvent() {
	// Change to next place based on current cursor position and current view
}

func GenerateTitleMenu() TitleMenu {
	newGamePageItem := TitleMenuItem{0, "New game", nil}
	loadGamePageItem := TitleMenuItem{1, "Load game", nil}
	exitGameItem := TitleMenuItem{2, "Exit", nil}

	titleHeaderAnimation := []string{TitleMenuHeaderAnim1, TitleMenuHeaderAnim2, TitleMenuHeaderAnim3, TitleMenuHeaderAnim4, TitleMenuHeaderAnim5, TitleMenuHeaderAnim6}

	mainMenu := TitleMenuPage{
		MainMenuPageOrder,
		titleHeaderAnimation,
		0,
		0,
		map[int]TitleMenuItem{
			0: newGamePageItem,
			1: loadGamePageItem,
			2: exitGameItem,
		},
	}

	newGamePage := TitleMenuPage{
		NewGamePageOrder,
		titleHeaderAnimation,
		0,
		0,
		GenerateNewGameMapList(),
	}

	tm := TitleMenu{0, 0, map[int]*TitleMenuPage{MainMenuPageOrder: &mainMenu, NewGamePageOrder: &newGamePage}, false}
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
