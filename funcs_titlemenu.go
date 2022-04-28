package main

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

func (titleMenu *TitleMenu) Draw(screen tcell.Screen, m *sync.Mutex) {
	m.Lock()
	screen.Clear()
	titleMenu.DrawPage(screen, titleMenu.pageState)
	screen.Show()
	m.Unlock()
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
		pageItem := (pageItems)[i]
		centerX := (widthScreen / 2) - (len(pageItem.text))
		if i == currentPage.cursorState {
			screen.SetContent(centerX-1, currentY, '>', nil, tcell.StyleDefault)
		}
		for i, c := range pageItem.text {
			screen.SetContent(i+centerX, currentY, c, nil, tcell.StyleDefault)
		}
		if strings.Contains(pageItem.text, "Height") || strings.Contains(pageItem.text, "Width") { // Also draw the values for height and width if they exist
			for i, c := range pageItem.value.(string) {
				screen.SetContent(i+len(pageItem.text)+centerX, currentY, c, nil, tcell.StyleDefault)
			}
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

func (titleMenu *TitleMenu) InputHandler(screen tcell.Screen, m *sync.Mutex) {
	for {
		if titleMenu.exit {
			return
		}
		ev := screen.PollEvent()
		m.Lock()
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
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				titleMenu.HandleBackspaceEvent()
			}
			switch ev.Rune() {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				titleMenu.HandleNumberInputEvent(ev.Rune())
			}
		case *tcell.EventResize:
			screen.Sync()
		}
		m.Unlock()
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
	case "Generate new map":
		titleMenu.pageState = GenerateMapPageOrder
	case "Start":
		// Convert the width and height values to integers
		titleMenu.generatedMapWidth, _ = strconv.Atoi(titleMenu.titleMenuPages[GenerateMapPageOrder].titleMenuItems[0].value.(string))
		titleMenu.generatedMapHeight, _ = strconv.Atoi(titleMenu.titleMenuPages[GenerateMapPageOrder].titleMenuItems[1].value.(string))
		titleMenu.selectedMap = "generatemap"
		titleMenu.exit = true
	default:
		if pageItems[pageCursorState].value != nil && titleMenu.pageState == NewGamePageOrder {
			titleMenu.selectedMap = pageItems[pageCursorState].value.(string)
			titleMenu.exit = true
		}
	}
}

func (titleMenu *TitleMenu) HandleBackspaceEvent() {
	pageCursorState := titleMenu.titleMenuPages[titleMenu.pageState].cursorState
	pageItems := titleMenu.titleMenuPages[titleMenu.pageState].titleMenuItems
	currentPageItem := pageItems[pageCursorState]

	if strings.Contains(currentPageItem.text, "Height") || strings.Contains(currentPageItem.text, "Width") { // Also draw the values for height and width if they exist
		if len(currentPageItem.value.(string)) > 0 {
			currentPageItem.value = currentPageItem.value.(string)[:len(currentPageItem.value.(string))-1]
		}
	}
}

func (titleMenu *TitleMenu) HandleNumberInputEvent(r rune) {
	pageCursorState := titleMenu.titleMenuPages[titleMenu.pageState].cursorState
	pageItems := titleMenu.titleMenuPages[titleMenu.pageState].titleMenuItems
	currentPageItem := pageItems[pageCursorState]

	if strings.Contains(currentPageItem.text, "Height") || strings.Contains(currentPageItem.text, "Width") { // Also draw the values for height and width if they exist
		if len(currentPageItem.value.(string)) < 3 {
			currentPageItem.value = currentPageItem.value.(string) + string(r)
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
		map[int]*TitleMenuItem{
			0: &newGamePageItem,
			1: &loadGamePageItem,
			2: &exitGameItem,
		},
	}

	newGamePage := TitleMenuPage{
		NewGamePageOrder,
		titleHeaderAnimation,
		0,
		0,
		GenerateNewGameMapList(),
	}

	generateMapPage := TitleMenuPage{
		GenerateMapPageOrder,
		titleHeaderAnimation,
		0,
		0,
		GenerateMapPageItems(),
	}

	tm := TitleMenu{0,
		MainMenuPageOrder,
		map[int]*TitleMenuPage{MainMenuPageOrder: &mainMenu, NewGamePageOrder: &newGamePage, GenerateMapPageOrder: &generateMapPage},
		"",
		DefaultGeneratedMapWidth,
		DefaultGeneratedMapHeight,
		false}
	return tm
}

func GenerateMapPageItems() map[int]*TitleMenuItem {
	return map[int]*TitleMenuItem{
		0: {
			0,
			"Width: ",
			"20",
		},
		1: {
			1,
			"Height: ",
			"20",
		},
		2: {
			2,
			"Start",
			nil,
		},
	}
}

func GenerateNewGameMapList() map[int]*TitleMenuItem {

	titleMenuItems := make(map[int]*TitleMenuItem)
	files, err := os.ReadDir("kartor/")

	if err != nil {
		return nil
	}

	maxI := 0

	for i, file := range files {
		titleMenuItems[i] = &TitleMenuItem{
			i,
			file.Name(),
			file.Name(),
		}
		maxI++
	}

	titleMenuItems[maxI] = &TitleMenuItem{ // Generate new map menu item
		maxI,
		"Generate new map",
		nil,
	}

	maxI++

	titleMenuItems[maxI] = &TitleMenuItem{ // Go back menu item
		maxI,
		"Go back",
		nil,
	}

	return titleMenuItems
}
