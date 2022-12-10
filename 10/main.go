package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	cursor "chwojkofrank.com/cursor"
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

type CRT [6][40]byte

type CPU struct {
	register            int
	cycle               int
	totalSignalStrength int
	addState            int
	addValue            int
	crt                 CRT
}

const (
	noop      = 0
	startAdd  = 1
	finishAdd = 2
)

func (cpu *CPU) isSpriteVisible() byte {
	spritePosition := cpu.register
	currentLineIndex := (cpu.cycle - 1) % 40
	if currentLineIndex >= spritePosition-1 && currentLineIndex <= spritePosition+1 {
		return '#'
	}
	return '.'
}

func (cpu *CPU) start() {
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
	crtLine := (cpu.cycle - 1) / 40
	crtOffset := (cpu.cycle - 1) % 40
	cpu.crt[crtLine][crtOffset] = cpu.isSpriteVisible()
	cursor.Clear()
	cursor.Position(0, 0)
	fmt.Printf("%v\n", cpu.crt)
	time.Sleep(10 * time.Millisecond)
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

func (cpu CPU) String() string {
	result := fmt.Sprintf("cycle: %d\nregister: %d\n\n", cpu.cycle, cpu.register)
	return result
}

func (crt CRT) String() string {
	result := ""
	for i := 0; i < 6; i++ {
		result = result + string(crt[i][:]) + "\n"
	}
	return result
}

func run(input string) {
	instructions := strings.Split(input, "\n")
	cpu := CPU{
		register: 1,
	}
	for _, instruction := range instructions {
		cpu.step(strings.Split(instruction, " "))
	}
	fmt.Printf("\n\nTotal signal strength = %d\n\n", cpu.totalSignalStrength)

	fmt.Printf("%v", cpu.crt)
}
