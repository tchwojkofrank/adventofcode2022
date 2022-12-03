package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func readInput(fname string) string {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string
	return string(content)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Enter input file name.\n")
		return
	}
	params := os.Args[1]
	inputName := strings.Split(params, " ")[0]
	text := readInput(inputName)
	run(text)
}

func run(input string) {
}
