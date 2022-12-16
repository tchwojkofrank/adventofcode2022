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

type Pt struct {
	x int
	y int
}

type Sensor struct {
	location      Pt
	nearestBeacon Pt
	radius        int
}

func newPt(info string) Pt {
	var pt Pt
	info = strings.TrimPrefix(info, "x=")
	coordinates := strings.Split(info, ", y=")
	pt.x, _ = strconv.Atoi(coordinates[0])
	pt.y, _ = strconv.Atoi(coordinates[1])
	return pt
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func distance(a, b Pt) int {
	return abs(b.x-a.x) + abs(b.y-a.y)
}

func getSensors(readings []string) ([]Sensor, Pt, Pt) {
	min := Pt{0, 0}
	max := Pt{0, 0}
	sensors := make([]Sensor, len(readings))
	for i, r := range readings {
		r = strings.TrimPrefix(r, "Sensor at ")
		info := strings.Split(r, ": closest beacon is at ")
		sensors[i].location = newPt(info[0])
		sensors[i].nearestBeacon = newPt(info[1])
		sensors[i].radius = distance(sensors[i].location, sensors[i].nearestBeacon)
		if (sensors[i].location.x - sensors[i].radius) < min.x {
			min.x = sensors[i].location.x - sensors[i].radius
		}
		if (sensors[i].location.y - sensors[i].radius) < min.y {
			min.y = sensors[i].location.y - sensors[i].radius
		}
		if (sensors[i].location.x + sensors[i].radius) > max.x {
			max.x = sensors[i].location.x + sensors[i].radius
		}
		if (sensors[i].location.y + sensors[i].radius) > min.y {
			max.y = sensors[i].location.y + sensors[i].radius
		}
	}
	return sensors, min, max
}

func pointInSensorRange(sensor Sensor, pt Pt) bool {
	return distance(sensor.location, pt) <= sensor.radius
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

var check int
var limit int

func run(input string) {
	readings := strings.Split(input, "\n")
	sensors, min, max := getSensors(readings)

	switch os.Args[1] {
	case "test":
		check = 10
		limit = 20
	case "input":
		check = 2000000
		limit = 4000000
	}

	beacons := make(map[Pt]struct{}, 0)
	for _, s := range sensors {
		if s.nearestBeacon.y == check {
			beacons[s.nearestBeacon] = struct{}{}
		}
	}

	count := 0
	pt := Pt{min.x, check}
	for ; pt.x <= max.x; pt.x++ {
		foundSensor := false
		for i := 0; i < len(sensors) && !foundSensor; i++ {
			if pointInSensorRange(sensors[i], pt) {
				if _, ok := beacons[pt]; !ok {
					count++
					foundSensor = true
				}
			}
		}
	}
	fmt.Printf("Total points in range of sensor at y=%d is %d\n", check, count)

	// part 2

	for y := 0; y <= limit; y++ {
		intervalslice := Intervals(make([]Pt, 0))
		intervals := &intervalslice
		for _, s := range sensors {
			intervals.addSensorInfo(s, y)
		}
		if len(*intervals) > 1 {
			fmt.Printf("Intervals on line %d\n%v\n", y, *intervals)
			x := (*intervals)[0].y + 1
			frequency := limit*x + y
			fmt.Printf("tuning frequency = %d\n", frequency)
			break
		}
	}
}

type Intervals []Pt

func (intervals *Intervals) addInterval(interval Pt) {

	// first interval whose max >= this interval's min
	first := 0
	for ; first < len(*intervals) && (*intervals)[first].y < interval.x; first++ {
	}
	if first >= len(*intervals) {
		*intervals = append(*intervals, interval)
		return
	}

	// last interval whose min >= this interval's max
	last := first
	for ; last < len(*intervals) && (*intervals)[last].x <= interval.y; last++ {
	}

	// no overlaps
	if first == last {
		// insert the interval before first

		*intervals = append((*intervals)[:first+1], (*intervals)[first:]...)
		(*intervals)[first] = interval
		return
	}

	// overlaps from first to last-1
	// create a new interval, first -> interval, delete from first +1 to last-1
	newInterval := Pt{min((*intervals)[first].x, interval.x), max((*intervals)[last-1].y, interval.y)}
	(*intervals)[first] = newInterval
	*intervals = append((*intervals)[:first+1], (*intervals)[last:]...)
}

func (intervals *Intervals) addSensorInfo(s Sensor, y int) {
	minY := s.location.y - s.radius
	maxY := s.location.y + s.radius
	if y >= minY && y <= maxY {
		dy := abs(y - s.location.y)
		dx := abs(s.radius - dy)
		interval := Pt{max(0, s.location.x-dx), min(limit, s.location.x+dx)}
		intervals.addInterval(interval)
	}
}
