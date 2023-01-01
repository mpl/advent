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
	"sync"
	"time"
)

var (
	debug   = flag.Bool("debug", false, `debug mode`)
	verbose = flag.Bool("verbose", false, `verbose mode`)
	demo    = flag.Bool("demo", false, `use demo input`)
	destX   = flag.Int("destX", 0, `destx`)
	destY   = flag.Int("destY", 0, `desty`)
	gseen   = flag.Bool("seen", false, `seen so far`)
)

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
	RU, UL, LD, DR        = 5, 6, 7, 8
	bottom                = int('a')
	top                   = int('z')
	width                 = 0
	height                = 0
	hike                  []point
	debugX, debugY        = 41, 15
	depthDebug            = 30
	pause                 = 10 * time.Millisecond

	seenMu     sync.Mutex
	globalSeen = make(map[position]bool)
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

	if *gseen {
		go func() {
			time.Sleep(30 * time.Second)
			seenMu.Lock()
			printSeen("./seen.txt")
			seenMu.Unlock()
			log.Fatal("TIMEOUT")
		}()
	}

	println("STARTING AT ", start.x, start.y)
	if *destX != 0 {
		dest.x = *destX
	}
	if *destY != 0 {
		dest.y = *destY
	}

	zstart, _ := topomap[start]
	pos := point{x: start.x, y: start.y, z: zstart}

	seen := make(map[point]bool)
	hike, err := visit(pos, seen, 0)
	if err != nil {
		println(err.Error())
	}

	seenMu.Lock()
	printSeen("./seen.txt")
	seenMu.Unlock()
	printHike(hike, point{})
	println("STEPS: ", len(hike))
}

func copySeen(seen map[point]bool) map[point]bool {
	cps := make(map[point]bool)
	for k, v := range seen {
		cps[k] = v
	}
	return cps
}

