package main

import (
	"Nand2Tetris/class"
	"fmt"
	"os"
)

func main() {
	myfile, err := os.Open("test.vm")

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer myfile.Close()

	parser := class.NewParser(myfile)

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
}
