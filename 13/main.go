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

func splitList(list string) []string {
	items := make([]string, 0)
	for i := 0; i < len(list); {
		openBracket := 0
		switch list[i] {
		case '[':
			startIndex := i
			openBracket++
			for i = i + 1; i < len(list) && openBracket > 0; i++ {
				switch list[i] {
				case '[':
					openBracket++
				case ']':
					openBracket--
				}
			}
			endIndex := i
			items = append(items, list[startIndex:endIndex])
		case ',':
			i++
		default:
			startIndex := i
			for i = i + 1; i < len(list) && list[i] != ','; i++ {
			}
			endIndex := i
			items = append(items, list[startIndex:endIndex])
		}
	}
	return items
}

func isPairCorrect(list1 string, list2 string) int {
	if list1[0] == '[' && list2[0] == '[' {
		items1 := splitList(list1[1 : len(list1)-1])
		items2 := splitList(list2[1 : len(list2)-1])
		for i, j := 0, 0; i < len(items1) && j < len(items2); i, j = i+1, j+1 {
			switch isPairCorrect(items1[i], items2[j]) {
			case -1:
				return -1
			case 0:
				// keep comparing
			case 1:
				return 1
			}
		}
		if len(items1) < len(items2) {
			return 1
		} else if len(items1) > len(items2) {
			return -1
		} else {
			return 0
		}
	} else if list1[0] == '[' || list2[0] == '[' {
		if list1[0] != '[' {
			list1 = "[" + list1 + "]"
		}
		if list2[0] != '[' {
			list2 = "[" + list2 + "]"
		}
		return isPairCorrect(list1, list2)
	} else {
		val1, _ := strconv.Atoi(list1)
		val2, _ := strconv.Atoi(list2)
		if val1 < val2 {
			return 1
		} else if val1 > val2 {
			return -1
		} else {
			return 0
		}
	}
	return 0
}

// func listValue(list string) int {
// 	list = strings.ReplaceAll(list, "[", "")
// 	list = strings.ReplaceAll(list, "]", "")
// 	items := strings.Split(list, ",")
// 	sum := 0
// 	for _, i := range items {
// 		value, _ := strconv.Atoi(i)
// 		sum = sum + value
// 	}
// 	return sum
// }

func run(input string) {
	pairs := strings.Split(input, "\n\n")
	correctLists := make([]int, 0)
	sum := 0
	packets := make([]string, 0)
	for i, pair := range pairs {
		lists := strings.Split(pair, "\n")
		packets = append(packets, lists...)
		if isPairCorrect(lists[0], lists[1]) == 1 {
			correctLists = append(correctLists, i+1)
			sum = sum + i + 1
		}
	}
	fmt.Printf("Correct lists: %v\n\n", correctLists)
	fmt.Printf("Correct sum: %v\n\n", sum)

	packets = append(packets, "[[2]]", "[[6]]")
	sort.Slice(packets, func(i, j int) bool {
		return isPairCorrect(packets[i], packets[j]) == 1
	})
	fmt.Println(packets)
	v1, v2 := 0, 0
	for i, p := range packets {
		if p == "[[2]]" {
			v1 = i + 1
		}
		if p == "[[6]]" {
			v2 = i + 1
		}
	}
	fmt.Printf("%d * %d = %d\n", v1, v2, v1*v2)
}
