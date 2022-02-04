package interpreter

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type Screen interface {
	Init() error
	WriteAtPos(x int, y int, str []rune)
	Close()
	Print(str []rune)
}

// Console screen

type ConsoleScreen struct {
	screen         tcell.Screen
	currentCursorX int
	currentCursorY int
	defaultStyle   tcell.Style
	blinkStyle     tcell.Style
	boldStyle      tcell.Style
}

func NewConsoleScreen() *ConsoleScreen {
	return &ConsoleScreen{
		screen:         nil,
		currentCursorX: 0,
		currentCursorY: 0,
		defaultStyle: tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorBlack),
		blinkStyle: tcell.StyleDefault.
			Foreground(tcell.ColorYellow).
			Background(tcell.ColorRed),
		boldStyle: tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDarkBlue),
	}
}

func (cs *ConsoleScreen) Init() error {
	if s, e := tcell.NewScreen(); e != nil {
		return e
	} else if e = s.Init(); e != nil {
		return e
	} else {
		cs.screen = s
		cs.screen.ShowCursor(0, 0)
		return nil
	}
}

func (cs *ConsoleScreen) updateCursorPos(x int, y int) {

	cs.currentCursorX = x
	cs.currentCursorY = y

	cs.screen.ShowCursor(x, y)

}

func (cs *ConsoleScreen) Print(str []rune) {
	cs.WriteAtPos(cs.currentCursorX, cs.currentCursorY, str)
	cs.updateCursorPos(0, cs.currentCursorY+1)
}

func (cs *ConsoleScreen) WriteAtPos(x int, y int, str []rune) {

	style := cs.defaultStyle

	for i := 0; i < len(str); i++ {

		char := str[i]

		// 0x09 indicates the end of the formating of a string
		if char == 0x09 {
			style = cs.defaultStyle
			continue
		}

		// 0x08 indicates the start of the formating of a string
		if char == 0x08 {

			i++

			char := str[i]

			if char == 'A' { // Default (reversed)
				style = cs.defaultStyle.Reverse(true)
			} else if char == 'B' { // Blink
				style = cs.blinkStyle
			} else if char == 'C' { // Blink (reversed)
				style = cs.blinkStyle.Reverse(true)
			} else if char == 'D' { // Bold
				style = cs.boldStyle
			} else if char == 'E' { // Bold (reversed)
				style = cs.boldStyle.Reverse(true)
			}

			continue
		}

		cs.screen.SetContent(x, y, rune(char), nil, style)
		x++

	}

	cs.updateCursorPos(x, y)

	cs.screen.Show()
}

func (cs *ConsoleScreen) Close() {
	cs.screen.Fini()
}

// Test screen

type SimulationScreen struct {
	screen         tcell.Screen
	currentCursorX int
	currentCursorY int
	defaultStyle   tcell.Style
	blinkStyle     tcell.Style
	boldStyle      tcell.Style
}

func NewSimulationScreen() *SimulationScreen {
	return &SimulationScreen{
		screen:         nil,
		currentCursorX: 0,
		currentCursorY: 0,
		defaultStyle: tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorBlack),
		blinkStyle: tcell.StyleDefault.
			Foreground(tcell.ColorYellow).
			Background(tcell.ColorRed),
		boldStyle: tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDarkBlue),
	}
}

func (cs *SimulationScreen) Init() error {
	cs.screen = tcell.NewSimulationScreen("")
	cs.screen.ShowCursor(0, 0)
	return nil
}

func (cs *SimulationScreen) updateCursorPos(x int, y int) {

	cs.currentCursorX = x
	cs.currentCursorY = y

	cs.screen.ShowCursor(x, y)

}

func (cs *SimulationScreen) Print(str []rune) {

	fmt.Println(string(str))

	// cs.WriteAtPos(cs.currentCursorX, cs.currentCursorY, str)
	// cs.updateCursorPos(0, cs.currentCursorY+1)
}

func (cs *SimulationScreen) WriteAtPos(x int, y int, str []rune) {

	style := cs.defaultStyle

	for i := 0; i < len(str); i++ {

		char := str[i]

		// 0x09 indicates the end of the formating of a string
		if char == 0x09 {
			style = cs.defaultStyle
			continue
		}

		// 0x08 indicates the start of the formating of a string
		if char == 0x08 {

			i++

			char := str[i]

			if char == 'A' { // Default (reversed)
				style = cs.defaultStyle.Reverse(true)
			} else if char == 'B' { // Blink
				style = cs.blinkStyle
			} else if char == 'C' { // Blink (reversed)
				style = cs.blinkStyle.Reverse(true)
			} else if char == 'D' { // Bold
				style = cs.boldStyle
			} else if char == 'E' { // Bold (reversed)
				style = cs.boldStyle.Reverse(true)
			}

			continue
		}

		cs.screen.SetContent(x, y, char, nil, style)
		x++

	}

	cs.updateCursorPos(x, y)

	cs.screen.Show()
}

func (cs *SimulationScreen) Close() {
	cs.screen.Fini()
}
