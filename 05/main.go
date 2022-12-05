package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readInput(fname string) string {
	content, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string
	return string(content)
}

func main() {
	args := os.Args
	params := "input"
	if len(args) >= 2 {
		params = os.Args[1]
	}
	inputName := strings.Split(params, " ")[0]
	text := readInput(inputName)
	run(text)
}

func reverseSlice[T interface{}](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

type Stack []string

func initStacks(crateStrings []string) []Stack {
	stackCount := (len(crateStrings[0]) + 3) / 4
	stacks := make([]Stack, stackCount)
	crateStrings = crateStrings[1:]
	for _, crateString := range crateStrings {
		for sid, j := 0, 1; j < len(crateString); sid, j = sid+1, j+4 {
			if crateString[j:j+1] != " " {
				stacks[sid] = append(stacks[sid], crateString[j:j+1])
			}
		}
	}
	return stacks
}

func moveCrate9000(stacks []Stack, count int, from int, to int) []Stack {
	for i := 0; i < count; i++ {
		fromstacksize := len(stacks[from])
		if fromstacksize > 0 {
			box := stacks[from][fromstacksize-1]
			stacks[from] = stacks[from][0 : fromstacksize-1]
			stacks[to] = append(stacks[to], box)
		}
	}

	return stacks
}

func moveCrate9001(stacks []Stack, count int, from int, to int) []Stack {
	fromstacksize := len(stacks[from])
	stacks[to] = append(stacks[to], stacks[from][fromstacksize-count:]...)
	stacks[from] = stacks[from][0 : fromstacksize-count]

	return stacks
}

func doInstruction(instruction string, stacks []Stack, moveCrate func([]Stack, int, int, int) []Stack) []Stack {
	parts := strings.Split(instruction, " ")
	countString := parts[1]
	fromString := parts[3]
	toString := parts[5]
	count, _ := strconv.Atoi(countString)
	from, _ := strconv.Atoi(fromString)
	to, _ := strconv.Atoi(toString)
	from = from - 1
	to = to - 1
	stacks = moveCrate(stacks, count, from, to)
	return stacks
}

func doAllInstructions(instructions []string, stacks []Stack, moveCrate func([]Stack, int, int, int) []Stack) []Stack {
	for _, i := range instructions {
		stacks = doInstruction(i, stacks, moveCrate)
		fmt.Println(i)
		fmt.Println(stacks)
	}
	return stacks
}

func printTop(stacks []Stack) {
	for _, stack := range stacks {
		fmt.Printf("%v", stack[len(stack)-1])
	}
	fmt.Println()
}

func run(input string) {
	inputParts := strings.Split(input, "\n\n")
	crateStackStrings := strings.Split(inputParts[0], "\n")
	instructionsString := inputParts[1]
	instructions := strings.Split(instructionsString, "\n")
	crateStackStrings = reverseSlice(crateStackStrings)
	stacks := initStacks(crateStackStrings)
	stacks = doAllInstructions(instructions, stacks, moveCrate9000)
	printTop(stacks)

	//part 2
	inputParts = strings.Split(input, "\n\n")
	crateStackStrings = strings.Split(inputParts[0], "\n")
	instructionsString = inputParts[1]
	instructions = strings.Split(instructionsString, "\n")
	crateStackStrings = reverseSlice(crateStackStrings)
	stacks = initStacks(crateStackStrings)
	fmt.Println(stacks)
	stacks = doAllInstructions(instructions, stacks, moveCrate9001)
	fmt.Println(stacks)
	printTop(stacks)

}
