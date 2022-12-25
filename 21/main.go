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
	name        string
	value       int
	valueString string
	operation   string
	operands    *[2]string
}

type Monkeys map[string]Monkey

func newMonkey(info string) Monkey {
	var monkey Monkey
	instructions := strings.Split(info, " ")
	monkey.name = strings.TrimSuffix(instructions[0], ":")
	if len(instructions) == 2 {
		var ok error
		monkey.value, ok = strconv.Atoi(instructions[1])
		if ok != nil {
			monkey.value = math.MaxInt
		}
		monkey.valueString = instructions[1]
	} else if len(instructions) == 4 {
		var operands [2]string
		operands[0] = instructions[1]
		operands[1] = instructions[3]
		monkey.operands = &operands
		monkey.operation = instructions[2]
		monkey.value = math.MaxInt
		monkey.valueString = ""
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
		if value1 == math.MaxInt || value2 == math.MaxInt {
			return math.MaxInt
		}
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

func (monkeys Monkeys) monkeyEquation(monkey Monkey) (string, int) {
	if monkey.operands == nil {
		return monkey.valueString, monkey.value
	} else {
		result := "( "
		result1, value1 := monkeys.monkeyEquation(monkeys[monkey.operands[0]])
		result += result1
		result += " "
		result += monkey.operation + " "
		result2, value2 := monkeys.monkeyEquation(monkeys[monkey.operands[1]])
		result += result2
		result += " ) "
		var mValue int
		if value1 == math.MaxInt || value2 == math.MaxInt {
			mValue = math.MaxInt
		} else {
			switch monkey.operation {
			case "+":
				mValue = value1 + value2
			case "-":
				mValue = value1 - value2
			case "*":
				mValue = value1 * value2
			case "/":
				mValue = value1 / value2
			default:
				mValue = math.MaxInt
			}
		}
		return result, mValue
	}
}

func (monkeys Monkeys) calcShout(unknown string, target int) int {
	m := monkeys[unknown]
	if m.name == "humn" {
		return target
	}
	if m.value != math.MaxInt {
		return m.value
	}
	m1 := m.operands[0]
	m2 := m.operands[1]
	m1Value := monkeys.monkeyNumber(monkeys[m1])
	m2Value := monkeys.monkeyNumber(monkeys[m2])
	if m1Value == math.MaxInt {
		var newTarget int
		switch m.operation {
		case "+":
			newTarget = target - m2Value
		case "-":
			newTarget = target + m2Value
		case "*":
			newTarget = target / m2Value
		case "/":
			newTarget = target * m2Value
		}
		return monkeys.calcShout(m1, newTarget)
	} else if m2Value == math.MaxInt {
		var newTarget int
		switch m.operation {
		case "+":
			newTarget = target - m1Value
		case "-":
			newTarget = m1Value - target
		case "*":
			newTarget = target / m1Value
		case "/":
			newTarget = m1Value / target
		}
		return monkeys.calcShout(m2, newTarget)
	} else {
		switch m.operation {
		case "+":
			return m1Value + m2Value
		case "-":
			return m1Value - m2Value
		case "*":
			return m1Value * m2Value
		case "/":
			return m1Value / m2Value
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

	monkeys = Monkeys(make(map[string]Monkey))
	for _, m := range monkeyInfo {
		monkey := newMonkey(m)
		monkeys[monkey.name] = monkey
	}

	root := monkeys["root"]
	root.operation = "="
	monkeys["root"] = root
	humn := monkeys["humn"]
	humn.valueString = "X"
	humn.value = math.MaxInt
	monkeys["humn"] = humn

	equation, _ := monkeys.monkeyEquation(root)
	fmt.Printf("Equation: \n%v\n", equation)

	monkeys = Monkeys(make(map[string]Monkey))
	for _, m := range monkeyInfo {
		monkey := newMonkey(m)
		monkeys[monkey.name] = monkey
	}

	root = monkeys["root"]
	root.operation = "="
	monkeys["root"] = root
	humn = monkeys["humn"]
	humn.valueString = "X"
	humn.value = math.MaxInt
	monkeys["humn"] = humn

	m1name := root.operands[0]
	m2name := root.operands[1]
	m1Value := monkeys.monkeyNumber(monkeys[m1name])
	m2Value := monkeys.monkeyNumber(monkeys[m2name])
	var shout int
	if m1Value == math.MaxInt {
		shout = monkeys.calcShout(m1name, m2Value)
	} else if m2Value == math.MaxInt {
		shout = monkeys.calcShout(m2name, m1Value)
	}
	fmt.Printf("Shout %v\n", shout)
}