// TODO: optimizations to fail faster
func visit(pos point, seen map[point]bool, depth int) (hike []point, err error) {
	// println("VISITING ", pos.x, pos.y)
	if depth == 10000 {
		// TODO: is this makeing us loop forever??
		// yep, ffs.
		return nil, errors.New("CRITICAL DEPTH")
	}

	if *debug && depth > depthDebug {
		// time.Sleep(pause)
		// printMap(pos, "map.txt")
		// printMap(pos, "")
	}

	seenMu.Lock()
	globalSeen[position{x: pos.x, y: pos.y}] = true
	seenMu.Unlock()

	seen[pos] = true
	if pos.x == dest.x && pos.y == dest.y {
		pos.next = 'X'
		if *verbose || *debug {
			println("Found dest ", dest.x, dest.y, "in ", depth, "steps")
		}
		return []point{pos}, nil
	}

	nb := neighbours(pos)

	obstacle := noClimbing(pos.z, nb)
	obstacle = append(obstacle, noFalling(pos.z, nb)...)
	notTwice(nb, seen)

	if *debug && (pos.x == debugX && pos.y == debugY || depth > depthDebug) {
		println(depth, pos.x, pos.y, "NEIGHBOURS: ", len(nb))
	}

	if len(nb) == 0 {
		return nil, fmt.Errorf("dead end at %d, %d", pos.x, pos.y)
	}

	above := higher(pos.z, nb)
	if *debug && (pos.x == debugX && pos.y == debugY || depth > depthDebug) {
		msh, _ := json.Marshal(above)
		msh2, _ := json.Marshal(nb)
		_ = msh
		_ = msh2
		println(depth, pos.x, pos.y, "ABOVE: ", string(msh), "BELOW: ", string(msh2))
	}

	// TODO: do we really always want to give priority to above first?
	sorted := append(sortByDirection(pos, above, obstacle), sortByDirection(pos, nb, obstacle)...)
	if *debug && (pos.x == debugX && pos.y == debugY || depth > depthDebug) {
		msh, _ := json.Marshal(sorted)
		println(depth, pos.x, pos.y, "SORTED: ", string(msh))
	}

	if pos.x == 47 && pos.y == 15 {
		// println("SORTEDDDDD: ", len(sorted))
	}

	// merge them again
	for k, v := range above {
		nb[k] = v
	}

	if len(sorted) == 0 {
		return nil, fmt.Errorf("dead way at %d, %d", pos.x, pos.y)
	}

	// TODO: remove?
	if isWrongWay(pos, sorted[0], obstacle) {
		if pos.x == 47 && pos.y == 15 {
			// println("SORTEDeuoueDDDD: ", len(sorted))
		}
		return nil, fmt.Errorf("wrong way from %d,%d going %d", pos.x, pos.y, sorted[0])
	}

	if depth > 0 && atEdge(pos) {

		if pos.x == 47 && pos.y == 15 {
			println("AYOOOOO: ", len(sorted))
		}

		return nil, fmt.Errorf("wrong way: at edge: %d,%d", pos.x, pos.y)
	}

	// TODO: do better
	fewestSteps := 1000000000
	oneDirection := 0
	var bestHike []point
	//	for i, v := range sorted {
	// not using range, because we might want to skip over one of them
	for i := 0; i < len(sorted); i++ {
		v := sorted[i]
		newPos, ok := nb[v]
		if !ok {
			panic("UNEXPECTED UNFOUND POS")
		}
		cps := copySeen(seen)
		hike, err := visit(newPos, cps, depth+1)
		if err != nil {
			if *debug {
				_ = i
				// println(fmt.Sprintf("route %d, starting from %d, %d, direction %d: %v", i, pos.x, pos.y, v, err))
			}
			continue
		}

		if len(hike) > 0 && len(hike) <= fewestSteps {
			// note as best hike so far
			oneDirection = v
			fewestSteps = len(hike)
			bestHike = hike
		}

		// still haven't found a route
		if oneDirection == 0 {
			continue
		}

		if i == len(sorted)-1 {
			break
		}

		if len(hike) == 1 {
			// we were the last stop, it would be stupid to try other routes
			// break
		}

		// should we still attempt other routes?

		if len(sorted) == 2 {
			newPos2, _ := nb[sorted[1]]

			if len(obstacle) > 0 {
				continue
			}

			if _, ok := isRubbing(newPos2, nb); ok {
				continue
			}

			break

			// println("FROM", pos.x, pos.y, "WE WOULD TRY", sorted[1], dir2.x, dir2.y)
			// break

			// we're touching an obstacle
			if len(obstacle) == 1 {
				continue
			}

			if _, ok := isRubbing(pos, nb); ok {
				// continue
			}

			break

			//			dir2, _ := nb[sorted[1]]
			if newPos.x == newPos2.x || newPos.y == newPos2.y {
				continue
			}
			if *debug {
				// println(fmt.Sprintf("%d IS BRANCHING2 FOR %d,%d and %d,%d", depth, dir1.x, dir1.y, dir2.x, dir2.y))
			}

			if isBranch(newPos, newPos2) {
				continue
			}

			break

			// TODO: remove all below?

			break
		}

		break

		/*
			if len(sorted) == 3 {
				log.Fatal("CAN IT HAPPEN?")
				dir2, _ := nb[sorted[i+1]]
				if *debug {
					// println(fmt.Sprintf("%d IS BRANCHING3 FOR %d,%d and %d,%d", depth, posTried.x, posTried.y, dir2.x, dir2.y))
				}
				if isBranch(posTried, dir2) {
					continue
				}
				if i > 0 {
					break
				}

				dir3, _ := nb[sorted[i+2]]
				if *debug {
					// println(fmt.Sprintf("%d IS BRANCHING3 BIS FOR %d,%d and %d,%d", depth, posTried.x, posTried.y, dir3.x, dir3.y))
				}
				if isBranch(posTried, dir3) {
					i++
					continue
				}
				break
			}

			if len(sorted) == 2 {
				dir1, _ := nb[sorted[0]]
				dir2, _ := nb[sorted[1]]
				if dir1.x == dir2.x || dir1.y == dir2.y {
					continue
				}
				if *debug {
					// println(fmt.Sprintf("%d IS BRANCHING2 FOR %d,%d and %d,%d", depth, dir1.x, dir1.y, dir2.x, dir2.y))
				}
				if isBranch(dir1, dir2) {
					continue
				}
				break
			}
		*/

		break

	}

	if bestHike == nil {
		return nil, fmt.Errorf("dead end 2 at %d, %d", pos.x, pos.y)
	}

	if *debug && (pos.x == debugX && pos.y == debugY || depth > depthDebug) {
		// println(fmt.Sprintf("best route segment from %d, %d: %d, in %d steps", pos.x, pos.y, oneDirection, fewestSteps))
	}

	pos.next = oneDirection
	return append([]point{pos}, bestHike...), nil
}

func isWrongWay(pos point, dir int, obstacle []int) bool {
	if len(obstacle) > 0 {
		return false
	}
	diffX := dest.x - pos.x
	diffY := dest.y - pos.y
	currentDist := diffX*diffX + diffY*diffY
	switch dir {
	case UP:
		pos.y--
	case DOWN:
		pos.y++
	case LEFT:
		pos.x--
	case RIGHT:
		pos.x++
	}
	diffX = dest.x - pos.x
	diffY = dest.y - pos.y
	newDist := diffX*diffX + diffY*diffY
	return newDist > currentDist
}

