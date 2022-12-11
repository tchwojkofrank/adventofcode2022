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
	worries     []int64
	operation   int
	operand     string
	testValue   int
	trueMonkey  int
	falseMonkey int
	activity    int
}

func getWorry(info string) []int64 {
	info = strings.TrimPrefix(info, "  Starting items: ")
	worryStrings := strings.Split(info, ", ")
	worries := make([]int64, len(worryStrings))
	for i, w := range worryStrings {
		value, _ := strconv.Atoi(w)
		worries[i] = int64(value)
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

func newWorryValue(worry int64, monkey MonkeyInfo) int64 {
	var operandValue int64
	if monkey.operand == "old" {
		operandValue = worry
	} else {
		ov, _ := strconv.Atoi(monkey.operand)
		operandValue = int64(ov)
	}
	switch monkey.operation {
	case add:
		z := worry + operandValue
		return z
	case multiply:
		z := worry * operandValue
		return z
	}
	return worry
}

func doRound(monkeys []MonkeyInfo, worryFactor int64, worryLimit int64) []MonkeyInfo {
	for i := 0; i < len(monkeys); i++ {
		for _, worry := range monkeys[i].worries {
			newWorry := (newWorryValue(worry, monkeys[i]) / worryFactor) % worryLimit
			newMonkey := monkeys[i].falseMonkey
			testResult := newWorry % int64(monkeys[i].testValue)
			if testResult == 0 {
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
	var worryLimit int64
	worryLimit = 1
	for i, mi := range monkeyInfo {
		monkeys[i] = getMonkeyInfo(strings.Split(mi, "\n"))
		worryLimit = worryLimit * int64(monkeys[i].testValue)
	}

	for r := 0; r < 20; r++ {
		monkeys = doRound(monkeys, 3, worryLimit)
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].activity > monkeys[j].activity
	})

	fmt.Println(monkeys[0].activity * monkeys[1].activity)

	monkeys = make([]MonkeyInfo, len(monkeyInfo))
	worryLimit = 1
	for i, mi := range monkeyInfo {
		monkeys[i] = getMonkeyInfo(strings.Split(mi, "\n"))
		worryLimit = worryLimit * int64(monkeys[i].testValue)
	}

	for r := 0; r < 10000; r++ {
		monkeys = doRound(monkeys, 1, worryLimit)
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].activity > monkeys[j].activity
	})

	fmt.Println(monkeys[0].activity * monkeys[1].activity)

}
