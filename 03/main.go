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

	// Convert []byte to string and print to screen
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
	fmt.Printf("inputName = %v\n", inputName)
	text := readInput(inputName)
	run(text)
}

func getRucksacks(input string) []string {
	return strings.Split(input, "\n")
}

func getCompartments(rucksack string) [2]string {
	var compartments [2]string
	compartments[0] = rucksack[0 : len(rucksack)/2]
	compartments[1] = rucksack[len(rucksack)/2:]
	return compartments
}

func getPriority(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r-'a') + 1
	}
	if r >= 'A' && r <= 'Z' {
		return int(r-'A') + 27
	}
	return 0
}

func getDuplicate(compartments [2]string) rune {
	for _, r := range compartments[0] {
		if strings.ContainsRune(compartments[1], r) {
			return r
		}
	}
	return '-'
}

func getBadge(group []string) rune {
	for _, r := range group[0] {
		if strings.ContainsRune(group[1], r) && strings.ContainsRune(group[2], r) {
			return r
		}
	}
	return '-'
}

func getBadgesSum(rucksacks []string) int {
	sum := 0
	for i := 0; i < len(rucksacks); i += 3 {
		sum += getPriority(getBadge(rucksacks[i : i+3]))
	}
	return sum
}

func run(input string) {
	rucksacks := getRucksacks(input)
	prioritySum := 0
	for _, r := range rucksacks {
		duplicate := getDuplicate(getCompartments(r))
		prioritySum += getPriority(duplicate)
	}
	fmt.Printf("priority sum = %d\n", prioritySum)
	badgeSum := getBadgesSum(rucksacks)
	fmt.Printf("badge sum = %d\n", badgeSum)
}