func atEdge(pos point) bool {
	if pos.x == 0 || pos.x == width-1 {
		return true
	}
	return pos.y == 0 || pos.y == height-1
}

func isBranch(dir1, dir2 point) bool {
	if dir1.x == dir2.x || dir1.y == dir2.y {
		return false
	}
	// we assume dir1 and dir2 are perpendicular
	dir1X := math.Abs(float64(dir1.x))
	dir1Y := math.Abs(float64(dir1.y))
	dir2X := math.Abs(float64(dir2.x))
	dir2Y := math.Abs(float64(dir2.y))
	if dir1.z != dir2.z {
		println(fmt.Sprintf("COMPARING Z of %d,%d and %d,%d", dir1.x, dir1.y, dir2.x, dir2.y))
	}
	var pos position
	if dir1X > dir2X {
		pos.x = dir1.x
	} else {
		pos.x = dir2.x
	}
	if dir1Y > dir2Y {
		pos.y = dir1.y
	} else {
		pos.y = dir2.y
	}
	z, ok := topomap[pos]
	if !ok {
		panic("AIAIAIAI")
	}
	// lower, or unclimbable == obstacle
	return (z < dir1.z || z > dir1.z+1) && (z < dir2.z || z > dir2.z+1)
}

// not thread safe
func notTwice(nb map[int]point, seen map[point]bool) {
	for k, v := range nb {
		if _, ok := seen[v]; ok {
			delete(nb, k)
		}
	}
}

func printPoints(p map[int]point) {
	for k, v := range p {
		println(string(dirToByte(k)), v.x, v.y)
	}
}

func printMap(center point, outFile string) {
	var out io.Writer
	if outFile == "" {
		out = os.Stdout
	} else {
		var err error
		filename := outFile
		out, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer out.(*os.File).Close()
	}
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
			if x == center.x && y == center.y {
				elevation = 'X'
			}
			line = append(line, byte(elevation))
		}
		fmt.Fprintln(out, string(line))
	}
}

