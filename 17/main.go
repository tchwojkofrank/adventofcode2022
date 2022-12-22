package main

import (
	"fmt"
	"log"
	"os"
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

type Pt struct {
	x int
	y int
}

const (
	width  = 7
	xStart = 2
)

var shapeHeight = [5]int{1, 3, 3, 4, 2}
var shapeWidth = [5]int{4, 3, 3, 1, 2}

func initShapes() [][]Pt {
	shapes := make([][]Pt, 5)
	shapes[0] = []Pt{{0, 0}, {1, 0}, {2, 0}, {3, 0}}
	shapes[1] = []Pt{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}}
	shapes[2] = []Pt{{2, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2}}
	shapes[3] = []Pt{{0, 0}, {0, 1}, {0, 2}, {0, 3}}
	shapes[4] = []Pt{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	return shapes
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

type Jets struct {
	jets  []int
	index int
}

type Shapes struct {
	shapes [][]Pt
	index  int
}

type RockMap struct {
	rocks  map[Pt]struct{}
	height int
}

func isCollision(rockMap *RockMap, shape []Pt, offset Pt) bool {
	for _, p := range shape {
		newP := Pt{p.x + offset.x, p.y + offset.y}
		if _, ok := rockMap.rocks[newP]; ok {
			return true
		}
	}
	return false
}

func dropRock(rockMap *RockMap, shapes *Shapes, jets *Jets) {
	offset := Pt{2, -rockMap.height - 3 - shapeHeight[shapes.index]}
	collision := false
	for !collision {
		j := jets.jets[jets.index]
		jets.index = (jets.index + 1) % len(jets.jets)
		if (offset.x+j >= 0) && (offset.x+j+shapeWidth[shapes.index]-1 < 7) {
			offset.x += j
			if isCollision(rockMap, shapes.shapes[shapes.index], offset) {
				offset.x -= j
			}
		}
		offset.y++
		if isCollision(rockMap, shapes.shapes[shapes.index], offset) {
			offset.y--
			collision = true
		}
	}
	for _, p := range shapes.shapes[shapes.index] {
		rockMap.rocks[Pt{p.x + offset.x, p.y + offset.y}] = struct{}{}
		if -(p.y + offset.y) > rockMap.height {
			rockMap.height = -(p.y + offset.y)
		}
	}
	shapes.index = (shapes.index + 1) % len(shapes.shapes)
}

func (rm RockMap) String() string {
	result := ""
	for y := -rm.height - 4; y <= 0; y++ {
		line := "|"
		for x := 0; x < 7; x++ {
			if _, ok := rm.rocks[Pt{x, y}]; ok {
				line = line + "#"
			} else {
				line = line + "."
			}
		}
		line = line + "|\n"
		result = result + line
	}
	return result
}

func run(input string) {
	var jets Jets
	var shapes Shapes
	var rockMap RockMap
	jets.jets = make([]int, len(input))
	for i, r := range input {
		switch r {
		case '<':
			jets.jets[i] = -1
		case '>':
			jets.jets[i] = 1
		}
	}
	shapes.shapes = initShapes()
	rockMap.rocks = make(map[Pt]struct{})
	for i := 0; i < 7; i++ {
		rockMap.rocks[Pt{i, 0}] = struct{}{}
		rockMap.height = 0
	}
	for i := 0; i < 2022; i++ {
		dropRock(&rockMap, &shapes, &jets)
	}
	fmt.Printf("Height = %v\n", rockMap.height)
}
