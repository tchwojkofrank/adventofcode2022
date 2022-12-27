package main

import (
	"fmt"
	"log"
	"math"
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

type Pt struct {
	x, y int
}

var NW = Pt{-1, -1}
var N = Pt{0, -1}
var NE = Pt{1, -1}
var W = Pt{-1, 0}
var E = Pt{1, 0}
var SW = Pt{-1, 1}
var S = Pt{0, 1}
var SE = Pt{1, 1}

type Direction struct {
	d     Pt
	check []Pt
}

var North = []Pt{NW, N, NE}
var East = []Pt{NE, E, SE}
var South = []Pt{SW, S, SE}
var West = []Pt{NW, W, SW}
var All = []Pt{NW, N, NE, E, SE, S, SW, W}

func add(a Pt, b Pt) Pt {
	return Pt{a.x + b.x, a.y + b.y}
}

func checkDirection(elves Elves, p Pt, direction []Pt) int {
	count := 0

	for _, d := range direction {
		count += elves[add(p, d)]
	}

	return count
}

type ProposedTargets map[Pt]int
type ProposedDestinations map[Pt]Pt
type Elves map[Pt]int
type Bounds struct {
	min, max Pt
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func boundsSize(b Bounds) int {
	return (b.max.x - b.min.x + 1) * (b.max.y - b.min.y + 1)
}

func getElves(input string) (Elves, Bounds) {
	lines := strings.Split(input, "\n")
	elves := make(Elves)
	bounds := Bounds{Pt{math.MaxInt, math.MaxInt}, Pt{math.MinInt, math.MinInt}}
	for y, l := range lines {
		for x, r := range l {
			if r == '#' {
				p := Pt{x, y}
				elves[p] = 1
				bounds.min = Pt{min(bounds.min.x, p.x), min(bounds.min.y, p.y)}
				bounds.max = Pt{max(bounds.max.x, p.x), max(bounds.max.y, p.y)}
			}
		}
	}
	return elves, bounds
}

func doRound(elves Elves, directions []Direction) (Elves, Bounds, int) {
	targets := make(ProposedTargets)
	destinations := make(ProposedDestinations)
	movedCount := 0
	for e, _ := range elves {
		if checkDirection(elves, e, All) != 0 {
			found := false
			for _, dir := range directions {
				if checkDirection(elves, e, dir.check) == 0 {
					targets[add(e, dir.d)] = targets[add(e, dir.d)] + 1
					destinations[e] = add(e, dir.d)
					found = true
					movedCount++
					break
				}
			}
			if !found {
				destinations[e] = e
			}
		} else {
			destinations[e] = e
		}
	}
	newElves := make(Elves)
	bounds := Bounds{Pt{math.MaxInt, math.MaxInt}, Pt{math.MinInt, math.MinInt}}
	for e, _ := range elves {
		d := e
		if targets[destinations[e]] == 1 {
			d = destinations[e]
		}
		if newElves[d] > 0 {
			fmt.Printf("Collision at %v\n", d)
			for k, v := range destinations {
				if v == d {
					fmt.Printf("%v has destination %v\n", k, v)
				}
			}
			panic("Collision!")
		}
		newElves[d] = 1
		bounds.min = Pt{min(bounds.min.x, d.x), min(bounds.min.y, d.y)}
		bounds.max = Pt{max(bounds.max.x, d.x), max(bounds.max.y, d.y)}
	}
	return newElves, bounds, movedCount
}

func showElves(elves Elves, bounds Bounds) {
	for y := bounds.min.y; y <= bounds.max.y; y++ {
		for x := bounds.min.x; x <= bounds.max.x; x++ {
			if elves[Pt{x, y}] > 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func run(input string) {
	elves, bounds := getElves(input)
	fmt.Printf("After round %d\n", 0)
	showElves(elves, bounds)
	directions := make([]Direction, 4)
	directions[0] = Direction{d: N, check: North}
	directions[1] = Direction{d: S, check: South}
	directions[2] = Direction{d: W, check: West}
	directions[3] = Direction{d: E, check: East}
	movedCount := 0
	i := 0
	for ; i < 10; i++ {
		elves, bounds, movedCount = doRound(elves, directions)
		first := directions[0]
		directions = append(directions[1:], first)
		fmt.Printf("After round %d, %d elves moved\n", i+1, movedCount)
		// showElves(elves, bounds)
	}
	fmt.Printf("empty space = %d\n", boundsSize(bounds)-len(elves))
	fmt.Printf("After round %d, %d elves moved\n", i, movedCount)
	for ; movedCount > 0; i++ {
		elves, bounds, movedCount = doRound(elves, directions)
		first := directions[0]
		directions = append(directions[1:], first)
		fmt.Printf("After round %d, %d elves moved\n", i+1, movedCount)
		// showElves(elves, bounds)
	}
	showElves(elves, bounds)
	fmt.Printf("First round no elves moved = %v\n", i)
}
