package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
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

func getElfStrings(input string) []string {
	elfStrings := strings.Split(input, "\n\n")
	return elfStrings
}

func getElfInventoryStrings(elfString string) []string {
	inventoryStrings := strings.Split(elfString, "\n")
	return inventoryStrings
}

func getElfInventory(elfInventory []string) []int {
	inventory := make([]int, len(elfInventory))
	for i, es := range elfInventory {
		inventory[i], _ = strconv.Atoi(es)
	}
	return inventory
}

func getElfTotal(elfInventory []int) int {
	sum := 0
	for _, calories := range elfInventory {
		sum += calories
	}
	return sum
}

func getMaxCalores(elfs [][]int) (int, int) {
	max := 0
	elfIndex := -1
	for i, elf := range elfs {
		calories := getElfTotal(elf)
		if calories > max {
			max = getElfTotal(elf)
			elfIndex = i
		}
	}
	return elfIndex, max
}

func getTopThree(elfs [][]int) int {
	sort.Slice(elfs, func(i, j int) bool {
		return getElfTotal(elfs[i]) > getElfTotal(elfs[j])
	})
	return getElfTotal(elfs[0]) + getElfTotal(elfs[1]) + getElfTotal(elfs[2])
}

func run(input string) {
	elfStrings := getElfStrings(input)
	elfs := make([][]int, len(elfStrings))
	fmt.Printf("%v\n", elfStrings)
	for i, es := range elfStrings {
		inventoryStrings := getElfInventoryStrings(es)
		elfs[i] = getElfInventory(inventoryStrings)
		fmt.Printf("%d:\n\t%v\n", i, elfs[i])
	}
	elfIndex, maxCalories := getMaxCalores(elfs)
	fmt.Printf("Max calories: Elf %d, %d\n", elfIndex, maxCalories)
	fmt.Printf("Top three total: %d\n", getTopThree(elfs))
}
