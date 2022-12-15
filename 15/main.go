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

func isBeacon(beacons []Pt, pt Pt) bool {
	for _, b := range beacons {
		if b == pt {
			return true
		}
	}
	return false
}

func run(input string) {
	readings := strings.Split(input, "\n")
	sensors, min, max := getSensors(readings)

	check := 0

	switch os.Args[1] {
	case "test":
		check = 10
	case "input":
		check = 2000000
	}

	beacons := make([]Pt, 0)
	for _, s := range sensors {
		if s.nearestBeacon.y == check {
			beacons = append(beacons, s.nearestBeacon)
		}
	}

	count := 0
	pt := Pt{min.x, check}
	for ; pt.x <= max.x; pt.x++ {
		foundSensor := false
		for i := 0; i < len(sensors) && !foundSensor; i++ {
			if pointInSensorRange(sensors[i], pt) {
				if !isBeacon(beacons, pt) {
					count++
					foundSensor = true
				}
			}
		}
	}
	fmt.Printf("Total points in range of sensor at y=%d is %d\n", check, count)
}
