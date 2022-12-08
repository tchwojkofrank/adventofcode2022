package main

import (
	"fmt"
	"log"
	"os"
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

func makeMap(input string) [][]int {
	rows := strings.Split(input, "\n")
	width := len(rows[0])
	height := len(rows)
	treemap := make([][]int, height)
	for j, row := range rows {
		treemap[j] = make([]int, width)
		for i := 0; i < len(row); i++ {
			treemap[j][i] = int(row[i] - '0')
		}
	}
	return treemap
}

func getDimensions(treemap [][]int) (int, int) {
	height := len(treemap)
	width := len(treemap[0])
	return width, height
}

func isVisibleInDirection(treemap [][]int, x int, y int, dx int, dy int) bool {
	width, height := getDimensions(treemap)
	treeheight := treemap[y][x]
	for i, j := x+dx, y+dy; i >= 0 && i < width && j >= 0 && j < height; i, j = i+dx, j+dy {
		if treeheight <= treemap[j][i] {
			return false
		}
	}
	return true
}

func isVisible(treemap [][]int, x int, y int) bool {
	return isVisibleInDirection(treemap, x, y, 0, -1) ||
		isVisibleInDirection(treemap, x, y, 1, 0) ||
		isVisibleInDirection(treemap, x, y, 0, 1) ||
		isVisibleInDirection(treemap, x, y, -1, 0)
}

func countVisible(treemap [][]int) int {
	width, height := getDimensions(treemap)
	count := 0
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			if isVisible(treemap, i, j) {
				count = count + 1
			}
		}
	}
	return count
}

func scenicScoreInDirection(treemap [][]int, x int, y int, dx int, dy int) int {
	width, height := getDimensions(treemap)
	treeheight := treemap[y][x]
	score := 0
	for i, j := x+dx, y+dy; i >= 0 && i < width && j >= 0 && j < height; i, j = i+dx, j+dy {
		score = score + 1
		if treeheight <= treemap[j][i] {
			break
		}
	}
	return score
}

func scenicScore(treemap [][]int, x int, y int) int {
	return scenicScoreInDirection(treemap, x, y, 0, -1) *
		scenicScoreInDirection(treemap, x, y, 1, 0) *
		scenicScoreInDirection(treemap, x, y, 0, 1) *
		scenicScoreInDirection(treemap, x, y, -1, 0)
}

func bestScenicScore(treemap [][]int) (int, int, int) {
	width, height := getDimensions(treemap)
	bestScore := 0
	besti := -1
	bestj := -1
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			score := scenicScore(treemap, i, j)
			if score > bestScore {
				bestScore, besti, bestj = score, i, j
			}
		}
	}
	return bestScore, besti, bestj
}

func run(input string) {
	treemap := makeMap(input)
	fmt.Printf("%d visible trees\n", countVisible(treemap))
	score, x, y := bestScenicScore(treemap)
	fmt.Printf("Best score of %d at [%d,%d]\n", score, x, y)
}
