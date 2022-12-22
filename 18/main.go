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

func (outside *PointMap) outsideNeighbors(p Pt) int {
	count := 0
	for d := -1; d <= 1; d = d + 2 {
		if _, ok := (*outside)[Pt{p.x + d, p.y, p.z}]; ok {
			count++
		}
		if _, ok := (*outside)[Pt{p.x, p.y + d, p.z}]; ok {
			count++
		}
		if _, ok := (*outside)[Pt{p.x, p.y, p.z + d}]; ok {
			count++
		}
	}
	return count
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

func ptMin(a, b Pt) Pt {
	return Pt{min(a.x, b.x), min(a.y, b.y), min(a.z, b.z)}
}

func ptMax(a, b Pt) Pt {
	return Pt{max(a.x, b.x), max(a.y, b.y), max(a.z, b.z)}
}

func (outside *PointMap) addNeighbors(inside *PointMap, p Pt, min Pt, max Pt) {
	if p.x < min.x || p.y < min.y || p.y < min.z {
		return
	}
	if p.x > max.x || p.y > max.y || p.y > max.z {
		return
	}
	if _, ok := (*inside)[p]; ok {
		return
	}
	if _, ok := (*outside)[p]; ok {
		return
	}

	var px [2]Pt
	var py [2]Pt
	var pz [2]Pt
	var okx [2]bool
	var oky [2]bool
	var okz [2]bool
	for d, i := -1, 0; d <= 1; d, i = d+2, i+1 {
		px[i] = Pt{p.x + d, p.y, p.z}
		py[i] = Pt{p.x, p.y + d, p.z}
		pz[i] = Pt{p.x, p.y, p.z + d}
		_, okx[i] = (*outside)[px[i]]
		_, oky[i] = (*outside)[py[i]]
		_, okz[i] = (*outside)[pz[i]]
	}
	if okx[0] || okx[1] || oky[0] || oky[1] || okz[0] || okz[1] {
		(*outside)[p] = -1
	}
	for i := 0; i <= 1; i = i + 1 {
		outside.addNeighbors(inside, px[i], min, max)
		outside.addNeighbors(inside, py[i], min, max)
		outside.addNeighbors(inside, pz[i], min, max)
	}
}

func run(input string) {
	lines := strings.Split(input, "\n")
	points := PointMap(make(map[Pt]int))
	min := Pt{0, 0, 0}
	max := Pt{0, 0, 0}
	for _, l := range lines {
		p := newPt(l)
		min = ptMin(min, p)
		max = ptMax(max, p)
		points[p] = -1
	}
	surfaceArea := 0
	for k := range points {
		surfaceArea += points.neighbors(k)
	}
	fmt.Printf("surface area = %v\n", 6*len(points)-surfaceArea)

	outside := PointMap(make(map[Pt]int))
	for z := min.z - 1; z <= max.z+1; z++ {
		for y := min.y - 1; y <= max.y+1; y++ {
			for x := min.x - 1; x <= max.x+1; x++ {
				p := Pt{x, y, z}
				if x == min.x-1 || x == max.x+1 || y == min.y-1 || y == max.y+1 || z == min.z-1 || z == max.z+1 {
					outside[p] = -1
				}
			}
		}
	}
	outside.addNeighbors(&points, min, min, max)
	fmt.Printf("From %v to %v, %v outside points, %v lava points\n",
		Pt{min.x - 1, min.y - 1, min.z - 1},
		Pt{max.x + 1, max.y + 1, max.z + 1},
		len(outside), len(points))

	exteriorSurfaceArea := 0
	for k := range points {
		exteriorSurfaceArea = exteriorSurfaceArea + outside.outsideNeighbors(k)
	}
	fmt.Printf("Exterior surface area = %v\n", exteriorSurfaceArea)

}
