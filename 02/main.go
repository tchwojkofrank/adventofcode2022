package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func readInput(fname string) string {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	return string(content)
}

const (
	rock     = 1
	paper    = 2
	scissors = 3
)

const (
	lose = 0
	draw = 3
	win  = 6
)

var playMap = map[string]int{
	"A": rock,
	"B": paper,
	"C": scissors,
	"X": rock,
	"Y": paper,
	"Z": scissors,
}

type Play struct {
	opponent int
	self     int
}

var scoreMap = map[Play]int{
	{1, 1}: 1 + 3,
	{1, 2}: 2 + 6,
	{1, 3}: 3 + 0,
	{2, 1}: 1 + 0,
	{2, 2}: 2 + 3,
	{2, 3}: 3 + 6,
	{3, 1}: 1 + 6,
	{3, 2}: 2 + 0,
	{3, 3}: 3 + 3,
}

// lose = 0, draw = 3, win = 6

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Enter input file name.\n")
		return
	}
	params := os.Args[1]
	inputName := strings.Split(params, " ")[0]
	fmt.Printf("inputName = %v\n", inputName)
	text := readInput(inputName)
	run(text)
}

func getPlayStrings(input string) []string {
	return strings.Split(input, "\n")
}

func getPlays(playStrings []string) []Play {
	plays := make([]Play, len(playStrings))
	for i, ps := range playStrings {
		playStrings := strings.Split(ps, " ")
		o := playStrings[0]
		s := playStrings[1]
		plays[i].opponent = playMap[o]
		plays[i].self = playMap[s]
	}
	return plays
}

func getScores(plays []Play) int {
	sum := 0
	for _, play := range plays {
		sum += scoreMap[play]
	}
	return sum
}

var resultMap = map[string]int{
	"X": lose,
	"Y": draw,
	"Z": win,
}

type Round struct {
	opponent int
	result   int
}

var findPlayMap = map[Round]int{
	{rock, lose}:     scissors,
	{rock, draw}:     rock,
	{rock, win}:      paper,
	{paper, lose}:    rock,
	{paper, draw}:    paper,
	{paper, win}:     scissors,
	{scissors, lose}: paper,
	{scissors, draw}: scissors,
	{scissors, win}:  rock,
}

func roundScore(round Round) int {
	return round.result + findPlayMap[round]
}

func getRounds(playStrings []string) []Round {
	rounds := make([]Round, len(playStrings))
	for i, ps := range playStrings {
		pr := strings.Split(ps, " ")
		rounds[i].opponent = playMap[pr[0]]
		rounds[i].result = resultMap[pr[1]]
	}
	return rounds
}

func getTotalScore(rounds []Round) int {
	sum := 0
	for _, r := range rounds {
		sum += roundScore(r)
	}
	return sum
}

func run(input string) {
	playStrings := getPlayStrings(input)
	plays := getPlays(playStrings)
	fmt.Printf("game total: %v\n", getScores(plays))
	rounds := getRounds(playStrings)
	fmt.Printf("rounds: %v\n", rounds)
	fmt.Printf("real game total: %v\n", getTotalScore(rounds))
}
