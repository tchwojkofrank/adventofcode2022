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

type Node struct {
	name        string
	size        int
	isDirectory bool
	children    map[string]*Node
	parent      *Node
}

func cd(root *Node, cursor *Node, d string) (*Node, *Node) {
	if strings.HasPrefix(d, "/") {
		return root, root
	} else if d == ".." {
		return root, cursor.parent
	} else {
		return root, cursor.children[d]
	}
	return root, cursor
}

func ls(cursor *Node, output []string) *Node {
	for _, entry := range output {
		fields := strings.Split(entry, " ")
		children := make(map[string]*Node)
		newNode := Node{
			name:        fields[1],
			size:        0,
			isDirectory: false,
			children:    children,
			parent:      cursor,
		}
		if fields[0] == "dir" {
			newNode.isDirectory = true
		} else {
			newNode.size, _ = strconv.Atoi(fields[0])
		}
		cursor.children[newNode.name] = &newNode
	}
	return cursor
}

func processCommands(root *Node, commands []string) *Node {
	cursor := root
	for _, c := range commands {
		output := strings.Split(c, "\n")
		commandLine := output[0]
		output = output[1:]
		command := strings.Split(commandLine, " ")
		switch command[1] {
		case "cd":
			root, cursor = cd(root, cursor, command[2])
		case "ls":
			cursor = ls(cursor, output)
		}
	}
	return root
}

func calcSizes(cursor *Node) int {
	sum := cursor.size
	if cursor.isDirectory {
		if cursor.size == 0 {
			for _, child := range cursor.children {
				sum += calcSizes(child)
			}
		}
		cursor.size = sum
	}
	return sum
}

func part1calc(cursor *Node) int {
	sum := 0
	if cursor.isDirectory {
		if cursor.size <= 100000 {
			sum += cursor.size
		}
		for _, child := range cursor.children {
			sum += part1calc(child)
		}
	}
	return sum
}

const (
	totalSize = 70000000
	needSize  = 30000000
)

func findBest(cursor *Node, needToFree int, best *Node) *Node {
	thisDirSize := 0
	if cursor.isDirectory {
		thisDirSize = cursor.size
	}
	nextBest := best
	if thisDirSize >= needToFree && best.size > thisDirSize {
		nextBest = cursor
	}
	for _, child := range cursor.children {
		if child.isDirectory {
			nextBest = findBest(child, needToFree, nextBest)
		}
	}
	return nextBest
}

func run(input string) {
	children := make(map[string]*Node)
	root := &Node{
		"/",
		0,
		true,
		children,
		nil,
	}

	commands := strings.Split(input, "\n$")
	_ = processCommands(root, commands)
	fmt.Println(root)

	calcSizes(root)
	fmt.Printf("part 1 = %d\n", part1calc(root))
	usedSpace := root.size
	freeSpace := totalSize - usedSpace
	needToFree := needSize - freeSpace

	bestToFree := findBest(root, needToFree, root)
	fmt.Printf("part 2 best directory = %v size = %d\n", bestToFree.name, bestToFree.size)
}
