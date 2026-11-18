package class

import (
	"Nand2Tetris/enum"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CodeWriters Interface for the VM code writer
type CodeWriters interface {
	SetFileName(string)
	WriteInit() error
	WriteArithmetic() error
	WritePushPop(string) error
	WriteLabel(string) error
	WriteGoto(string) error
	WriteIf(string) error
	WriteFunction(string, int) error
	WriteCall(string, int) error
	WriteReturn(string) error
}

// CodeWriter Struct for the VM code writer
type CodeWriter struct {
	writer          *bufio.Writer
	fileName        string
	arithmeticCount map[string]int
}

// NewCodeWriter Constructor for the parser struct
func NewCodeWriter(f *os.File) CodeWriter {
	writer := bufio.NewWriter(f)
	arithmeticCount := map[string]int{
		"eq": 0,
		"gt": 0,
		"lt": 0,
	}
	return CodeWriter{writer: writer, arithmeticCount: arithmeticCount}
}

func (cw *CodeWriter) SetFileName(name string) {
	cw.fileName = name
}

func (cw *CodeWriter) WriteInit() error {
	_, err := cw.writer.WriteString("@SP\nM=256")
	if err != nil {
		return err
	}

	// TODO : add WriteCall("call Sys.init") once the implementation of the function is done

	return nil
}

func (cw *CodeWriter) WritePushPop(commandType enum.CommandType, segment string, index int) error {
	command := ""
	if commandType == enum.C_POP {
		command = "pop"
	} else if commandType == enum.C_PUSH {
		command = "push"
	} else {
		return fmt.Errorf("forbidden command %d", commandType)
	}

	var code strings.Builder
	code.WriteString(fmt.Sprintf("\n// %s %s %d", command, segment, index))

	if commandType == enum.C_PUSH {

		pushHeader := ""

		switch segment {
		case "local":
			pushHeader = fmt.Sprintf("\n@LCL\nD=M\n@%d\nA=D+A\nD=M", index)
		case "argument":
			pushHeader = fmt.Sprintf("\n@ARG\nD=M\n@%d\nA=D+A\nD=M", index)
		case "this":
			pushHeader = fmt.Sprintf("\n@THIS\nD=M\n@%d\nA=D+A\nD=M", index)
		case "that":
			pushHeader = fmt.Sprintf("\n@THAT\nD=M\n@%d\nA=D+A\nD=M", index)
		case "constant":
			pushHeader = fmt.Sprintf("\n@%d\nD=A", index)
		case "static":
			pushHeader = fmt.Sprintf("\n@%s.%d\nD=M", cw.fileName, index)
		case "temp":
			if index <= 8 {
				pushHeader = fmt.Sprintf("\n@R5\nD=A\n@%d\nA=D+A\nD=M", index)
			}
		case "pointer":
			if index == 0 {
				pushHeader = "\n@THIS\nD=M"
			} else if index == 1 {
				pushHeader = "\n@THAT\nD=M"
			}

		}

		if pushHeader != "" {
			code.WriteString(pushHeader)
		} else {
			return fmt.Errorf("error with push command with segment %q and index %d", segment, index)
		}

		code.WriteString("\n@SP\nA=M\nM=D\n@SP\nM=M+1")

	} else if commandType == enum.C_POP {

		code.WriteString("\n@SP\nM=M-1")

		popHeader := ""

		switch segment {
		case "local":
			popHeader = fmt.Sprintf("\n@LCL\nD=M\n@%d\nA=D+A\nD=A\n@R13\nM=D\n@SP\nA=M\nD=M\n@R13\nA=M", index)
		case "argument":
			popHeader = fmt.Sprintf("\n@ARG\nD=M\n@%d\nA=D+A\nD=A\n@R13\nM=D\n@SP\nA=M\nD=M\n@R13\nA=M", index)
		case "this":
			popHeader = fmt.Sprintf("\n@THIS\nD=M\n@%d\nA=D+A\nD=A\n@R13\nM=D\n@SP\nA=M\nD=M\n@R13\nA=M", index)
		case "that":
			popHeader = fmt.Sprintf("\n@THAT\nD=M\n@%d\nA=D+A\nD=A\n@R13\nM=D\n@SP\nA=M\nD=M\n@R13\nA=M", index)
		case "static":
			popHeader = fmt.Sprintf("\n@SP\nA=M\nD=M\n@%s.%d", cw.fileName, index)
		case "temp":
			if index <= 8 {
				popHeader = fmt.Sprintf("\n@R5\nD=A\n@%d\nA=D+A\nD=A\n@R13\nM=D\n@SP\nA=M\nD=M\n@R13\nA=M", index)
			}
		case "pointer":
			if index == 0 {
				popHeader = "\n@SP\nA=M\nD=M\n@THIS"
			} else if index == 1 {
				popHeader = "\n@SP\nA=M\nD=M\n@THAT"
			}

		}

		if popHeader != "" {
			code.WriteString(popHeader)
		} else {
			return fmt.Errorf("error with pop command with segment %q and index %d", segment, index)
		}

		code.WriteString("\nM=D")

	} else {
		return fmt.Errorf("forbidden command %d", commandType)
	}

	_, err := cw.writer.WriteString(code.String())
	if err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) WriteArithmetic(command string) error {
	var code strings.Builder
	code.WriteString(fmt.Sprintf("\n// %s", command))

	arithmeticCode := ""
	switch command {
	case "add":
		arithmeticCode = "\n@R6\nD=M\n@R5\nM=D+M"
	case "sub":
		arithmeticCode = "\n@R6\nD=M\n@R5\nM=D-M"
	case "neg":
		arithmeticCode = "\n@R5\nM=-M"
	case "eq":
		count := cw.arithmeticCount[command]
		arithmeticCode = fmt.Sprintf("\n@R6\nD=M\n@R5\nD=D-M\n@EQ_%d\nD;JEQ\nD=0\n@EQ_%d_END\nA;JMP\n(EQ_%d)\nD=-1\n(EQ_%d_END)\n@R5\nM=D", count, count, count, count)
		cw.arithmeticCount[command]++
	case "gt":
		count := cw.arithmeticCount[command]
		arithmeticCode = fmt.Sprintf("\n@R6\nD=M\n@R5\nD=D-M\n@GT_%d\nD;JGT\nD=0\n@GT_%d_END\nA;JMP\n(GT_%d)\nD=-1\n(GT_%d_END)\n@R5\nM=D", count, count, count, count)
		cw.arithmeticCount[command]++
	case "lt":
		count := cw.arithmeticCount[command]
		arithmeticCode = fmt.Sprintf("\n@R6\nD=M\n@R5\nD=D-M\n@LT_%d\nD;JLT\nD=0\n@LT_%d_END\nA;JMP\n(LT_%d)\nD=-1\n(LT_%d_END)\n@R5\nM=D", count, count, count, count)
		cw.arithmeticCount[command]++
	case "and":
		arithmeticCode = "\n@R6\nD=M\n@R5\nM=D&M"
	case "or":
		arithmeticCode = "\n@R6\nD=M\n@R5\nM=D|M"
	case "not":
		arithmeticCode = "\n@R5\nM=!M"
	}

	if arithmeticCode != "" {
		code.WriteString(arithmeticCode)
	} else {
		return fmt.Errorf("error with %s command", command)
	}

	if command == "neg" || command == "not" {
		err := cw.WritePushPop(enum.C_POP, "temp", 0)
		if err != nil {
			return err
		}
	} else {
		err := cw.WritePushPop(enum.C_POP, "temp", 0)
		if err != nil {
			return err
		}

		err = cw.WritePushPop(enum.C_POP, "temp", 1)
		if err != nil {
			return err
		}
	}

	_, err := cw.writer.WriteString(code.String())
	if err != nil {
		return err
	}

	err = cw.WritePushPop(enum.C_PUSH, "temp", 0)
	if err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) WriteLabel(label string) error {
	_, err := cw.writer.WriteString(fmt.Sprintf("\n(%s)", label))
	if err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) WriteGoto(label string) error {
	_, err := cw.writer.WriteString(fmt.Sprintf("\n@%s\n0;JMP", label))
	if err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) WriteIf(label string) error {
	err := cw.WritePushPop(enum.C_POP, "temp", 0)
	if err != nil {
		return err
	}

	_, err = cw.writer.WriteString(fmt.Sprintf("\n@R5\nD=M\n@%s\nD;JNE", label))
	if err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) Close() error {
	return cw.writer.Flush()
}
