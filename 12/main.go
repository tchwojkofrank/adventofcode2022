package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"chwojkofrank.com/dijkstra"
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

func getElevation(r rune) int {
	if r == 'S' {
		return 0
	}
	if r == 'E' {
		return int('z' - 'a')
	}
	return int(r - 'a')
}

type Pos struct {
	x int
	y int
}

var isNeighbor func(p Pos, n Pos) bool

func isNeighborUp(p Pos, n Pos) bool {
	pHeight, ok := elevationMap[p]
	if !ok {
		return false
	}
	nHeight, ok := elevationMap[n]
	if !ok {
		return false
	}
	return nHeight-pHeight <= 1
}

func isNeighborDown(p Pos, n Pos) bool {
	pHeight, ok := elevationMap[p]
	if !ok {
		return false
	}
	nHeight, ok := elevationMap[n]
	if !ok {
		return false
	}
	return nHeight-pHeight >= -1
}

// func abs(x int) int {
// 	if x < 0 {
// 		return -x
// 	}
// 	return x
// }

func (p Pos) Neighbors() ([]dijkstra.Node, []int) {
	neighbors := make([]dijkstra.Node, 0)
	distances := make([]int, 0)
	if isNeighbor(p, Pos{p.x - 1, p.y}) {
		neighbors = append(neighbors, Pos{p.x - 1, p.y})
		distances = append(distances, 1)
	}
	if isNeighbor(p, Pos{p.x + 1, p.y}) {
		neighbors = append(neighbors, Pos{p.x + 1, p.y})
		distances = append(distances, 1)
	}
	if isNeighbor(p, Pos{p.x, p.y - 1}) {
		neighbors = append(neighbors, Pos{p.x, p.y - 1})
		distances = append(distances, 1)
	}
	if isNeighbor(p, Pos{p.x, p.y + 1}) {
		neighbors = append(neighbors, Pos{p.x, p.y + 1})
		distances = append(distances, 1)
	}
	return neighbors, distances
}

var elevationMap map[Pos]int

func run(input string) {
	rows := strings.Split(input, "\n")
	var start Pos
	var end Pos
	allNodes := make([]dijkstra.Node, 0)
	// allReverseNodes := make([]dijkstra.Node, 0)
	isNeighbor = isNeighborUp

	elevationMap = make(map[Pos]int)
	for y := range rows {
		for x, r := range rows[y] {
			allNodes = append(allNodes, Pos{x, y})
			elevationMap[Pos{x, y}] = getElevation(r)
			if r == 'S' {
				start = Pos{x, y}
			} else if r == 'E' {
				end = Pos{x, y}
			}
		}
	}
	path := dijkstra.GetShortestPath(allNodes, start, end)
	fmt.Println(path)
	fmt.Println(len(path) - 1)

	isNeighbor = isNeighborDown
	distances, prev := dijkstra.GetShortestDistances(allNodes, Pos(end))
	bestDistance := 10000
	var bestNode dijkstra.Node
	for _, n := range allNodes {
		if elevationMap[n.(Pos)] == 0 {
			if distances[n] < bestDistance {
				bestDistance = distances[n]
				bestNode = n
			}
		}
	}
	reversepath := make([]dijkstra.Node, 0)
	u := bestNode
	for u != nil {
		reversepath = append([]dijkstra.Node{u}, reversepath...)
		u = prev[u]
	}
	fmt.Printf("Path = %v, length = %v\n", reversepath, len(reversepath))
	fmt.Printf("Best starting point at %v with distance %v\n", bestNode, bestDistance)

}
