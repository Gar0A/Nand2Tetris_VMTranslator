package class

import (
	"Nand2Tetris/enum"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CodeWriters Interface for the VM code writer
type CodeWriters interface {
	WriteArithmetic() error
	WritePushPop(string) error
}

// CodeWriter Struct for the VM code writer
type CodeWriter struct {
	writer   *bufio.Writer
	fileName string
}

// NewCodeWriter Constructor for the parser struct
func NewCodeWriter(f *os.File) CodeWriter {
	writer := bufio.NewWriter(f)
	return CodeWriter{writer: writer, fileName: strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))}
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
			pushHeader = fmt.Sprintf("\n@LCL\nD=A\n@%d\nA=D+A\nD=M", index)
		case "argument":
			pushHeader = fmt.Sprintf("\n@ARG\nD=A\n@%d\nA=D+A\nD=M", index)
		case "this":
			pushHeader = fmt.Sprintf("\n@THIS\nD=A\n@%d\nA=D+A\nD=M", index)
		case "that":
			pushHeader = fmt.Sprintf("\n@THAT\nD=A\n@%i\nA=D+A\nD=M", index)
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
			popHeader = fmt.Sprintf("\n@LCL\nD=A\n@%d\nA=D+A\nD=A\n@R13\nM=D\n@SP\nA=M\nD=M\n@R13\nA=M", index)
		case "argument":
			popHeader = fmt.Sprintf("\n@ARG\nD=A\n@%d\nA=D+A\nD=A\n@R13\nM=D\n@SP\nA=M\nD=M\n@R13\nA=M", index)
		case "this":
			popHeader = fmt.Sprintf("\n@THIS\nD=A\n@%d\nA=D+A\nD=A\n@R13\nM=D\n@SP\nA=M\nD=M\n@R13\nA=M", index)
		case "that":
			popHeader = fmt.Sprintf("\n@THAT\nD=A\n@%d\nA=D+A\nD=A\n@R13\nM=D\n@SP\nA=M\nD=M\n@R13\nA=M", index)
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

	fmt.Print(code.String())

	_, err := cw.writer.WriteString(code.String())
	if err != nil {
		return err
	}

	return nil
}

func (cw *CodeWriter) Close() error {
	return cw.writer.Flush()
}
