package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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
	x         int
	y         int
	z         int
	next      int
	direction int
}

var (
	start                 position
	dest                  position
	topomap               = make(map[position]int)
	DOWN, UP, RIGHT, LEFT = 1, 2, 3, 4
	bottom                = int('a')
	top                   = int('z')
	width                 = 0
	height                = 0
	hike                  []point
	seen                  = make(map[point]bool)
)

func initMap() {
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
	for sc.Scan() {
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
}

func main() {
	flag.Parse()

	initMap()

	if *debug {
		// printMap(point{})
	}

	zstart, _ := topomap[start]
	pos := point{x: start.x, y: start.y, z: zstart}

	hike, err := visit(pos)
	if err != nil {
		println(err)
	}

	printHike(hike, point{})
	println("STEPS: ", len(hike))
}

/*
TODO:

func visit(apoint, from point, directionmaybe point, bestHikeSofar) (hike []point, err error) {
  - get neighbours
  - remove falls,
  - remove climbing,
  - remove backtracking,
  - sort the rest by best direction
    for 1,2,3 nb {
    if hike, err := visit(1, fromhere)
    if err != nil {
    if diag entre 1 et 2 occupée visit(2)
    if err != nil {
    if diag entre 2 et 3 occupée visit(3)
    }
    continue
    }
    if diag entre 1 et 2 occupée visit(2)
    if err != nil {
    ...
    }
    if diag entre 2 et 3 occupée visit(3)
    if err != nil {
    ...
    }
    }
    }
*/

// TODO: optimizations to fail faster
func visit(pos point) (hike []point, err error) {

	if pos.x == dest.x && pos.y == dest.y {
		pos.next = 'X'
		return []point{pos}, nil
	}

	nb := neighbours(pos)

	noClimbing(pos.z, nb)
	nb = noFalling(pos.z, nb)
	// nb = noBackTracking(prev, nb)
	notTwice(nb)

	if len(nb) == 0 {
		return nil, fmt.Errorf("dead end at %d, %d", pos.x, pos.y)
	}

	sorted := higher(pos.z, nb)
	sortByDirection(pos, nb, &sorted)

	fewestSteps := 0
	oneDirection := 0
	var bestHike []point
	for i, v := range sorted {
		posTried, ok := nb[v]
		if !ok {
			panic("UNEXPECTED UNFOUND POS")
		}
		hike, err := visit(posTried)
		seen[posTried] = true
		if err != nil {
			if *debug {
				println(fmt.Sprintf("route %d, starting from %d, %d, direction %d: %v", i, pos.x, pos.y, v, err))
			}
			continue
		}
		if len(hike) > 0 && len(hike) < fewestSteps {
			oneDirection = v
			fewestSteps = len(hike)
			bestHike = hike
		}
		// TODO: do not try all the routes
	}

	if *debug {
		println(fmt.Sprintf("best route segment from %d, %d: %d, in %d steps", pos.x, pos.y, oneDirection, fewestSteps))
	}

	return append([]point{pos}, bestHike...), nil
}

// not thread safe
func notTwice(nb map[int]point) {
	for k, v := range nb {
		if _, ok := seen[v]; ok {
			delete(nb, k)
		}
	}
}

/*
v..v<<<<
>v.vv<<^
.>vv>E^^
..v>>>^^
..>>>>>^
*/

func printPoints(p map[int]point) {
	for k, v := range p {
		println(string(dirToByte(k)), v.x, v.y)
	}
}

func printMap(center point) {
	startx, starty := 0, 0
	endx, endy := width, height
	if center.x != 0 && center.y != 0 {
		startx = center.x - 5
		starty = center.y - 5
		endx = center.x + 5
		endy = center.y + 5
		if startx < 0 {
			startx = 0
		}
		if starty < 0 {
			starty = 0
		}
	}
	for y := starty; y < endy; y++ {
		var line []byte
		for x := startx; x < endx; x++ {
			elevation, _ := topomap[position{x: x, y: y}]
			line = append(line, byte(elevation))
		}
		println(string(line))
	}
}

func dirToByte(dir int) byte {
	switch dir {
	case UP:
		return '^'
	case DOWN:
		return 'v'
	case LEFT:
		return '<'
	case RIGHT:
		return '>'
	case 'X':
		return 'X'
	}
	return 64
}

func printHike(hike []point, center point) {

	startx, starty := 0, 0
	endx, endy := width, height
	if center.x != 0 && center.y != 0 {
		startx = center.x - 5
		starty = center.y - 5
		endx = center.x + 5
		endy = center.y + 5
		if startx < 0 {
			startx = 0
		}
		if starty < 0 {
			starty = 0
		}
	}

	var out io.Writer
	if *demo || center.x != 0 && center.y != 0 {
		out = os.Stdout
	} else {
		var err error
		out, err = os.Create("./output.txt")
		if err != nil {
			panic(err)
		}
		defer out.(*os.File).Close()
	}
	hikeMap := make(map[position]int)
	for _, v := range hike {
		hikeMap[position{x: v.x, y: v.y}] = v.next
	}
	for y := starty; y < endy; y++ {
		var line []byte
		for x := startx; x < endx; x++ {
			dir, ok := hikeMap[position{x: x, y: y}]
			if ok {
				line = append(line, dirToByte(dir))
				continue
			}
			line = append(line, '.')

		}
		fmt.Fprintln(out, string(line))
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

// TODO: do not copy
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

// TODO: do not copy
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

func higher(currentZ int, neighbours map[int]point) []int {
	var elevated []int
	for k, v := range neighbours {
		if v.z > currentZ {
			elevated = append(elevated, k)
		}
	}
	return elevated
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

func appendIfExists(sorted *[]int, neighbours map[int]point, direction ...int) {
	for _, d := range direction {
		if _, ok := neighbours[d]; ok {
			*sorted = append(*sorted, d)
		}
	}
}

func sortByDirection(pos point, neighbours map[int]point, sorted *[]int) {
	if len(neighbours) == 1 {
		for k, _ := range neighbours {
			*sorted = append(*sorted, k)
			return
		}
	}

	lower := make(map[int]point)
	for k, v := range neighbours {
		lower[k] = v
	}
	// the higher elevation one(s) are already a priority, so we only sort the ones
	// that are left, on the same elevation.

	for _, v := range *sorted {
		delete(lower, v)
	}

	direction := direction(pos)
	_, okLEFT := lower[LEFT]
	_, okRIGHT := lower[RIGHT]
	_, okUP := lower[UP]
	_, okDOWN := lower[DOWN]
	X := math.Abs(float64(direction.x))
	Y := math.Abs(float64(direction.y))

	// no good direction available
	// go in the bad direction that moves us away the least

	if direction.x >= 0 && !okRIGHT && direction.y >= 0 && !okDOWN {
		if X < Y {
			appendIfExists(sorted, lower, LEFT, UP)
			return
		}
		appendIfExists(sorted, lower, UP, LEFT)
		return
	}

	if direction.x >= 0 && !okRIGHT && direction.y <= 0 && !okUP {
		if X < Y {
			appendIfExists(sorted, lower, LEFT, DOWN)
			return
		}
		appendIfExists(sorted, lower, DOWN, LEFT)
		return
	}

	if direction.x <= 0 && !okLEFT && direction.y <= 0 && !okUP {
		if X < Y {
			appendIfExists(sorted, lower, RIGHT, DOWN)
			return
		}
		appendIfExists(sorted, lower, DOWN, RIGHT)
		return
	}

	if direction.x <= 0 && !okLEFT && direction.y >= 0 && !okDOWN {
		if X < Y {
			appendIfExists(sorted, lower, RIGHT, UP)
			return
		}
		appendIfExists(sorted, lower, UP, RIGHT)
		return
	}

	// exactly one good direction available, prioritize it.

	if direction.x >= 0 && okRIGHT && (direction.y >= 0 && !okDOWN || direction.y <= 0 && !okUP) {
		appendIfExists(sorted, lower, RIGHT, UP, LEFT, DOWN)
		return
	}

	if direction.x <= 0 && okLEFT && (direction.y >= 0 && !okDOWN || direction.y <= 0 && !okUP) {
		appendIfExists(sorted, lower, LEFT, UP, RIGHT, DOWN)
		return
	}

	if direction.y >= 0 && okDOWN && (direction.x >= 0 && !okRIGHT || direction.x <= 0 && !okLEFT) {
		appendIfExists(sorted, lower, DOWN, LEFT, UP, RIGHT)
		return
	}

	if direction.y <= 0 && okUP && (direction.x >= 0 && !okRIGHT || direction.x <= 0 && !okLEFT) {
		appendIfExists(sorted, lower, UP, DOWN, LEFT, RIGHT)
		return
	}

	// two good directions available. Prioritize the one that gets us the closest.

	if X >= Y {
		if direction.x >= 0 && okRIGHT {
			appendIfExists(sorted, lower, RIGHT, UP, LEFT, DOWN)
			return
		}
		if direction.x <= 0 && okLEFT {
			appendIfExists(sorted, lower, LEFT, RIGHT, UP, DOWN)
			return
		}
	}
	if direction.y >= 0 && okDOWN {
		appendIfExists(sorted, lower, DOWN, LEFT, UP, RIGHT)
		return
	}
	if direction.y <= 0 && okUP {
		appendIfExists(sorted, lower, UP, DOWN, LEFT, RIGHT)
		return
	}

	panic("NO SORTING DONE")
	return
}
