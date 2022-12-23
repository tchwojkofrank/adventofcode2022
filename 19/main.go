package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
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

type Blueprint struct {
	index         int
	oreRobot      int    // ore cost
	clayRobot     int    // ore cost
	obsidianRobot [2]int // ore and clay cost
	geodeRobot    [2]int // ore and obsidian cost
}

//Blueprint 4: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 20 clay. Each geode robot costs 2 ore and 8 obsidian.

func newBlueprint(input string) Blueprint {
	re := regexp.MustCompile(`Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.`)
	result := re.FindAllStringSubmatch(input, -1)
	var b Blueprint
	b.index, _ = strconv.Atoi(result[0][1])
	b.oreRobot, _ = strconv.Atoi(result[0][2])
	b.clayRobot, _ = strconv.Atoi(result[0][3])
	b.obsidianRobot[0], _ = strconv.Atoi(result[0][4])
	b.obsidianRobot[1], _ = strconv.Atoi(result[0][5])
	b.geodeRobot[0], _ = strconv.Atoi(result[0][6])
	b.geodeRobot[1], _ = strconv.Atoi(result[0][7])
	return b
}

const timelimit = 24

type Inventory struct {
	ore           int
	oreRobot      int
	clay          int
	clayRobot     int
	obisidian     int
	obsidianRobot int
	geode         int
	geodeRobot    int
}

var attempt int = 0

func (inv Inventory) String() string {
	result := fmt.Sprintf("Ore: %d Clay: %d Obs: %d Geo: %d\n", inv.ore, inv.clay, inv.obisidian, inv.geode)
	result = result + fmt.Sprintf("OrR: %d CR %d ObR %d GeR %d\n\n", inv.oreRobot, inv.clayRobot, inv.obsidianRobot, inv.geodeRobot)
	return result
}

func runBlueprint(bp Blueprint, inventory Inventory, timeLeft int, actionList string) (Inventory, string) {
	if timeLeft == 0 {
		if attempt%1000 == 0 {
			fmt.Printf("Attempt %d action list\n%v\n", attempt, actionList)
		}
		attempt++
		return inventory, actionList
	}
	newinventory := inventory
	newinventory.ore += inventory.oreRobot
	newinventory.clay += inventory.clayRobot
	newinventory.obisidian += inventory.obsidianRobot
	newinventory.geode += inventory.geodeRobot

	// try not building anything this time
	noBuild, noBuildActionList := runBlueprint(bp, newinventory, timeLeft-1, actionList+"no build\n")
	bestInventory := noBuild
	bestAction := noBuildActionList

	if inventory.ore >= bp.geodeRobot[0] && inventory.obisidian >= bp.geodeRobot[1] {
		geodeInventory := newinventory
		var geodeActionList string
		geodeInventory.ore -= bp.geodeRobot[0]
		geodeInventory.obisidian -= bp.geodeRobot[1]
		geodeInventory.geodeRobot++
		geodeInventory, geodeActionList = runBlueprint(bp, geodeInventory, timeLeft-1, actionList+"build geode robot\n")
		if geodeInventory.geode > bestInventory.geode {
			bestInventory = geodeInventory
			bestAction = geodeActionList
		}
	}

	if inventory.ore >= bp.obsidianRobot[0] && inventory.clay >= bp.obsidianRobot[1] {
		obsidianInventory := newinventory
		var obsidianActionList string
		obsidianInventory.ore -= bp.obsidianRobot[0]
		obsidianInventory.clay -= bp.obsidianRobot[1]
		obsidianInventory.obsidianRobot++
		obsidianInventory, obsidianActionList = runBlueprint(bp, obsidianInventory, timeLeft-1, actionList+"build obsidian robot\n")
		if obsidianInventory.geode > bestInventory.geode {
			bestInventory = obsidianInventory
			bestAction = obsidianActionList
		}
	}

	if inventory.ore >= bp.clayRobot {
		clayInventory := newinventory
		var clayActionList string
		clayInventory.ore -= bp.clayRobot
		clayInventory.clayRobot++
		clayInventory, clayActionList = runBlueprint(bp, clayInventory, timeLeft-1, actionList+"build clay robot\n")
		if clayInventory.geode > bestInventory.geode {
			bestInventory = clayInventory
			bestAction = clayActionList
		}
	}

	if inventory.ore >= bp.oreRobot {
		oreInventory := newinventory
		var oreActionList string
		oreInventory.ore -= bp.oreRobot
		oreInventory.oreRobot++
		oreInventory, oreActionList = runBlueprint(bp, oreInventory, timeLeft-1, actionList+"build ore robot\n")
		if oreInventory.geode > bestInventory.geode {
			bestInventory = oreInventory
			bestAction = oreActionList
		}
	}

	return bestInventory, bestAction
}

func run(input string) {
	blueprints := strings.Split(input, "\n")
	for _, b := range blueprints {
		blueprint := newBlueprint(b)
		fmt.Println(blueprint)
		var inventory Inventory
		inventory.oreRobot = 1
		result, actionList := runBlueprint(blueprint, inventory, timelimit, "")
		fmt.Printf("Blueprint %d Inventory: %v\n", blueprint.index, result)
		fmt.Printf("Best set of actions:\n%v\n", actionList)
	}

}
