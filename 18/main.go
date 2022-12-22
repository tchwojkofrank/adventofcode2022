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

type Pt struct {
	x int
	y int
	z int
}

func newPt(input string) Pt {
	var p Pt
	fields := strings.Split(input, ",")
	p.x, _ = strconv.Atoi(fields[0])
	p.y, _ = strconv.Atoi(fields[1])
	p.z, _ = strconv.Atoi(fields[2])
	return p
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

type PointMap map[Pt]int

func (points *PointMap) neighbors(p Pt) int {
	count, ok := (*points)[p]
	if ok && count >= 0 {
		return count
	}
	count = 0
	for d := -1; d <= 1; d = d + 2 {
		if _, ok := (*points)[Pt{p.x + d, p.y, p.z}]; ok {
			count++
		}
		if _, ok := (*points)[Pt{p.x, p.y + d, p.z}]; ok {
			count++
		}
		if _, ok := (*points)[Pt{p.x, p.y, p.z + d}]; ok {
			count++
		}
	}
	(*points)[p] = count
	return count
}

func run(input string) {
	lines := strings.Split(input, "\n")
	points := PointMap(make(map[Pt]int, len(lines)))
	for _, l := range lines {
		p := newPt(l)
		points[p] = -1
	}
	count := 0
	for k, _ := range points {
		count += points.neighbors(k)
	}
	fmt.Printf("surface area = %v\n", 6*len(points)-count)
}
