package shared

const (
	Digits              = "0123456789"
	Letters             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LettersAndUnderline = Letters + "_"
	LettersAndDigits    = Letters + Digits
	DigitsAndDot        = Digits + "."
	ValidPathChars      = LettersAndDigits + "/_."
	WhitespaceChars     = " \t "
)
