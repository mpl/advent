package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
)

var debug = flag.Bool("debug", false, `debug mode`)
var demo = flag.Bool("demo", false, `use demo input`)

type position struct {
	x int
	y int
}

type point struct {
	x    int
	y    int
	z    int
	next int
}

var start position
var dest position
var topomap = make(map[position]int)
var DOWN, UP, RIGHT, LEFT = 1, 2, 3, 4
var bottom = 97
var top = 122

var hike []point

func main() {
	flag.Parse()
	var f *os.File
	var err error
	if *demo {
		f, err = os.Open("./input.txt.demo")
	} else {
		f, err = os.Open("./input.txt")
	}
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	width := 0
	height := 0
	debugCount := 0
	for sc.Scan() {
		if *debug {
			if debugCount > 34 {
				break
			}
		}
		debugCount++

		l := sc.Text()
		width = len(l)
		for k, v := range []byte(l) {
			if string(v) == "S" {
				start = position{x: k, y: height}
				topomap[position{x: k, y: height}] = bottom
				continue
			}
			if string(v) == "E" {
				dest = position{x: k, y: height}
				topomap[position{x: k, y: height}] = top
				continue
			}
			topomap[position{x: k, y: height}] = int(v)
		}
		height++
	}

	for y := 0; y < height; y++ {
		var line []byte
		for x := 0; x < width; x++ {
			elevation, _ := topomap[position{x: x, y: y}]
			line = append(line, byte(elevation))
		}
		println(string(line))
	}

	zstart, _ := topomap[start]
	pos := point{x: start.x, y: start.y, z: zstart}
	vector := direction(pos)
	prevPos := pos
	hike = append(hike, pos)
	steps := 0
	for {
		if steps >= 30 {
			if *debug {
				println("EMERGENCY BRAKING")
			}
			break
		}
		if pos.x == dest.x && pos.y == dest.y {
			break
		}

		z := pos.z
		nb := neighbours(pos)
		if *debug {
			println("STEP", steps, "1", "DIRECTIONS LEFT:", len(nb))
		}
		noClimbing(z, nb)
		if *debug {
			println("STEP", steps, "2", "DIRECTIONS LEFT:", len(nb))
		}
		nb = noFalling(z, nb)
		if *debug {
			println("STEP", steps, "3", "DIRECTIONS LEFT:", len(nb))
		}
		nb = noBackTracking(prevPos, nb)
		if *debug {
			println("STEP", steps, "4", "DIRECTIONS LEFT:", len(nb))
		}
		nb = noWrongVerticalWay(vector, nb)
		if *debug {
			println("STEP", steps, "5", "DIRECTIONS LEFT:", len(nb))
		}
		nb = noWrongHorizontalWay(vector, nb)
		if *debug {
			println("STEP", steps, "6", "DIRECTIONS LEFT:", len(nb))
		}
		nb = bestDirection(vector, nb)

		if *debug {
			println("STEP", steps, "DIRECTIONS LEFT:", len(nb))
		}
		if len(nb) != 1 {
			panic("NOPE")
		}

		prevPos = pos
		for _, v := range nb {
			pos = v
			break
		}
		hike = append(hike, pos)
		vector = direction(pos)
		steps++
	}

}

func direction(current point) position {
	return position{
		x: dest.x - current.x,
		y: dest.y - current.y,
	}
}

func neighbours(pos point) map[int]point {
	var points = make(map[int]point)
	down := position{x: pos.x, y: pos.y + 1}
	up := position{x: pos.x, y: pos.y - 1}
	right := position{x: pos.x + 1, y: pos.y}
	left := position{x: pos.x - 1, y: pos.y}
	downE, ok := topomap[down]
	if ok {
		points[DOWN] = point{x: down.x, y: down.y, z: downE}
	}
	upE, ok := topomap[up]
	if ok {
		points[UP] = point{x: up.x, y: up.y, z: upE}
	}
	rightE, ok := topomap[right]
	if ok {
		points[RIGHT] = point{x: right.x, y: right.y, z: rightE}
	}
	leftE, ok := topomap[left]
	if ok {
		points[LEFT] = point{x: left.x, y: left.y, z: leftE}
	}
	return points
}

func noClimbing(currentZ int, neighbours map[int]point) {
	// println(currentZ)
	for k, v := range neighbours {
		// println("VS", v.z)
		if v.z > currentZ+1 {
			delete(neighbours, k)
		}
	}
}

func noFalling(currentZ int, neighbours map[int]point) map[int]point {
	withoutFalls := make(map[int]point)
	for k, v := range neighbours {
		if v.z < currentZ {
			continue
		}
		withoutFalls[k] = v
	}
	return withoutFalls
}

func noBackTracking(previous point, neighbours map[int]point) map[int]point {
	withoutBackTracking := make(map[int]point)
	for k, v := range neighbours {
		if v.x == previous.x && v.y == previous.y {
			continue
		}
		withoutBackTracking[k] = v
	}
	return withoutBackTracking
}

func noWrongVerticalWay(direction position, neighbours map[int]point) map[int]point {
	without := make(map[int]point)
	_, okUP := neighbours[UP]
	_, okDOWN := neighbours[DOWN]
	if !okUP || !okDOWN {
		return neighbours
	}
	for k, v := range neighbours {
		without[k] = v
	}

	if direction.y > 0 {
		delete(without, UP)
		return without
	}
	delete(without, DOWN)
	return without
}

func noWrongHorizontalWay(direction position, neighbours map[int]point) map[int]point {
	without := make(map[int]point)
	_, okLEFT := neighbours[LEFT]
	_, okRIGHT := neighbours[RIGHT]
	if !okLEFT || !okRIGHT {
		return neighbours
	}
	for k, v := range neighbours {
		without[k] = v
	}

	if direction.x > 0 {
		delete(without, LEFT)
		return without
	}
	delete(without, RIGHT)
	return without
}

func bestDirection(direction position, neighbours map[int]point) map[int]point {
	without := make(map[int]point)
	for k, v := range neighbours {
		without[k] = v
	}
	_, okLEFT := neighbours[LEFT]
	_, okRIGHT := neighbours[RIGHT]
	_, okUP := neighbours[UP]
	_, okDOWN := neighbours[DOWN]
	X := math.Abs(float64(direction.x))
	Y := math.Abs(float64(direction.y))
	if X > Y {
		if okLEFT || okRIGHT {
			delete(without, UP)
			delete(without, DOWN)
		}
		return without
	}
	if okUP || okDOWN {
		delete(without, LEFT)
		delete(without, RIGHT)
	}
	return without
}
