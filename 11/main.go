package main

import (
	"fmt"
	"log"
	"os"
	"sort"
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
	if len(args) < 2 {
		fmt.Printf("Enter input file name.\n")
		return
	}
	params := os.Args[1]
	inputName := strings.Split(params, " ")[0]
	text := readInput(inputName)
	run(text)
}

const (
	add      = 1
	multiply = 2
)

type MonkeyInfo struct {
	worries     []int
	operation   int
	operand     string
	testValue   int
	trueMonkey  int
	falseMonkey int
	activity    int
}

func getWorry(info string) []int {
	info = strings.TrimPrefix(info, "  Starting items: ")
	worryStrings := strings.Split(info, ", ")
	worries := make([]int, len(worryStrings))
	for i, w := range worryStrings {
		worries[i], _ = strconv.Atoi(w)
	}
	return worries
}

func getOperation(info string) (int, string) {
	info = strings.TrimPrefix(info, "  Operation: new = old ")
	operationInfo := strings.Split(info, " ")
	operation := 0
	switch {
	case operationInfo[0] == "+":
		operation = add
	case operationInfo[0] == "*":
		operation = multiply
	}
	operand := operationInfo[1]
	return operation, operand
}

func getTest(info string) int {
	info = strings.TrimPrefix(info, "  Test: divisible by ")
	test, _ := strconv.Atoi(info)
	return test
}

func getTrueMonkey(info string) int {
	info = strings.TrimPrefix(info, "    If true: throw to monkey ")
	trueMonkey, _ := strconv.Atoi(info)
	return trueMonkey
}

func getFalseMonkey(info string) int {
	info = strings.TrimPrefix(info, "    If false: throw to monkey ")
	falseMonkey, _ := strconv.Atoi(info)
	return falseMonkey
}

func getMonkeyInfo(info []string) MonkeyInfo {
	var monkeyInfo MonkeyInfo
	monkeyInfo.worries = getWorry(info[1])
	monkeyInfo.operation, monkeyInfo.operand = getOperation(info[2])
	monkeyInfo.testValue = getTest(info[3])
	monkeyInfo.trueMonkey = getTrueMonkey(info[4])
	monkeyInfo.falseMonkey = getFalseMonkey(info[5])

	return monkeyInfo
}

func newWorryValue(worry int, monkey MonkeyInfo) int {
	operandValue := 0
	if monkey.operand == "old" {
		operandValue = worry
	} else {
		operandValue, _ = strconv.Atoi(monkey.operand)
	}
	switch monkey.operation {
	case add:
		return worry + operandValue
	case multiply:
		return worry * operandValue
	}
	return worry
}

func doRound(monkeys []MonkeyInfo) []MonkeyInfo {
	for i := 0; i < len(monkeys); i++ {
		for _, worry := range monkeys[i].worries {
			newWorry := newWorryValue(worry, monkeys[i]) / 3
			newMonkey := monkeys[i].falseMonkey
			if newWorry%monkeys[i].testValue == 0 {
				newMonkey = monkeys[i].trueMonkey
			}
			monkeys[newMonkey].worries = append(monkeys[newMonkey].worries, newWorry)
			monkeys[i].activity++
		}
		monkeys[i].worries = monkeys[i].worries[:0]
	}
	return monkeys
}

func run(input string) {
	monkeyInfo := strings.Split(input, "\n\n")
	monkeys := make([]MonkeyInfo, len(monkeyInfo))
	for i, mi := range monkeyInfo {
		monkeys[i] = getMonkeyInfo(strings.Split(mi, "\n"))
	}

	for r := 0; r < 20; r++ {
		fmt.Println(monkeys)
		fmt.Println()
		monkeys = doRound(monkeys)
	}
	fmt.Println(monkeys)

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].activity > monkeys[j].activity
	})

	fmt.Println(monkeys[0].activity * monkeys[1].activity)
}