func printSeen(outFile string) {
	var out io.Writer
	if outFile == "" {
		out = os.Stdout
	} else {
		var err error
		filename := outFile
		out, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer out.(*os.File).Close()
	}
	startx, starty := 0, 0
	endx, endy := width, height

	for y := starty; y < endy; y++ {
		var line []byte
		for x := startx; x < endx; x++ {
			if _, ok := globalSeen[position{x: x, y: y}]; ok {
				line = append(line, byte('X'))
				continue
			}
			elevation, _ := topomap[position{x: x, y: y}]
			line = append(line, byte(elevation))
		}
		fmt.Fprintln(out, string(line))
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
		filename := "./output.txt"
		out, err = os.Create(filename)
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

func noClimbing(currentZ int, neighbours map[int]point) []int {
	var obstacle []int
	for k, v := range neighbours {
		if v.z > currentZ+1 {
			obstacle = append(obstacle, k)
			delete(neighbours, k)
		}
	}
	return obstacle
}

func noFalling(currentZ int, neighbours map[int]point) []int {
	var obstacle []int
	for k, v := range neighbours {
		if v.z < currentZ {
			obstacle = append(obstacle, k)
			delete(neighbours, k)
		}
	}
	return obstacle
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

func appendIfExists(neighbours map[int]point, direction ...int) []int {
	var sorted []int
	for _, d := range direction {
		if _, ok := neighbours[d]; ok {
			sorted = append(sorted, d)
		}
	}
	return sorted
}

func isObstacle(pos, obs point) bool {
	return obs.z > pos.z+1 || obs.z < pos.z
}

func isRubbing(pos point, nb map[int]point) (int, bool) {
	// TODO: maybe isObstacle should compare z of pos to z of diag. not z of nb to z of diag.
	for k, v := range nb {
		nb2 := neighbours(v)
		switch k {
		case UP, DOWN:
			l, ok := nb2[LEFT]
			if ok && isObstacle(v, l) {
				if k == UP {
					return UL, true
				}
				return LD, true
			}
			r, ok := nb2[RIGHT]
			if ok && isObstacle(v, r) {
				if k == UP {
					return RU, true
				}
				return DR, true
			}
		case RIGHT, LEFT:
			u, ok := nb2[UP]
			if ok && isObstacle(v, u) {
				if k == RIGHT {
					return RU, true
				}
				return UL, true
			}
			d, ok := nb2[DOWN]
			if ok && isObstacle(v, d) {
				if k == RIGHT {
					return DR, true
				}
				return LD, true
			}

		}
	}
	return 0, false
}

func sortByDirection(pos point, nb map[int]point, obstacle []int) []int {
	var sorted []int
	if len(nb) == 0 {
		return sorted
	}
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

	// we're touching 2 obstacles, only one direction we can go
	// TODO: fuck, no, we could be coming from the hole. so 2 obstacles, and still possible 2 directions then.
	if len(obstacle) >= 2 {
		if len(nb) > 1 {
			// log.Fatalf("IMPOSSIBLE: at %d,%d: %d obstacles, and %d neighbours", pos.x, pos.y, len(obstacle), len(nb))
		}
	}

	if len(obstacle) == 1 {
		if pos.x == 40 && pos.y == 15 {
			// println("INDEED TOUCHING: ")
		}
		// not allowed to detach from obstacle if target is on the other side of obstacle
		switch obstacle[0] {
		case RIGHT:
			// target is on the other side obstacle, go along the obstacle in both directions
			if direction.x >= 0 {
				if direction.y > 0 {
					return appendIfExists(nb, DOWN, UP)
				}
				return appendIfExists(nb, UP, DOWN)
			}
			// target is on our side of the obstacle, we are allowed to detach
			if direction.y == 0 {
				return appendIfExists(nb, LEFT)
			}
			if direction.y > 0 {
				if X > Y {
					return appendIfExists(nb, LEFT, DOWN)
				}
				return appendIfExists(nb, DOWN, LEFT)
			}
			if X > Y {
				return appendIfExists(nb, LEFT, UP)
			}
			return appendIfExists(nb, UP, LEFT)
		case LEFT:
			if direction.x <= 0 {
				if direction.y > 0 {
					return appendIfExists(nb, DOWN, UP)
				}
				return appendIfExists(nb, UP, DOWN)
			}
			if direction.y == 0 {
				return appendIfExists(nb, RIGHT)
			}
			if direction.y > 0 {
				if X > Y {
					return appendIfExists(nb, RIGHT, DOWN)
				}
				return appendIfExists(nb, DOWN, RIGHT)
			}
			if X > Y {
				return appendIfExists(nb, RIGHT, UP)
			}
			return appendIfExists(nb, UP, RIGHT)
		case UP:
			if direction.y <= 0 {
				if direction.x > 0 {
					return appendIfExists(nb, RIGHT, LEFT)
				}
				return appendIfExists(nb, LEFT, RIGHT)
			}
			if direction.x == 0 {
				return appendIfExists(nb, DOWN)
			}
			if direction.x > 0 {
				if X > Y {
					return appendIfExists(nb, RIGHT, DOWN)
				}
				return appendIfExists(nb, DOWN, RIGHT)
			}
			if X > Y {
				return appendIfExists(nb, LEFT, DOWN)
			}
			return appendIfExists(nb, DOWN, LEFT)
		case DOWN:
			if direction.y >= 0 {
				if direction.x > 0 {
					return appendIfExists(nb, RIGHT, LEFT)
				}
				return appendIfExists(nb, LEFT, RIGHT)
			}
			if direction.x == 0 {
				return appendIfExists(nb, UP)
			}
			if direction.x > 0 {
				if X > Y {
					return appendIfExists(nb, RIGHT, UP)
				}
				return appendIfExists(nb, UP, RIGHT)
			}
			if X > Y {
				return appendIfExists(nb, LEFT, UP)
			}
			return appendIfExists(nb, UP, LEFT)
		}
		panic("WTF")
	}

	// TODO: what if we rub two different obstacles?
	diag, ok := isRubbing(pos, nb)
	if ok {
		if pos.x == 41 && pos.y == 15 {
			// println("INDEED RUBBING: ", diag)
		}
		if pos.x == 47 && pos.y == 15 {
			// println("INDEED WE RUBBING: ", diag)
		}
		// not allowed to detach from the obstacle,
		// unless if it's to go in the right direction
		switch diag {
		case RU:
			if direction.x == 0 {
				if direction.y < 0 {
					return appendIfExists(nb, UP, RIGHT)
				}
				return appendIfExists(nb, DOWN)
			}
			if direction.y == 0 {
				if direction.x > 0 {
					return appendIfExists(nb, RIGHT)
				}
				return appendIfExists(nb, LEFT)
			}
			return appendIfExists(nb, RIGHT, UP)
		case UL:
			if direction.x == 0 {
				if direction.y < 0 {
					return appendIfExists(nb, UP, LEFT)
				}
				return appendIfExists(nb, DOWN)
			}
			if direction.y == 0 {
				if direction.x > 0 {
					return appendIfExists(nb, RIGHT)
				}
				return appendIfExists(nb, LEFT)
			}
			return appendIfExists(nb, UP, LEFT)
		case LD:
			if pos.x == 41 && pos.y == 15 {
				// println("INDEED LD")
			}
			if direction.x == 0 {
				if direction.y > 0 {
					return appendIfExists(nb, DOWN, LEFT)
				}
				return appendIfExists(nb, UP)
			}
			if direction.y == 0 {
				if direction.x > 0 {
					// TODO: too strict?
					return appendIfExists(nb, RIGHT)
					// return appendIfExists(nb, RIGHT, DOWN)
				}
				return appendIfExists(nb, LEFT)
			}
			// TODO: compare X and Y ?
			return appendIfExists(nb, LEFT, DOWN)
		case DR:
			if pos.x == 47 && pos.y == 15 {
				// println("I CHOOSE YOU")
			}
			if direction.x == 0 {
				if direction.y > 0 {
					return appendIfExists(nb, DOWN, RIGHT)
				}
				return appendIfExists(nb, UP)
			}
			if direction.y == 0 {
				if direction.x > 0 {
					return appendIfExists(nb, RIGHT)
				}
				return appendIfExists(nb, LEFT)
			}
			// TODO: compare X and Y ?
			return appendIfExists(nb, RIGHT, DOWN)
		}
	}

	// no obstacle.
	// two good directions available. Prioritize the one that gets us the closest.

	if len(obstacle) == 0 {
		if X >= Y {
			if direction.x > 0 && okRIGHT {
				if direction.y == 0 {
					return appendIfExists(nb, RIGHT)
				}
				if direction.y > 0 && okDOWN {
					return appendIfExists(nb, RIGHT, DOWN)
				}
				return appendIfExists(nb, RIGHT, UP)
			}
			if direction.x < 0 && okLEFT {
				if direction.y == 0 {
					return appendIfExists(nb, LEFT)
				}
				if direction.y >= 0 && okDOWN {
					return appendIfExists(nb, LEFT, DOWN)
				}
				return appendIfExists(nb, LEFT, UP)
			}
		}

		if direction.y > 0 && okDOWN {
			if direction.x == 0 {
				return appendIfExists(nb, DOWN)
			}
			if direction.x > 0 && okRIGHT {
				return appendIfExists(nb, DOWN, RIGHT)
			}
			return appendIfExists(nb, DOWN, LEFT)
		}
		if direction.y < 0 && okUP {
			if direction.x == 0 {
				return appendIfExists(nb, UP)
			}
			if direction.x >= 0 && okRIGHT {
				return appendIfExists(nb, UP, RIGHT)
			}
			return appendIfExists(nb, UP, LEFT)
		}
	}

	// TODO: 2 obstacles

	return nil

	// TODO: clean up below

	// no good direction available.
	// try going along the obstacle, and not further from it.
	// otherwise just go in the bad direction that moves us the least possible.

	if direction.x >= 0 && !okRIGHT && direction.y >= 0 && !okDOWN {
		// in practice there's maximum one obstacle at this point
		//		if len(obstacle) > 0 && obstacle[0] == DOWN || X < Y {
		if len(obstacle) > 0 && obstacle[0] == DOWN {
			return appendIfExists(nb, LEFT, UP)
		}
		return appendIfExists(nb, UP, LEFT)

	}

	if direction.x >= 0 && !okRIGHT && direction.y <= 0 && !okUP {
		//		if len(obstacle) > 0 && obstacle[0] == UP || X < Y {
		if len(obstacle) > 0 && obstacle[0] == UP {
			return appendIfExists(nb, LEFT, DOWN)
		}
		return appendIfExists(nb, DOWN, LEFT)

	}

	if direction.x <= 0 && !okLEFT && direction.y <= 0 && !okUP {
		//		if len(obstacle) > 0 && obstacle[0] == UP || X < Y {
		if len(obstacle) > 0 && obstacle[0] == UP {
			return appendIfExists(nb, RIGHT, DOWN)
		}
		return appendIfExists(nb, DOWN, RIGHT)

	}

	if direction.x <= 0 && !okLEFT && direction.y >= 0 && !okDOWN {
		//		if len(obstacle) > 0 && obstacle[0] == DOWN || X < Y {
		if len(obstacle) > 0 && obstacle[0] == DOWN {
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
