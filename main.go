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
}
