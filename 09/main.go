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

type Pos struct {
	x int
	y int
}

var up = Pos{0, 1}
var down = Pos{0, -1}
var left = Pos{-1, 0}
var right = Pos{1, 0}
var directionMap = map[string]Pos{
	"U": up,
	"D": down,
	"L": left,
	"R": right,
}

type Rope []Pos

type Move struct {
	direction Pos
	count     int
}

func movePos(start Pos, delta Pos) Pos {
	return Pos{start.x + delta.x, start.y + delta.y}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	} else if x == 0 {
		return 0
	}
	return 1
}

/* returns the adjustment needed for the next know based on the current knot's position */
func getNextMove(rope Rope, i int) Pos {
	dx := rope[i].x - rope[i+1].x
	dy := rope[i].y - rope[i+1].y
	if abs(dx) > 1 || abs(dy) > 1 {
		return Pos{sign(dx), sign(dy)}
	}
	return Pos{0, 0}
}

func move(rope Rope, move Move, tailmap map[Pos]struct{}) (Rope, map[Pos]struct{}) {
	for i := 0; i < move.count; i++ {
		rope[0] = movePos(rope[0], move.direction)
		for k := 0; k < len(rope)-1; k++ {
			nextmove := getNextMove(rope, k)
			rope[k+1] = movePos(rope[k+1], nextmove)
		}
		tailmap[rope[len(rope)-1]] = struct{}{}
	}
	return rope, tailmap
}

func moveAll(rope Rope, moves []Move, tailmap map[Pos]struct{}) (Rope, map[Pos]struct{}) {
	for _, m := range moves {
		rope, tailmap = move(rope, m, tailmap)
	}
	return rope, tailmap
}

func getMoves(input string) []Move {
	moveStrings := strings.Split(input, "\n")
	moves := make([]Move, len(moveStrings))
	for i, m := range moveStrings {
		fields := strings.Split(m, " ")
		moves[i].direction = directionMap[fields[0]]
		moves[i].count, _ = strconv.Atoi(fields[1])
	}
	return moves
}

func createRope(knots int) Rope {
	return Rope(make([]Pos, knots))
}

func run(input string) {
	moves := getMoves(input)
	rope := createRope(2)
	tailmap := make(map[Pos]struct{})
	tailmap[Pos{0, 0}] = struct{}{}
	_, tailmap = moveAll(rope, moves, tailmap)
	fmt.Printf("Part 1: Tail visited %d positions.\n", len(tailmap))

	rope = createRope(10)
	tailmap = make(map[Pos]struct{})
	tailmap[Pos{0, 0}] = struct{}{}
	_, tailmap = moveAll(rope, moves, tailmap)
	fmt.Printf("Part 2: Tail visited %d positions.\n", len(tailmap))
}
