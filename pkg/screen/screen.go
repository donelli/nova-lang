package screen

type Screen interface {
	Init() error
	Say(y int, x int, s []rune)
	SayWithModifier(y int, x int, s []rune, modifiers int64)
	Close()
	Print(str []rune)
	Inkey(seconds int)
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
