package main

import (
	"fmt"
	"log"
	"os"
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

func isUnique(t string) bool {
	chars := make(map[rune]struct{}, len(t))
	for _, c := range t {
		if _, ok := chars[c]; ok {
			return false
		}
		chars[c] = struct{}{}
	}
	return true
}

const (
	packetUnique  = 4
	messageUnique = 14
)

func getMarkerIndex(m string, unique int) int {
	for i := unique; i <= len(m); i++ {
		testMarker := m[i-unique : i]
		if isUnique(testMarker) {
			return i
		}
	}
	return -1
}

func run(input string) {
	messages := strings.Split(input, "\n")
	for _, m := range messages {
		fmt.Printf("Start of packet %v: %d\n", m, getMarkerIndex(m, packetUnique))
		fmt.Printf("Start of message %v: %d\n", m, getMarkerIndex(m, messageUnique))
	}
}
