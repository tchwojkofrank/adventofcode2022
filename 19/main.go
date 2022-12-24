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

var earliestGeode = 0

func (inv Inventory) String() string {
	result := fmt.Sprintf("Ore: %d Clay: %d Obs: %d Geo: %d\n", inv.ore, inv.clay, inv.obisidian, inv.geode)
	result = result + fmt.Sprintf("OrR: %d CR %d ObR %d GeR %d\n\n", inv.oreRobot, inv.clayRobot, inv.obsidianRobot, inv.geodeRobot)
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func canBuildGeodeRobotInTime(bp Blueprint, inv Inventory, timeLeft int) bool {
	if inv.geodeRobot > 0 {
		return true
	}
	oreRate := inv.oreRobot
	obsidianRate := inv.obsidianRobot
	oreNeed := bp.geodeRobot[0] - inv.ore
	obsidianNeed := bp.geodeRobot[1] - inv.obisidian
	if (timeLeft-1)*obsidianRate >= obsidianNeed && (timeLeft-1)*oreRate >= oreNeed {
		return true
	}
	// time to build obsidianRobot?
	if (timeLeft-1)*obsidianRate < obsidianNeed {
		clayRate := inv.clayRobot
		if clayRate < 1 && timeLeft >= 21 {
			// too early, we need to build a clay robot
			return true
		}
		oreNeed = bp.obsidianRobot[0]
		clayNeed := bp.obsidianRobot[1]
		var timeToObsidian int
		if clayRate > 0 && oreRate > 0 {
			timeToObsidian = max(clayNeed/clayRate, oreNeed/oreRate)
			timeLeftAfterBuild := timeLeft - timeToObsidian
			obsidianRate = inv.obsidianRobot + 1
			if (timeLeftAfterBuild-1)*obsidianRate >= obsidianNeed {
				return true
			}
		}
		// time to build clayRobot?
		if timeLeft >= 10 {
			return true
		}

		timeToClayRobot := make([]int, 3)
		timeToFinishOreRobots := make([]int, 3)
		// build 0 to 2 ore robots
		for i := 0; i < 3; i++ {
			if i == 0 {
				timeToFinishOreRobots[i] = 0
			} else {
				timeToFinishOreRobots[i] = bp.oreRobot/(oreRate+i-1) + timeToFinishOreRobots[i-1]
			}
			oreUsed := i * bp.oreRobot
			oreMade := inv.oreRobot * timeToFinishOreRobots[i]
			if i == 2 {
				oreMade += timeToFinishOreRobots[i] - timeToFinishOreRobots[i-1]
			}
			timeToClayRobot[i] = (bp.clayRobot - inv.ore + oreUsed - oreMade) / (oreRate + i)
			oreUsed += bp.clayRobot
			timeToBuildObsidianRobot := max(bp.obsidianRobot[0]-inv.ore+oreUsed-oreMade/(oreRate+i), (bp.obsidianRobot[1]-inv.clay)/(inv.clayRobot+1))
			if timeToClayRobot[i]+timeToBuildObsidianRobot+timeToFinishOreRobots[i] <= timeLeft-1 {
				return true
			}
		}
	}

	return false
}

var attempt int = 0

func runBlueprint(bp Blueprint, inventory Inventory, timeLeft int, actionList string) (Inventory, string) {
	if timeLeft == 0 {
		if attempt%1000000000 == 0 {
			// fmt.Printf("attempt %d best %d\n%v\n", attempt, earliestGeode, actionList)
		}
		attempt++
		return inventory, actionList
	}
	// check if it's possible to build a geode robot in time
	if !canBuildGeodeRobotInTime(bp, inventory, timeLeft) {
		if attempt%1000000000 == 0 {
			// fmt.Printf("attempt %d best %d\n%v\n", attempt, earliestGeode, actionList)
		}
		attempt++
		return inventory, actionList
	}
	newinventory := inventory
	newinventory.ore += inventory.oreRobot
	newinventory.clay += inventory.clayRobot
	newinventory.obisidian += inventory.obsidianRobot
	newinventory.geode += inventory.geodeRobot
	if newinventory.geode > 0 && timeLeft > earliestGeode {
		earliestGeode = timeLeft
	} else if timeLeft <= earliestGeode-2 && newinventory.geode == 0 {
		return inventory, actionList
	}

	bestInventory := inventory
	bestAction := actionList

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
	} else if inventory.ore >= bp.obsidianRobot[0] && inventory.clay >= bp.obsidianRobot[1] {
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
	} else {
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

		// try not building anything this time
		noBuildInventory, noBuildActionList := runBlueprint(bp, newinventory, timeLeft-1, actionList+"no build\n")
		if noBuildInventory.geode > bestInventory.geode {
			bestInventory = noBuildInventory
			bestAction = noBuildActionList
		}
	}

	return bestInventory, bestAction
}

func run(input string) {
	blueprints := strings.Split(input, "\n")
	qualitySum := 0
	for i, b := range blueprints {
		blueprint := newBlueprint(b)
		fmt.Println(blueprint)
		var inventory Inventory
		inventory.oreRobot = 1
		earliestGeode = 0
		attempt = 0
		result, actionList := runBlueprint(blueprint, inventory, timelimit, "")
		fmt.Printf("Blueprint %d Inventory: %v\n", blueprint.index, result)
		fmt.Printf("Best set of actions:\n%v\n", actionList)
		fmt.Printf("Earliest geode robot %d\n", earliestGeode)
		fmt.Printf("Quality number for blueprint %d is %d\n", i+1, (i+1)*result.geode)
		qualitySum += (i + 1) * result.geode
	}
	fmt.Printf("\nQuality total:\n%d\n", qualitySum)
}
