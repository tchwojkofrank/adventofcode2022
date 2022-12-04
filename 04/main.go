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

type Range struct {
	min int
	max int
}

type RangePair [2]Range

func getRange(rangeString string) Range {
	valueStrings := strings.Split(rangeString, "-")
	var r Range
	r.min, _ = strconv.Atoi(valueStrings[0])
	r.max, _ = strconv.Atoi(valueStrings[1])
	return r
}

func getRangePairs(input string) []RangePair {
	pairStrings := strings.Split(input, "\n")
	rangePairs := make([]RangePair, len(pairStrings))
	for i, ps := range pairStrings {
		rangeStrings := strings.Split(ps, ",")
		rangePairs[i][0] = getRange(rangeStrings[0])
		rangePairs[i][1] = getRange(rangeStrings[1])
	}
	return rangePairs
}

func isContaining(rp RangePair) bool {
	return (((rp[0].min <= rp[1].min) && (rp[0].max >= rp[1].max)) ||
		((rp[0].min >= rp[1].min) && (rp[0].max <= rp[1].max)))
}

func isOverlapping(rp RangePair) bool {
	return (rp[0].min >= rp[1].min && rp[0].min <= rp[1].max) || (rp[0].max >= rp[1].min && rp[0].max <= rp[1].max) ||
		(rp[1].min >= rp[0].min && rp[1].min <= rp[0].max) || (rp[1].max >= rp[0].min && rp[1].max <= rp[0].max)
}

func countCondition(rangePairs []RangePair, condition func(rp RangePair) bool) int {
	count := 0

	for _, rp := range rangePairs {

		if condition(rp) {
			count++
		}
	}

	return count
}

func run(input string) {
	rangePairs := getRangePairs(input)
	contains := countCondition(rangePairs, isContaining)
	fmt.Printf("Total contains = %d\n", contains)
	overlaps := countCondition(rangePairs, isOverlapping)
	fmt.Printf("Total overlaps = %d\n", overlaps)
}
