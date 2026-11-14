package main

import (
	"Nand2Tetris/class"
	"Nand2Tetris/enum"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	arg := os.Args[1]

	if arg == "" {
		panic("script requires a filename as an argument")
	}

	inputFile, err := os.Open(arg)

	if err != nil {
		panic(err)
	}

	defer inputFile.Close()

	outputFile, err := os.Create(strings.TrimSuffix(inputFile.Name(), filepath.Ext(inputFile.Name())) + ".asm")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	parser := class.NewParser(inputFile)

	codeWriter := class.NewCodeWriter(outputFile)

	for {
		hasCommand, err := parser.AdvanceNextCommand()
		if err != nil {
			panic(err)
		}
		if !hasCommand {
			break
		}

		commandType, err := parser.CommandType()
		if err != nil {
			panic(err)
		}

		if commandType == enum.C_POP || commandType == enum.C_PUSH {
			arg1, err := parser.Arg1()
			if err != nil {
				panic(err)
			}
			arg2, err := parser.Arg2()
			if err != nil {
				panic(err)
			}

			err = codeWriter.WritePushPop(commandType, arg1, arg2)
			if err != nil {
				panic(err)
			}
		} else if commandType == enum.C_ARITHMETIC {
			arg1, err := parser.Arg1()
			if err != nil {
				panic(err)
			}

			err = codeWriter.WriteArithmetic(arg1)
			if err != nil {
				panic(err)
			}
		}
	}

	err = codeWriter.Close()
	if err != nil {
		panic(err)
	}
}
