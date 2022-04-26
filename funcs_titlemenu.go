package main

import (
	"os"
	"strings"
	"time"

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
	widthScreen, _ := screen.Size()
	currentAnimation := pageContent[currentPage.animationState]

	for x := 0; x < len(currentAnimation); x++ {
		centerX := (widthScreen / 2) - (len(currentAnimation) / 2)
		screen.SetContent(x+centerX, 0, rune(currentAnimation[x]), nil, tcell.StyleDefault) // coords to center: (x + centerX, 0)
	}

	for i := 0; i < len(pageItems); i++ { // This ensures order
		pageItem := pageItems[i]
		centerX := (widthScreen / 2) - (len(pageItem.text))
		if i == currentPage.cursorState {
			screen.SetContent(centerX-1, currentY, '>', nil, tcell.StyleDefault)
		}
		for i, c := range pageItem.text {
			screen.SetContent(i+centerX, currentY, c, nil, tcell.StyleDefault)
		}
		currentY++
	}
}

func (titleMenu *TitleMenu) AnimationHandler() {
	for {
		if titleMenu.exit {
			return
		}
		time.Sleep(100 * time.Millisecond)
		currentPage, found := titleMenu.titleMenuPages[titleMenu.pageState]
		if !found {
			currentPage = titleMenu.titleMenuPages[MainMenuPageOrder]
		}

		if currentPage.animationState >= len(currentPage.content)-1 {
			currentPage.animationState = 0
		} else {
			currentPage.animationState++ // Increment animation state after drawing it
		}
	}

}

func (titleMenu *TitleMenu) Update(screen tcell.Screen) {

}

func (titleMenu *TitleMenu) InputHandler(screen tcell.Screen) {
	for {
		if titleMenu.exit {
			return
		}
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				screen.Clear()
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
	// Change the current page cursor state
	cursorState := &titleMenu.titleMenuPages[titleMenu.pageState].cursorState
	numMenuItems := len(titleMenu.titleMenuPages[titleMenu.pageState].titleMenuItems)
	switch dir {
	case DirUp:
		if *cursorState > 0 {
			*cursorState--
		} else {
			*cursorState = numMenuItems - 1
		}
	case DirDown:
		if *cursorState < numMenuItems-1 {
			*cursorState++
		} else {
			*cursorState = 0
		}
	}
}

func (titleMenu *TitleMenu) HandleEnterEvent() {
	// Change to next place based on current cursor position and current view
	pageCursorState := titleMenu.titleMenuPages[titleMenu.pageState].cursorState
	pageItems := titleMenu.titleMenuPages[titleMenu.pageState].titleMenuItems
	switch pageItems[pageCursorState].text {
	case "Exit":
		os.Exit(0)
	case "New game":
		titleMenu.pageState = NewGamePageOrder
	case "Go back":
		titleMenu.pageState = MainMenuPageOrder
	default:
		if pageItems[pageCursorState].value != nil {
			titleMenu.selectedMap = pageItems[pageCursorState].value.(string)
			titleMenu.exit = true
		}
	}
}

func GenerateTitleMenu() TitleMenu {
	newGamePageItem := TitleMenuItem{0, "New game", nil}
	loadGamePageItem := TitleMenuItem{1, "Load game", nil}
	exitGameItem := TitleMenuItem{2, "Exit", nil}

	titleHeaderAnimation := []string{TitleMenuHeaderAnim1, TitleMenuHeaderAnim2, TitleMenuHeaderAnim3, TitleMenuHeaderAnim4, TitleMenuHeaderAnim5, TitleMenuHeaderAnim6,
		TitleMenuHeaderAnim7, TitleMenuHeaderAnim8, TitleMenuHeaderAnim9, TitleMenuHeaderAnim10, TitleMenuHeaderAnim11, TitleMenuHeaderAnim12, TitleMenuHeaderAnim13}
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

	tm := TitleMenu{0, MainMenuPageOrder, map[int]*TitleMenuPage{MainMenuPageOrder: &mainMenu, NewGamePageOrder: &newGamePage}, "", false}
	return tm
}

func GenerateNewGameMapList() map[int]TitleMenuItem {

	titleMenuItems := make(map[int]TitleMenuItem)
	files, err := os.ReadDir("kartor/")

	if err != nil {
		return nil
	}

	maxI := 0

	for i, file := range files {
		titleMenuItems[i] = TitleMenuItem{
			i,
			file.Name(),
			file.Name(),
		}
		maxI++
	}

	titleMenuItems[maxI] = TitleMenuItem{
		maxI,
		"Go back",
		nil,
	}

	return titleMenuItems
}
