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
	if len(args) < 2 {
		fmt.Printf("Enter input file name.\n")
		return
	}
	params := os.Args[1]
	inputName := strings.Split(params, " ")[0]
	text := readInput(inputName)
	run(text)
}

type CPU struct {
	register            int
	cycle               int
	totalSignalStrength int
	addState            int
	addValue            int
}

const (
	noop      = 0
	startAdd  = 1
	finishAdd = 2
)

func (cpu *CPU) start() {
	fmt.Printf("Start : %v\n", *cpu)
	cpu.cycle++
}

func (cpu *CPU) during(command string, param int) {
	if (cpu.cycle-20)%40 == 0 {
		cpu.totalSignalStrength += cpu.cycle * cpu.register
	}
	switch {
	case command == "noop":
		cpu.addState = noop
		cpu.addValue = 0
	case command == "addx":
		switch cpu.addState {
		case noop:
			cpu.addState = startAdd
			cpu.addValue = param
		case startAdd:
			cpu.addState = finishAdd
		}
	}
	fmt.Printf("During: %v\n", *cpu)
}

func (cpu *CPU) end(command string, param int) {
	switch {
	case command == "noop":
	case command == "addx":
		switch cpu.addState {
		case noop:
		case startAdd:
		case finishAdd:
			cpu.register += param
			cpu.addState = noop
			cpu.addValue = 0
		}
	}
	fmt.Printf("End   : %v\n\n", *cpu)
}

func (cpu *CPU) step(instruction []string) {
	command := instruction[0]
	param := 0
	if command == "addx" {
		param, _ = strconv.Atoi(instruction[1])
	}
	cpu.start()
	cpu.during(command, param)
	cpu.end(command, param)
	if command == "addx" {
		cpu.start()
		cpu.during(command, param)
		cpu.end(command, param)
	}
}

func run(input string) {
	instructions := strings.Split(input, "\n")
	cpu := CPU{
		register: 1,
	}
	for _, instruction := range instructions {
		cpu.step(strings.Split(instruction, " "))
	}
	fmt.Printf("\n\nTotal signal strength = %d\n", cpu.totalSignalStrength)
}
