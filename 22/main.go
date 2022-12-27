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
	size = 4
	if inputName == "input" {
		size = 50
	}
	run(text)
	end := time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
}

type Pt struct {
	x, y int
}

type Limit struct {
	min, max int
}

type JungleMap map[Pt]rune
type JungleXBounds map[int]Limit
type JungleYBounds map[int]Limit
type Jungle struct {
	m  JungleMap
	xb JungleXBounds // min and max y coordinate for a given column
	yb JungleYBounds // min and max x coordinate for a given row
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

func makeJungle(input string) Jungle {
	var jungle Jungle
	jungle.m = make(JungleMap)
	jungle.xb = make(JungleXBounds)
	jungle.yb = make(JungleYBounds)
	lines := strings.Split(input, "\n")
	height = len(lines)
	for y, l := range lines {
		for x, r := range l {
			if r != ' ' {
				jungle.m[Pt{x, y}] = r
				if lim, ok := jungle.xb[x]; !ok {
					jungle.xb[x] = Limit{y, y}
				} else {
					jungle.xb[x] = Limit{min(y, lim.min), max(y, lim.max)}
				}
				if lim, ok := jungle.yb[y]; !ok {
					jungle.yb[y] = Limit{x, x}
				} else {
					jungle.yb[y] = Limit{min(x, lim.min), max(x, lim.max)}
				}
			}
		}
	}

	return jungle
}

func (j Jungle) up(p Pt) Pt {
	q := Pt{p.x, p.y - 1}
	if r, ok := j.m[q]; ok {
		switch r {
		case '.':
			return q
		case '#':
			return p
		}
	} else {
		q = Pt{p.x, j.xb[p.x].max}
		r, ok := j.m[q]
		if ok && r == '.' {
			return q
		}
	}
	return p
}

func (j Jungle) down(p Pt) Pt {
	q := Pt{p.x, p.y + 1}
	if r, ok := j.m[q]; ok {
		switch r {
		case '.':
			return q
		case '#':
			return p
		}
	} else {
		q = Pt{p.x, j.xb[p.x].min}
		r, ok := j.m[q]
		if ok && r == '.' {
			return q
		}
	}
	return p
}

func (j Jungle) left(p Pt) Pt {
	q := Pt{p.x - 1, p.y}
	if r, ok := j.m[q]; ok {
		switch r {
		case '.':
			return q
		case '#':
			return p
		}
	} else {
		q = Pt{j.yb[p.y].max, p.y}
		r, ok := j.m[q]
		if ok && r == '.' {
			return q
		}
	}
	return p
}

func (j Jungle) right(p Pt) Pt {
	q := Pt{p.x + 1, p.y}
	if r, ok := j.m[q]; ok {
		switch r {
		case '.':
			return q
		case '#':
			return p
		}
	} else {
		q = Pt{j.yb[p.y].min, p.y}
		r, ok := j.m[q]
		if ok && r == '.' {
			return q
		}
	}
	return p
}

const (
	Right = 0
	Down  = 1
	Left  = 2
	Up    = 3
)

type Location struct {
	p Pt
	d int
}

const (
	TurnLeft  = -1
	TurnRight = 1
)

func getNextCount(directions string) (int, string) {
	i := 0
	value := 0
	for ; i < len(directions) && directions[i] != 'L' && directions[i] != 'R'; i++ {
		value = value*10 + int(directions[i]-'0')
	}
	newDirections := directions[i:]
	return value, newDirections
}

func getNextTurn(directions string) (int, string) {
	switch directions[0] {
	case 'L':
		return TurnLeft, directions[1:]
	case 'R':
		return TurnRight, directions[1:]
	}
	return 0, directions
}

func nextDirectionIsTurn(directions string) bool {
	return directions[0] == 'L' || directions[0] == 'R'
}

func moveOne(jungle Jungle, location Location) Location {
	switch location.d {
	case Up:
		location.p = jungle.up(location.p)
	case Right:
		location.p = jungle.right(location.p)
	case Down:
		location.p = jungle.down(location.p)
	case Left:
		location.p = jungle.left(location.p)
	}
	return location
}

func showMap(j Jungle, l Location) {
	for y := 0; y < height; y++ {
		for x := 0; x <= j.yb[y].max; x++ {
			if x < j.yb[y].min {
				fmt.Print(" ")
			} else {
				if y == l.p.y && x == l.p.x {
					switch l.d {
					case Right:
						fmt.Print(">")
					case Down:
						fmt.Print("v")
					case Left:
						fmt.Print("<")
					case Up:
						fmt.Print("^")
					}
				} else {
					fmt.Printf("%v", string(j.m[Pt{x, y}]))
				}
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func move(jungle Jungle, location Location, directions string) (Location, string) {
	if nextDirectionIsTurn(directions) {
		d, newDirections := getNextTurn(directions)
		location.d = (location.d + d + 4) % 4
		directions = newDirections
	} else {
		m, newDirections := getNextCount(directions)
		for i := 0; i < m; i++ {
			location = moveOne(jungle, location)
		}
		directions = newDirections
	}
	return location, directions
}

func moveAll(jungle Jungle, location Location, directions string) Location {
	for len(directions) > 0 {
		// showMap(jungle, location)
		location, directions = move(jungle, location, directions)
	}
	// showMap(jungle, location)
	return location
}

var height int
var size int

// size x size square map
type Face struct {
	coord      Pt // {x/4, y/4}
	jungle     map[Pt]rune
	topFace    *Face // top edge
	leftFace   *Face // left edge
	bottomFace *Face // bottom edge
	rightFace  *Face // right edge
}

type Die struct {
	faces [6]Face
}

func sign(a int) int {
	if a > 0 {
		return 1
	} else if a < 0 {
		return -1
	}
	return 0
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func rotation(delta Pt) int {
	s := -sign(delta.x * delta.y)
	turns := (abs(delta.x+delta.y-1) % 4)
	return s * turns
}

func cubeMove()

func run(input string) {
	inputs := strings.Split(input, "\n\n")
	jungleInput := inputs[0]
	height = len(jungleInput)
	directions := inputs[1]
	jungle := makeJungle(jungleInput)
	location := Location{p: Pt{jungle.yb[0].min, 0}, d: Right}
	location = moveAll(jungle, location, directions)
	fmt.Printf("Final location: %v\n", location)
	fmt.Printf("Password = %v\n", 1000*(location.p.y+1)+4*(location.p.x+1)+location.d)
}
