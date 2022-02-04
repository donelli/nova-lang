package interpreter

import (
	"github.com/famz/SetLocale"
	gnc "github.com/rthornton128/goncurses"
)

type Screen interface {
	Init() error
	Say(y int, x int, s []rune)
	Close()
	Print(str []rune)
	// Inkey(seconds int)
	// AddPrompt(y int, x int, value string)
	// ReadPrompt(defIndex int) (int, string)
	// SendKeyboard(str string)
	// Lastkey() int
	// ClearTypeAhead()
	// Clear()
	// ResizeWindow(rows int, cols int)
	// SaveScreen() string
	// RestoreScreen(str string)
	// DialogBox(text string, title string)
}

const (
	KeyCode_Right     = 4
	KeyCode_Up        = 5
	KeyCode_Delete    = 7
	KeyCode_Backspace = 8
	KeyCode_Down      = 24
	KeyCode_Left      = 19
	KeyCode_Enter     = 10
	KeyCode_Esc       = 27
)

var keycodesDict = map[gnc.Key]int{
	1:   KeyCode_Up,
	259: KeyCode_Up,
	258: KeyCode_Down,
	260: KeyCode_Left,
	261: KeyCode_Right,
	263: KeyCode_Backspace,
	330: KeyCode_Delete,
}

type PromptOption struct {
	y     int
	x     int
	value string
}

type ConsoleScreen struct {
	rowCount     int
	columnCount  int
	lastKey      gnc.Key
	activePromps []PromptOption
}

func NewConsoleScreen() *ConsoleScreen {
	return &ConsoleScreen{
		rowCount:    24,
		columnCount: 80,
	}
}

const (
	COLOR_BLINK int16 = 1
	COLOR_BOLD  int16 = 2
)

func (c *ConsoleScreen) Init() error {

	SetLocale.SetLocale(SetLocale.LC_ALL, "")

	_, err := gnc.Init()

	if err != nil {
		return err
	}

	gnc.Echo(false)
	gnc.Raw(true)
	gnc.StartColor()

	gnc.InitPair(COLOR_BLINK, int16(gnc.C_YELLOW), int16(gnc.C_RED))
	gnc.InitPair(COLOR_BOLD, int16(gnc.C_BLUE), int16(gnc.C_WHITE))

	stdWin := gnc.StdScr()
	stdWin.Keypad(true)

	stdWin.Resize(c.rowCount, c.columnCount)

	return nil
}

func (c *ConsoleScreen) Close() {
	gnc.End()
}

func (c *ConsoleScreen) Say(y int, x int, s []rune) {

	for i := range s {

		charCode := int(s[i])

		printCharAt(y, x+i, charCode)

	}

	gnc.StdScr().Move(y, x+len(s))

	gnc.StdScr().Refresh()

}

func (c *ConsoleScreen) Print(s []rune) {

	y, x := gnc.StdScr().CursorYX()

	c.Say(y, x, s)

}

func printCharAt(y int, x int, charCode int) {

	// 1000000000000000000000 (BOLD)
	//  100000000000000000000 (DIM)
	//   10000000000000000000 (BLINK)
	//    1000000000000000000 (REVERSE)
	//     100000000000000000 (UNDERLINE)

	// TODO change this to the following rule:
	// Last 8 bits are the value, first 8 bits are the attributes

	stdWin := gnc.StdScr()

	if charCode >= int(gnc.A_BOLD) {
		stdWin.AttrOn(gnc.ColorPair(COLOR_BOLD))
		charCode = charCode - int(gnc.A_BOLD)
		defer stdWin.AttrOff(gnc.ColorPair(COLOR_BOLD))
	}

	if charCode >= int(gnc.A_DIM) {
		stdWin.AttrOn(gnc.Char(gnc.A_DIM))
		charCode = charCode - int(gnc.A_DIM)
		defer stdWin.AttrOff(gnc.Char(gnc.A_DIM))
	}

	if charCode >= int(gnc.A_BLINK) {
		stdWin.AttrOn(gnc.ColorPair(COLOR_BLINK))
		charCode = charCode - int(gnc.A_BLINK)
		defer stdWin.AttrOff(gnc.ColorPair(COLOR_BLINK))
	}

	if charCode >= int(gnc.A_REVERSE) {
		stdWin.AttrOn(gnc.Char(gnc.A_REVERSE))
		charCode = charCode - int(gnc.A_REVERSE)
		defer stdWin.AttrOff(gnc.Char(gnc.A_REVERSE))
	}

	if charCode >= int(gnc.A_UNDERLINE) {
		stdWin.AttrOn(gnc.Char(gnc.A_UNDERLINE))
		charCode = charCode - int(gnc.A_UNDERLINE)
		defer stdWin.AttrOff(gnc.Char(gnc.A_UNDERLINE))
	}

	// if gncChar, exits := keyboardRuneToChar[charCode]; exits {
	// 	stdWin.MoveAddChar(y, x, gncChar)
	// }

	stdWin.MoveAddChar(y, x, gnc.Char(charCode))

}
