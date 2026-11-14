package enum

// CommandType Enum type for the 9 possible commands of the VM code
type CommandType int8

const (
	C_UNKNOWN CommandType = iota
	C_ARITHMETIC
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)
