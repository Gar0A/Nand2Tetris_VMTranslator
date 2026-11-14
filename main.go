package main

import (
	"Nand2Tetris/class"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputFile, err := os.Open("test.vm")

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

	_, err = parser.AdvanceNextCommand()
	if err != nil {
		panic(err)
	}

	fmt.Println("command: ", parser.CurrentCommand)

	commandType, _ := parser.CommandType()

	fmt.Println("command type : ", commandType)

	arg1, _ := parser.Arg1()

	fmt.Println("arg1 : ", arg1)

	arg2, _ := parser.Arg2()

	fmt.Println("arg2 : ", arg2)

	err = codeWriter.WritePushPop(commandType, arg1, arg2)
	if err != nil {
		panic(err)
	}

	err = codeWriter.Close()
	if err != nil {
		panic(err)
	}
}
