package class

import (
	"Nand2Tetris/enum"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Parsers Interface for the VM code parser
type Parsers interface {
	AdvanceNextCommand() (bool, error)
	CommandType() (enum.CommandType, error)
	Arg1() (string, error)
	Arg2() (int, error)
}

// Parser Struct for the VM code parser
type Parser struct {
	scanner        *bufio.Scanner
	currentCommand string
	lineCount      int
}

// NewParser Constructor for the parser struct
func NewParser(f *os.File) Parser {
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	return Parser{scanner: scanner, lineCount: 0}
}

func (p *Parser) AdvanceNextCommand() (bool, error) {
	rWhitespaceLine := regexp.MustCompile("^\\s*$|^\\s*(?://+)+")
	rArithmetic := regexp.MustCompile("^\\s*(add|sub|neg|eq|gt|lt|and|or|not)(?:\\s*$|\\s+(?://+)+)")
	rPushConstant := regexp.MustCompile("^\\s*push\\s+constant\\s+(\\d+)(?:\\s*$|\\s+(?://+)+)")
	rPushPopTemp := regexp.MustCompile("^\\s*(pop|push)\\s+temp\\s+([0-8])(?:\\s*$|\\s+(?://+)+)")
	rPushPopPointer := regexp.MustCompile("^\\s*(pop|push)\\s+pointer\\s+([0-1])(?:\\s*$|\\s+(?://+)+)")
	rPushPopRest := regexp.MustCompile("^\\s*(pop|push)\\s+(local|argument|this|that|static)\\s+(\\d+)(?:\\s*$|\\s+(?://+)+)")

	if !p.scanner.Scan() {
		return false, nil
	}
	p.lineCount++

	line := p.scanner.Text()

	line = strings.ToLower(line)

	if rWhitespaceLine.MatchString(line) {
		return p.AdvanceNextCommand()
	} else if rArithmetic.MatchString(line) {
		parts := rArithmetic.FindStringSubmatch(line)
		p.currentCommand = parts[1]
		return true, nil
	} else if rPushConstant.MatchString(line) {
		parts := rPushConstant.FindStringSubmatch(line)
		n, _ := strconv.Atoi(parts[1])
		p.currentCommand = fmt.Sprintf("push constant %d", n)
		return true, nil
	} else if rPushPopTemp.MatchString(line) {
		parts := rPushPopTemp.FindStringSubmatch(line)
		p.currentCommand = fmt.Sprintf("%s temp %s", parts[1], parts[2])
		return true, nil
	} else if rPushPopPointer.MatchString(line) {
		parts := rPushPopPointer.FindStringSubmatch(line)
		p.currentCommand = fmt.Sprintf("%s pointer %s", parts[1], parts[2])
		return true, nil
	} else if rPushPopRest.MatchString(line) {
		parts := rPushPopRest.FindStringSubmatch(line)
		n, _ := strconv.Atoi(parts[3])
		p.currentCommand = fmt.Sprintf("%s %s %d", parts[1], parts[2], n)
		return true, nil
	}
	return false, fmt.Errorf("error in line %d: unknown command : %q", p.lineCount, line)
}

func (p *Parser) CommandType() (enum.CommandType, error) {
	rArithmetic := regexp.MustCompile("^(?:add|sub|neg|eq|gt|lt|and|or|not)")
	rPushPop := regexp.MustCompile("^(pop|push)")

	command := p.currentCommand

	if rArithmetic.MatchString(command) {
		return enum.C_ARITHMETIC, nil
	} else if rPushPop.MatchString(command) {
		parts := rPushPop.FindStringSubmatch(command)
		if parts[1] == "pop" {
			return enum.C_POP, nil
		} else if parts[1] == "push" {
			return enum.C_PUSH, nil
		}
	}

	return enum.C_UNKNOWN, fmt.Errorf("error in line %d: unknown command type : %q", p.lineCount, command)
}

func (p *Parser) Arg1() (string, error) {
	command := p.currentCommand
	commandType, err := p.CommandType()

	if err != nil {
		return "", err
	}

	if commandType == enum.C_ARITHMETIC {
		return command, nil
	} else if commandType != enum.C_RETURN && commandType != enum.C_UNKNOWN {
		arg1 := strings.Fields(command)[1]
		return arg1, nil
	}

	return "", fmt.Errorf("error in line %d: command does contain arguments : %q", p.lineCount, command)
}

func (p *Parser) Arg2() (int, error) {
	command := p.currentCommand
	commandType, err := p.CommandType()

	if err != nil {
		return -1, err
	}

	if commandType == enum.C_POP || commandType == enum.C_PUSH || commandType == enum.C_FUNCTION || commandType == enum.C_CALL {
		arg2 := strings.Fields(command)[2]
		n, _ := strconv.Atoi(arg2)
		return n, nil
	}

	return -1, fmt.Errorf("error in line %d: command does not contain a second argument : %q", p.lineCount, command)
}
