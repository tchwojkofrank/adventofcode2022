package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
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
	if len(args) >= 3 {
		f, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	start := time.Now()
	text := readInput(inputName)
	run(text)
	end := time.Now()
	fmt.Printf("Total time: %v\n", end.Sub(start))
}

type Pt struct {
	x int
	y int
}

type RockMap struct {
	rm       map[Pt]rune
	min      Pt
	max      Pt
	lastPath []Pt
}

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

func (rm RockMap) String() string {
	result := ""
	for y := rm.min.y; y <= rm.max.y; y++ {
		line := ""
		for x := rm.min.x; x <= rm.max.x; x++ {
			if x == 500 && y == 0 {
				line = line + "+"
			} else if r, ok := rm.rm[Pt{x, y}]; ok {
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
		rockMap.rm[Pt{x, y}] = '#'
	}
	rockMap.rm[end] = '#'
}

func (rockMap *RockMap) dropSand() bool {
	minPt := rockMap.min
	maxPt := rockMap.max
	abyss := false

	var from Pt

	// We don't have to check from the beginning every time.
	// We can use the previous path the sand took to figure out where the next one lands.
	// This resulted in taking 1/10th the time as without this optimization.
	if len(rockMap.lastPath) > 0 {
		from = rockMap.lastPath[len(rockMap.lastPath)-2]
		rockMap.lastPath = rockMap.lastPath[:len(rockMap.lastPath)-1]
	} else {
		from = Pt{500, 0}
		rockMap.lastPath = append(rockMap.lastPath, from)
	}
	ok := true
	for ok {
		if !(from.x >= minPt.x && from.y >= minPt.y && from.x <= maxPt.x && from.y <= maxPt.y) {
			ok = false
			abyss = true
			continue
		}
		cf, dir := mapCheck(*rockMap, from)
		if cf {
			from = Pt{from.x + dir.x, from.y + dir.y}
			rockMap.lastPath = append(rockMap.lastPath, from)
		} else {
			ok = false
			rockMap.rm[from] = 'o'
			if from.x == 500 && from.y == 0 {
				abyss = true
			}
		}
	}
	return abyss
}

func newRockMap(rockPaths []string) RockMap {
	var rockMap RockMap
	rockMap.rm = make(map[Pt]rune)
	rockMap.max = Pt{500, 0}
	rockMap.min = Pt{500, 0}
	for _, rockPath := range rockPaths {
		rockPoints := strings.Split(rockPath, " -> ")
		for i := 0; i < len(rockPoints)-1; i++ {
			start := newPt(rockPoints[i])
			end := newPt(rockPoints[i+1])
			if start.x < rockMap.min.x {
				rockMap.min.x = start.x
			}
			if end.x < rockMap.min.x {
				rockMap.min.x = end.x
			}
			if start.y < rockMap.min.y {
				rockMap.min.y = start.y
			}
			if end.y < rockMap.min.y {
				rockMap.min.y = end.y
			}
			if start.x > rockMap.max.x {
				rockMap.max.x = start.x
			}
			if end.x > rockMap.max.x {
				rockMap.max.x = end.x
			}
			if start.y > rockMap.max.y {
				rockMap.max.y = start.y
			}
			if end.y > rockMap.max.y {
				rockMap.max.y = end.y
			}
			rockMap.addLine(newPt(rockPoints[i]), newPt(rockPoints[i+1]))
		}
	}
	rockMap.lastPath = make([]Pt, 0)
	return rockMap
}

var mapCheck func(RockMap, Pt) (bool, Pt)

func mapCheck1(rockMap RockMap, pt Pt) (bool, Pt) {
	return rockMap.canFall(pt)
}

func mapCheck2(rockMap RockMap, pt Pt) (bool, Pt) {
	if pt.y == rockMap.max.y-1 {
		return false, Pt{0, 0}
	}
	return mapCheck1(rockMap, pt)
}

func (rockMap RockMap) canFall(from Pt) (bool, Pt) {
	if _, ok := rockMap.rm[Pt{from.x, from.y + 1}]; !ok {
		return true, Pt{0, 1}
	} else if _, ok := rockMap.rm[Pt{from.x - 1, from.y + 1}]; !ok {
		return true, Pt{-1, 1}
	} else if _, ok := rockMap.rm[Pt{from.x + 1, from.y + 1}]; !ok {
		return true, Pt{1, 1}
	}
	return false, Pt{0, 0}
}

func run(input string) {
	mapCheck = mapCheck1
	rockPaths := strings.Split(input, "\n")
	rockMap := newRockMap(rockPaths)
	abyss := false
	count := 0
	for !abyss {
		abyss = rockMap.dropSand()
		if !abyss {
			count++
		}
	}
	fmt.Printf("Part 1 count = %v\n", count)

	//part 2
	mapCheck = mapCheck2
	rockMap = newRockMap(rockPaths)
	rockMap.max.y += 2
	if rockMap.min.x > 500-rockMap.max.y-3 {
		rockMap.min.x = 500 - rockMap.max.y - 3
	}
	if rockMap.max.x < 500+rockMap.max.y+3 {
		rockMap.max.x = 500 + rockMap.max.y + 3
	}
	abyss = false
	count = 0
	for !abyss {
		abyss = rockMap.dropSand()
		count++
	}
	fmt.Println(count)

}
