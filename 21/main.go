package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
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
	if len(args) < 2 {
		fmt.Printf("Enter input file name.\n")
		return
	}
	params := os.Args[1]
	inputName := strings.Split(params, " ")[0]
	start := time.Now()
	text := readInput(inputName)
	run(text)
	end := time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
}

type Monkey struct {
	name      string
	value     int
	operation string
	operands  *[2]string
}

type Monkeys map[string]Monkey

func newMonkey(info string) Monkey {
	var monkey Monkey
	instructions := strings.Split(info, " ")
	monkey.name = strings.TrimSuffix(instructions[0], ":")
	if len(instructions) == 2 {
		monkey.value, _ = strconv.Atoi(instructions[1])
	} else if len(instructions) == 4 {
		var operands [2]string
		operands[0] = instructions[1]
		operands[1] = instructions[3]
		monkey.operands = &operands
		monkey.operation = instructions[2]
		monkey.value = math.MaxInt
	}
	return monkey
}

func (monkeys Monkeys) monkeyNumber(monkey Monkey) int {
	if monkey.operands == nil {
		return monkey.value
	} else {
		if monkey.value != math.MaxInt {
			return monkey.value
		}
		value1 := monkeys.monkeyNumber(monkeys[monkey.operands[0]])
		value2 := monkeys.monkeyNumber(monkeys[monkey.operands[1]])
		switch monkey.operation {
		case "+":
			return value1 + value2
		case "-":
			return value1 - value2
		case "*":
			return value1 * value2
		case "/":
			return value1 / value2
		}
	}
	return math.MaxInt
}

func run(input string) {
	monkeyInfo := strings.Split(input, "\n")
	monkeys := Monkeys(make(map[string]Monkey))
	for _, m := range monkeyInfo {
		monkey := newMonkey(m)
		monkeys[monkey.name] = monkey
	}
	fmt.Printf("Root monkey shouts %d\n", monkeys.monkeyNumber(monkeys["root"]))
}
