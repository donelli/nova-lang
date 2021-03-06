// Code generated by "stringer -type=CommandType -trimprefix=CommandType_"; DO NOT EDIT.

package parser

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CommandType_Close-1]
	_ = x[CommandType_Clear-2]
	_ = x[CommandType_Exit-3]
	_ = x[CommandType_Loop-4]
	_ = x[CommandType_Dialog-5]
	_ = x[CommandType_Compile-6]
	_ = x[CommandType_Alias-7]
	_ = x[CommandType_Eject-8]
	_ = x[CommandType_Sleep-9]
	_ = x[CommandType_Store-10]
	_ = x[CommandType_Release-11]
	_ = x[CommandType_Browse-12]
	_ = x[CommandType_Count-13]
	_ = x[CommandType_Do-14]
	_ = x[CommandType_Erase-15]
	_ = x[CommandType_Assert-16]
	_ = x[CommandType_Say-17]
	_ = x[CommandType_Get-18]
	_ = x[CommandType_Read-19]
}

const _CommandType_name = "CloseClearExitLoopDialogCompileAliasEjectSleepStoreReleaseBrowseCountDoEraseAssertSayGetRead"

var _CommandType_index = [...]uint8{0, 5, 10, 14, 18, 24, 31, 36, 41, 46, 51, 58, 64, 69, 71, 76, 82, 85, 88, 92}

func (i CommandType) String() string {
	i -= 1
	if i < 0 || i >= CommandType(len(_CommandType_index)-1) {
		return "CommandType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _CommandType_name[_CommandType_index[i]:_CommandType_index[i+1]]
}
