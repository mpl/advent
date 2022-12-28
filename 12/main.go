package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

var debug = flag.Bool("debug", false, `debug mode`)
var verbose = flag.Bool("verbose", false, `verbose mode`)
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
	RIGHT, UP, LEFT, DOWN = 1, 2, 3, 4
	bottom                = int('a')
	top                   = int('z')
	width                 = 0
	height                = 0
	hike                  []point
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

	seen := make(map[point]bool)
	hike, err := visit(pos, seen, 0)
	if err != nil {
		println(err.Error())
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

func copySeen(seen map[point]bool) map[point]bool {
	cps := make(map[point]bool)
	for k, v := range seen {
		cps[k] = v
	}
	return cps
}

// TODO: optimizations to fail faster
func visit(pos point, seen map[point]bool, depth int) (hike []point, err error) {
	if depth == 35 {
		return nil, errors.New("CRITICAL DEPTH")
	}
	seen[pos] = true
	if pos.x == dest.x && pos.y == dest.y {
		pos.next = 'X'
		return []point{pos}, nil
	}

	nb := neighbours(pos)

	noClimbing(pos.z, nb)
	nb = noFalling(pos.z, nb)
	// nb = noBackTracking(prev, nb)
	notTwice(nb, seen)

	if *debug {
		println(depth, pos.x, pos.y, "NEIGHBOURS: ", len(nb))
	}

	if len(nb) == 0 {
		return nil, fmt.Errorf("dead end at %d, %d", pos.x, pos.y)
	}

	above := higher(pos.z, nb)
	if *debug {
		msh, _ := json.Marshal(above)
		msh2, _ := json.Marshal(nb)
		println(depth, pos.x, pos.y, "ABOVE: ", string(msh), "BELOW: ", string(msh2))
	}
	sorted := append(sortByDirection(pos, above), sortByDirection(pos, nb)...)
	if *debug {
		msh, _ := json.Marshal(sorted)
		println(depth, pos.x, pos.y, "SORTED: ", string(msh))
	}

	// TODO: do better
	fewestSteps := 1000000000
	oneDirection := 0
	var bestHike []point
	for i, v := range sorted {
		posTried, ok := nb[v]
		if !ok {
			posTried, ok = above[v]
			if !ok {
				panic("UNEXPECTED UNFOUND POS")
			}
		}
		cps := copySeen(seen)
		hike, err := visit(posTried, cps, depth+1)
		if err != nil {
			if *debug {
				_ = i
				println(fmt.Sprintf("route %d, starting from %d, %d, direction %d: %v", i, pos.x, pos.y, v, err))
			}
			continue
		}
		if len(hike) > 0 && len(hike) < fewestSteps {
			oneDirection = v
			fewestSteps = len(hike)
			bestHike = hike
		}
		if (*debug || *verbose) && depth == 0 {
			printHike(hike, point{})
			println("IN ", len(hike), " STEPS")
		}
		// TODO: do not try all the routes
		// break
	}

	// TODO: how to still printHike to dead end, for debugging?

	if bestHike == nil {
		return nil, fmt.Errorf("dead end 2 at %d, %d", pos.x, pos.y)
	}

	if *debug {
		println(fmt.Sprintf("best route segment from %d, %d: %d, in %d steps", pos.x, pos.y, oneDirection, fewestSteps))
	}

	pos.next = oneDirection
	return append([]point{pos}, bestHike...), nil
}

// not thread safe
func notTwice(nb map[int]point, seen map[point]bool) {
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

func higher(currentZ int, neighbours map[int]point) map[int]point {
	above := make(map[int]point)
	for k, v := range neighbours {
		if v.z > currentZ {
			above[k] = v
			delete(neighbours, k)
		}
	}
	return above
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

func appendIfExists(neighbours map[int]point, direction ...int) []int {
	var sorted []int
	for _, d := range direction {
		if _, ok := neighbours[d]; ok {
			sorted = append(sorted, d)
		}
	}
	return sorted
}

func sortByDirection(pos point, nb map[int]point) []int {
	var sorted []int
	if len(nb) == 1 {
		for k, _ := range nb {
			return append(sorted, k)
		}
	}

	direction := direction(pos)
	_, okLEFT := nb[LEFT]
	_, okRIGHT := nb[RIGHT]
	_, okUP := nb[UP]
	_, okDOWN := nb[DOWN]
	X := math.Abs(float64(direction.x))
	Y := math.Abs(float64(direction.y))

	// no good direction available
	// go in the bad direction that moves us away the least

	if direction.x >= 0 && !okRIGHT && direction.y >= 0 && !okDOWN {
		if X < Y {
			return appendIfExists(nb, LEFT, UP)

		}
		return appendIfExists(nb, UP, LEFT)

	}

	if direction.x >= 0 && !okRIGHT && direction.y <= 0 && !okUP {
		if X < Y {
			return appendIfExists(nb, LEFT, DOWN)

		}
		return appendIfExists(nb, DOWN, LEFT)

	}

	if direction.x <= 0 && !okLEFT && direction.y <= 0 && !okUP {
		if X < Y {
			return appendIfExists(nb, RIGHT, DOWN)

		}
		return appendIfExists(nb, DOWN, RIGHT)

	}

	if direction.x <= 0 && !okLEFT && direction.y >= 0 && !okDOWN {
		if X < Y {
			return appendIfExists(nb, RIGHT, UP)

		}
		return appendIfExists(nb, UP, RIGHT)

	}

	// exactly one good direction available, prioritize it.

	if direction.x >= 0 && okRIGHT && (direction.y >= 0 && !okDOWN || direction.y <= 0 && !okUP) {
		return appendIfExists(nb, RIGHT, UP, LEFT, DOWN)

	}

	if direction.x <= 0 && okLEFT && (direction.y >= 0 && !okDOWN || direction.y <= 0 && !okUP) {
		return appendIfExists(nb, LEFT, UP, RIGHT, DOWN)

	}

	if direction.y >= 0 && okDOWN && (direction.x >= 0 && !okRIGHT || direction.x <= 0 && !okLEFT) {
		return appendIfExists(nb, DOWN, LEFT, UP, RIGHT)

	}

	if direction.y <= 0 && okUP && (direction.x >= 0 && !okRIGHT || direction.x <= 0 && !okLEFT) {
		return appendIfExists(nb, UP, DOWN, LEFT, RIGHT)

	}

	// two good directions available. Prioritize the one that gets us the closest.

	if X >= Y {
		if direction.x >= 0 && okRIGHT {
			return appendIfExists(nb, RIGHT, UP, LEFT, DOWN)

		}
		if direction.x <= 0 && okLEFT {
			return appendIfExists(nb, LEFT, RIGHT, UP, DOWN)

		}
	}
	if direction.y >= 0 && okDOWN {
		return appendIfExists(nb, DOWN, LEFT, UP, RIGHT)

	}
	if direction.y <= 0 && okUP {
		return appendIfExists(nb, UP, DOWN, LEFT, RIGHT)

	}

	panic("NO SORTING DONE")

}
