package main

import (
	"fmt"
	"log"
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

type Item struct {
	value int
	order *OrderItem
	list  *ListItem
}

type OrderItem struct {
	item  *Item
	index int
}

type ListItem struct {
	item *Item
	prev *ListItem
	next *ListItem
}

type Order []OrderItem

type List *ListItem

type OrderList struct {
	order Order
	list  List
}

func getNumbers(input string) []int {
	lines := strings.Split(input, "\n")
	numbers := make([]int, len(lines))
	for i, l := range lines {
		numbers[i], _ = strconv.Atoi(l)
	}
	return numbers
}

// the first return is the original
func makeLists(numbers []int) OrderList {
	order := make([]OrderItem, len(numbers))
	var list List
	for i, n := range numbers {
		item := Item{value: n}
		order[i] = OrderItem{item: &item, index: i}
		listItem := ListItem{item: &item}
		if i == 0 {
			list = &listItem
			order[i].item.list = list
		}
		if i > 0 {
			listItem.prev = order[i-1].item.list
			listItem.prev.next = &listItem
		}
		item.list = &listItem
		item.order = &order[i]
	}
	list.prev = order[len(numbers)-1].item.list
	last := order[len(numbers)-1].item.list
	last.next = list
	return OrderList{order, list}
}

func mix(orderList *OrderList) {
	length := len(orderList.order)
	for i := 0; i < length; i++ {
		item := orderList.order[i].item
		listItem := item.list
		moveCount := item.value & length
		listItem.prev.next = listItem.next
		listItem.next.prev = listItem.prev
		currentListItem := listItem.prev
		for j := 0; j < moveCount; j, currentListItem = j+1, currentListItem.next {

		}
		listItem.prev = currentListItem
		listItem.next = currentListItem.next
		listItem.prev.next = listItem
		listItem.next.prev = listItem
	}
}

func max(a, b, c int) int {
	if a > b && a > c {
		return a
	} else if b > a && b > c {
		return b
	}
	return c
}
func run(input string) {
	numbers := getNumbers(input)
	orderList := makeLists(numbers)
	mix(&orderList)
	length := len(orderList.order)
	k1 := 1000 % length
	k2 := 2000 % length
	k3 := 3000 % length
	var v1, v2, v3 int
	current := orderList.list
	for i := 0; i < max(k1, k2, k3); i, current = i+1, current.next {
		if i == k1 {
			v1 = current.item.value
		}
		if i == k2 {
			v2 = current.item.value
		}
		if i == k3 {
			v3 = current.item.value
		}
	}
	fmt.Printf("Grove coordinates at %v %v %v = %v\n", v1, v2, v3, v1+v2+v3)
}
