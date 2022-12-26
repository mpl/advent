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
	x        int
	y        int
	z        int
	next     int
	decision int
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

	debugCount := 0
	for sc.Scan() {
		if *debug {
			if debugCount > 34 {
				// break
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

	if *debug {
		// printMap(point{})
	}

	zstart, _ := topomap[start]
	pos := point{x: start.x, y: start.y, z: zstart}
	vector := direction(pos)
	prevPos := pos
	hike = append(hike, pos)
	retry := false
	steps := 0
	for {
		if steps >= 10000 {
			println("EMERGENCY BRAKING")
			break
		}
		if pos.x == dest.x && pos.y == dest.y {
			pos.next = 'X'
			hike[steps] = pos
			break
		}

		z := pos.z
		canRetry := false
		nb := neighbours(pos)

		decision := 0

		noClimbing(z, nb)
		decision++

		nb = noFalling(z, nb)
		decision++

		nb = noBackTracking(prevPos, nb)
		decision++

		nb = goHigher(nb)
		decision++

		// nb = noWrongVerticalWay(vector, nb)
		decision++

		// nb = noWrongHorizontalWay(vector, nb)
		decision++

		if len(nb) > 1 {
			canRetry = true
		}
		nb = bestDirection(vector, nb, retry)
		retry = false
		decision++
		_ = decision

		if *debug {
			// println("STEP", steps, "DIRECTIONS LEFT:", len(nb))
		}
		if len(nb) != 1 {
			panic("NOPE")
		}

		prevPos = pos
		var next int
		for k, v := range nb {
			next = k
			pos = v
			break
		}
		prevPos.next = next
		hike[steps] = prevPos

		if *debug {
			if steps > 45 {
				println(steps, prevPos.x, prevPos.y, string(dirToByte(next)), pos.x, pos.y)
			}
		}
		if _, ok := seen[pos]; ok {
			if *debug {
				println(steps, "CYCLE DETECTED ON", pos.x, pos.y)
				printMap(point{x: pos.x, y: pos.y})
				printHike(hike, point{x: pos.x, y: pos.y})
			}
			if !canRetry {
				break
			}
			retry = true
		}
		seen[pos] = true

		hike = append(hike, pos)
		vector = direction(pos)
		steps++
	}

	printHike(hike, point{})
	println("STEPS: ", steps)

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

func goHigher(neighbours map[int]point) map[int]point {
	highest := 0
	without := make(map[int]point)
	for k, v := range neighbours {
		if v.z > highest {
			highest = v.z
		}
		without[k] = v
	}

	for k, v := range without {
		if v.z < highest {
			delete(without, k)
		}
	}
	return without
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

func bestDirection(direction position, neighbours map[int]point, retry bool) map[int]point {
	if len(neighbours) == 1 {
		return neighbours
	}

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

	// no good direction available
	// go in the bad direction that moves us away the least
	// TODO: take into account retry here too?

	if direction.x >= 0 && !okRIGHT && direction.y >= 0 && !okDOWN {
		if X < Y {
			delete(without, UP)
			return without
		}
		delete(without, LEFT)
		return without
	}

	if direction.x >= 0 && !okRIGHT && direction.y <= 0 && !okUP {
		if X < Y {
			delete(without, DOWN)
			return without
		}
		delete(without, LEFT)
		return without
	}

	if direction.x <= 0 && !okLEFT && direction.y <= 0 && !okUP {
		if X < Y {
			delete(without, DOWN)
			return without
		}
		delete(without, RIGHT)
		return without
	}

	if direction.x <= 0 && !okLEFT && direction.y >= 0 && !okDOWN {
		if X < Y {
			delete(without, UP)
			return without
		}
		delete(without, RIGHT)
		return without
	}

	// exactly one good direction available, take it.
	// TODO: take into account retry here too?

	if direction.x >= 0 && okRIGHT && (direction.y >= 0 && !okDOWN || direction.y <= 0 && !okUP) {
		delete(without, UP)
		delete(without, DOWN)
		delete(without, LEFT)
		return without
	}

	if direction.x <= 0 && okLEFT && (direction.y >= 0 && !okDOWN || direction.y <= 0 && !okUP) {
		delete(without, UP)
		delete(without, DOWN)
		delete(without, RIGHT)
		return without
	}

	if direction.y >= 0 && okDOWN && (direction.x >= 0 && !okRIGHT || direction.x <= 0 && !okLEFT) {
		delete(without, UP)
		delete(without, RIGHT)
		delete(without, LEFT)
		return without
	}

	if direction.y <= 0 && okUP && (direction.x >= 0 && !okRIGHT || direction.x <= 0 && !okLEFT) {
		delete(without, DOWN)
		delete(without, RIGHT)
		delete(without, LEFT)
		return without
	}

	// two good directions available. Take the one that gets us the closest.

	if X >= Y && !retry {
		if direction.x >= 0 && okRIGHT {
			delete(without, UP)
			delete(without, DOWN)
			delete(without, LEFT)
			return without
		}
		if direction.x <= 0 && okLEFT {
			delete(without, UP)
			delete(without, DOWN)
			delete(without, RIGHT)
			return without
		}
	}
	if direction.y >= 0 && okDOWN {
		delete(without, UP)
		delete(without, RIGHT)
		delete(without, LEFT)
		return without
	}
	if direction.y <= 0 && okUP {
		delete(without, DOWN)
		delete(without, RIGHT)
		delete(without, LEFT)
		return without
	}
	return nil
}
