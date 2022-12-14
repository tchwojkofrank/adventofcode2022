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

type Pt struct {
	x int
	y int
}

type RockMap map[Pt]rune

func newPt(point string) Pt {
	var pt Pt
	fields := strings.Split(point, ",")
	pt.x, _ = strconv.Atoi(fields[0])
	pt.y, _ = strconv.Atoi(fields[1])
	return pt
}

func sign(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

func toString(rockMap RockMap, min Pt, max Pt) string {
	result := ""
	for y := min.y; y <= max.y; y++ {
		line := ""
		for x := min.x; x <= max.x; x++ {
			if x == 500 && y == 0 {
				line = line + "+"
			} else if r, ok := rockMap[Pt{x, y}]; ok {
				line = line + string(r)
			} else {
				line = line + "."
			}
		}
		result = result + line + "\n"
	}
	return result
}

func (rockMap RockMap) addLine(start Pt, end Pt) {
	dx := sign(end.x - start.x)
	dy := sign(end.y - start.y)
	for x, y := start.x, start.y; x != end.x || y != end.y; x, y = x+dx, y+dy {
		rockMap[Pt{x, y}] = '#'
	}
	rockMap[end] = '#'
}

func (rockMap RockMap) dropSand(minPt Pt, maxPt Pt) bool {
	ok := true
	from := Pt{500, 0}
	abyss := false
	for ok {
		if !(from.x >= minPt.x && from.y >= minPt.y && from.x <= maxPt.x && from.y <= maxPt.y) {
			ok = false
			abyss = true
			continue
		}
		cf, dir := rockMap.canFall(from)
		if cf {
			from = Pt{from.x + dir.x, from.y + dir.y}
		} else {
			ok = false
			rockMap[from] = 'o'
			if from.x == 500 && from.y == 0 {
				abyss = true
			}
		}
	}
	return abyss
}

func newRockMap(rockPaths []string) (RockMap, Pt, Pt) {
	maxPt := Pt{500, 0}
	minPt := Pt{500, 0}
	rockMap := RockMap(make(map[Pt]rune))
	for _, rockPath := range rockPaths {
		rockPoints := strings.Split(rockPath, " -> ")
		for i := 0; i < len(rockPoints)-1; i++ {
			start := newPt(rockPoints[i])
			end := newPt(rockPoints[i+1])
			if start.x < minPt.x {
				minPt.x = start.x
			}
			if end.x < minPt.x {
				minPt.x = end.x
			}
			if start.y < minPt.y {
				minPt.y = start.y
			}
			if end.y < minPt.y {
				minPt.y = end.y
			}
			if start.x > maxPt.x {
				maxPt.x = start.x
			}
			if end.x > maxPt.x {
				maxPt.x = end.x
			}
			if start.y > maxPt.y {
				maxPt.y = start.y
			}
			if end.y > maxPt.y {
				maxPt.y = end.y
			}
			rockMap.addLine(newPt(rockPoints[i]), newPt(rockPoints[i+1]))
		}
	}
	return rockMap, minPt, maxPt
}

func (rockMap RockMap) canFall(from Pt) (bool, Pt) {
	if _, ok := rockMap[Pt{from.x, from.y + 1}]; !ok {
		return true, Pt{0, 1}
	} else if _, ok := rockMap[Pt{from.x - 1, from.y + 1}]; !ok {
		return true, Pt{-1, 1}
	} else if _, ok := rockMap[Pt{from.x + 1, from.y + 1}]; !ok {
		return true, Pt{1, 1}
	}
	return false, Pt{0, 0}
}

func run(input string) {
	rockPaths := strings.Split(input, "\n")
	rockMap, minPt, maxPt := newRockMap(rockPaths)
	fmt.Print(toString(rockMap, minPt, maxPt))
	fmt.Printf("Box: %v, %v\n", minPt, maxPt)
	abyss := false
	count := 0
	for !abyss {
		abyss = rockMap.dropSand(minPt, maxPt)
		if !abyss {
			count++
		}
		fmt.Print(toString(rockMap, minPt, maxPt))
		fmt.Println(count)
	}
}
