package screen

import (
	"fmt"
	"os"

	"github.com/famz/SetLocale"
	gnc "github.com/rthornton128/goncurses"
)

type OutputType int

const (
	OutputType_Console OutputType = iota + 1
	OutputType_Test
)

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

const (
	COLOR_BLINK int16 = 1
	COLOR_BOLD  int16 = 2
)
const (
	Modif_Bold      = int64(gnc.A_BOLD)
	Modif_Dim       = int64(gnc.A_DIM)
	Modif_Blink     = int64(gnc.A_BLINK)
	Modif_Reverse   = int64(gnc.A_REVERSE)
	Modif_Underline = int64(gnc.A_UNDERLINE)
)

type PromptOption struct {
	y     int
	x     int
	value string
}

type ConsoleScreen struct {
	rowCount       int
	columnCount    int
	lastKey        gnc.Key
	activePromps   []PromptOption
	outputType     OutputType
	testOutputFile *os.File
}

func NewConsoleScreen(outputType OutputType) *ConsoleScreen {
	return &ConsoleScreen{
		rowCount:     24,
		columnCount:  80,
		lastKey:      0,
		outputType:   outputType,
		activePromps: []PromptOption{},
	}
}

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
	stdWin.ScrollOk(true)

	stdWin.Resize(c.rowCount, c.columnCount)

	if c.outputType == OutputType_Test {

		c.testOutputFile, err = os.Create("test.txt")

		if err != nil {
			return err
		}

	}

	return nil
}

func (c *ConsoleScreen) Close() {

	gnc.End()

	if c.outputType == OutputType_Test {
		c.testOutputFile.Close()
	}

}

func modifyStringRune(str []rune, modifer int64) []rune {

	newStr := make([]rune, len(str))

	for i, b := range str {
		newStr[i] = b | rune(modifer)
	}

	return newStr
}

func (c *ConsoleScreen) writeToTestFile(str string) {

	if c.outputType != OutputType_Test {
		return
	}

	c.testOutputFile.WriteString(str)

}

func (c *ConsoleScreen) SayWithModifier(y int, x int, s []rune, modifiers int64) {

	modiferStr := modifyStringRune(s, modifiers)

	c.Say(y, x, modiferStr)

}

func (c *ConsoleScreen) Say(y int, x int, s []rune) {

	c.writeToTestFile(fmt.Sprintf("\nSAY %d %d: ", y, x))

	for i := range s {

		charCode := int(s[i])

		char := printCharAt(y, x+i, charCode)

		c.writeToTestFile(string(char))

	}

	gnc.StdScr().Move(y, x+len(s))

	gnc.StdScr().Refresh()

}

func (c *ConsoleScreen) Print(s []rune) {

	y, _ := gnc.StdScr().CursorYX()

	if y >= c.rowCount-1 {
		gnc.StdScr().Scroll(1)
		gnc.StdScr().Refresh()
		y = c.rowCount - 1
	}

	c.Say(y, 0, s)

	gnc.StdScr().Move(y+1, 0)
	gnc.StdScr().Refresh()

}

func getChar() gnc.Key {

	char := gnc.StdScr().GetChar()

	keycode, found := keycodesDict[char]
	if found {
		return gnc.Key(keycode)
	}

	if char >= 0 && char <= 126 {
		return char
	}

	return 0
}

func (c *ConsoleScreen) Inkey(seconds int) {

	stdWin := gnc.StdScr()

	if seconds > 0 {
		stdWin.Timeout(seconds * 1000)
	}

	c.lastKey = getChar()

	if seconds > 0 {
		stdWin.Timeout(-1)
	}

}

func printCharAt(y int, x int, charCode int) int {

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

	stdWin.MoveAddChar(y, x, intToGncChar(charCode))

	return charCode
}

func intToGncChar(charCode int) gnc.Char {
	return gnc.Char(charCode)
}
