package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	start := time.Now()
	text := readInput(inputName)
	run(text)
	end := time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
}

type ValveNode struct {
	name string
	flow int
}

// var valveMap map[string]Valve
var valveNodeMap map[string]ValveNode
var valveNeighbors map[string][]string

func (v ValveNode) Neighbors() ([]dijkstra.Node, []int) {
	neighbors := make([]dijkstra.Node, len(valveNeighbors[v.name]))
	distances := make([]int, len(valveNeighbors[v.name]))
	for i, n := range valveNeighbors[v.name] {
		neighbors[i] = dijkstra.Node(valveNodeMap[n])
		distances[i] = 1
	}
	return neighbors, distances
}

func parse(line string) (ValveNode, []string) {
	line = strings.TrimPrefix(line, "Valve ")
	name, line, _ := strings.Cut(line, " has flow rate=")
	flowString, line, _ := strings.Cut(line, "; ")
	flow, _ := strconv.Atoi(flowString)
	neighbors := strings.Split(line, " ")
	neighbors = neighbors[4:]
	for i, n := range neighbors {
		neighbors[i] = strings.TrimSuffix(n, ",")
	}
	return ValveNode{name, flow}, neighbors
}

func newTunnelMap(input string) (map[string]ValveNode, map[string][]string) {
	valves := strings.Split(input, "\n")
	valveMap := make(map[string]ValveNode)
	valveNeighbors := make(map[string][]string)
	for _, v := range valves {
		node, neighbors := parse(v)
		valveMap[node.name] = node
		valveNeighbors[node.name] = neighbors
	}
	return valveMap, valveNeighbors
}

func closeBestValve(openValves map[string]struct{}, currentValve string, valveGraph []dijkstra.Node, timeRemaining int) (string, int, int) {
	distances, _ := dijkstra.GetShortestDistances(valveGraph, valveNodeMap[currentValve])
	bestValue := -1
	bestCost := 0
	bestValve := ""
	for k, _ := range openValves {
		node := dijkstra.Node(valveNodeMap[k])
		distance := distances[node]
		cost := distance + 1
		value := (timeRemaining - cost) * valveNodeMap[k].flow
		if value > bestValue {
			bestValve = k
			bestValue = value
			bestCost = cost
		}
	}
	return bestValve, bestCost, bestValue
}

func getAllCosts(valveGraph []dijkstra.Node) map[string]map[dijkstra.Node]int {
	costs := make(map[string]map[dijkstra.Node]int)
	for _, start := range valveGraph {
		costsFrom, _ := dijkstra.GetShortestDistances(valveGraph, start)
		costs[start.(ValveNode).name] = costsFrom
	}
	return costs
}

type ComboResult struct {
	list  []string
	value int
}

var allCosts map[string]map[dijkstra.Node]int

func valuePath(path []string) int {
	totalValue := 0
	for i, timeRemaining := 1, 30; i < len(path)-1 && timeRemaining > 0; i++ {
		valveNode := valveNodeMap[path[i]]
		cost := allCosts[path[i-1]][valveNodeMap[path[i]]]
		if cost+1 > timeRemaining {
			timeRemaining = 0
			continue
		}
		value := valveNode.flow * timeRemaining
		timeRemaining -= cost + 1
		totalValue += value
	}
	return totalValue
}

func combinations(input []string, totalLength int, combo chan []string) {
	if length == 1 {
		result := make([]string, 1)
		result[1] = input[0]
		return result
	}
	for i := 0; i < len(input); i++ {
		first := input[i]
		remaining := append(input[:i], input[i+1:]...)
		combinations(remaining, totalLength, combo)
		if len(input) == totalLength {
			result := make([]string, totalLength)

		}
	}
}

func tryCombinations(valves []string, count int, start int, result []string, values chan []string) {
	tryCombinations2(valves, count, start, result, values)
	close(values)
}

func tryCombinations2(valves []string, count int, start int, result []string, values chan []string) {
	if count == 0 {
		fmt.Println(result)
		return
	}
	for i := start; i <= len(valves)-count; i++ {
		result[len(result)-count] = valves[i]
		tryCombinations2(valves, count-1, i+1, result, values)
	}
}

func run(input string) {
	valveNodeMap, valveNeighbors = newTunnelMap(input)
	valveGraph := make([]dijkstra.Node, 0)
	for _, v := range valveNodeMap {
		valveGraph = append(valveGraph, v)
	}
	fmt.Println(valveNodeMap)
	valves := make([]string, len(valveNodeMap))
	valves[0] = "AA"
	i := 1
	for k := range valveNodeMap {
		if k != "AA" {
			valves[i] = k
			i++
		}
	}

	results := make(chan []string, 5)
	combo := make([]string, 10)
	allCosts = getAllCosts(valveGraph)
	tryCombinations(valves, 10, 0, combo, results)

	// bestCombo := ComboResult{make([]string, 0), 0}
	// for comboResult := range results {
	// 	if comboResult.value > bestCombo.value {
	// 		bestCombo = comboResult
	// 		fmt.Printf("Best so far = %v\n", bestCombo)
	// 	}
	// }
}
