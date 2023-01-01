package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	debug   = flag.Bool("debug", false, `debug mode`)
	verbose = flag.Bool("verbose", false, `verbose mode`)
	demo    = flag.Bool("demo", false, `use demo input`)
	destX   = flag.Int("destX", 0, `destx`)
	destY   = flag.Int("destY", 0, `desty`)
	startX  = flag.Int("startX", 0, `startx`)
	startY  = flag.Int("startY", 0, `starty`)
	gseen   = flag.Bool("seen", false, `seen so far`)
)

type position struct {
	x int
	y int
}

type point struct {
	x      int
	y      int
	z      int
	weight int
	next   int
}

var (
	start             position
	dest              position
	topomap           = make(map[position]int)
	seen              = make(map[position]bool)
	discovery         []point
	weights           = make(map[position]int)
	destinationFucked bool

	RIGHT, UP, LEFT, DOWN = 1, 2, 3, 4
	RU, UL, LD, DR        = 5, 6, 7, 8
	bottom                = int('a')
	top                   = int('z')
	width                 = 0
	height                = 0
	debugX, debugY        = 41, 15
	depthDebug            = 30
	pause                 = time.Second
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
			printSeen("./seen.txt")
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

	// zstart, _ := topomap[start]
	zdest, _ := topomap[dest]
	pos := point{x: dest.x, y: dest.y, z: zdest, weight: 0}

	discovery = append(discovery, pos)
	visitAll()

	for _, v := range discovery {
		weights[position{x: v.x, y: v.y}] = v.weight
	}

	hike := genHike()

	printSeen("./seen.txt")
	printWeights(dest, "./weights.txt")
	printHike(hike, point{})
	println("STEPS: ", len(hike))
}

func visitAll() {
	length := 1
	toVisit := discovery
	for {
		for _, v := range toVisit {
			visit(v)
		}
		index := length
		length = len(discovery)
		toVisit = discovery[index:]
		if destinationFucked {
			return
		}
	}
}

func visit(pos point) {
	if *debug {
		time.Sleep(pause)
		println("VISITING ", pos.x, pos.y)
		// printMap(pos, "map.txt")
		// printMap(pos, "")
	}

	if destinationFucked {
		return
	}
	seen[position{x: pos.x, y: pos.y}] = true
	if pos.x == start.x && pos.y == start.y {
		if *verbose || *debug {
			println("Found start ", start.x, start.y)
		}
		destinationFucked = true
		return
	}

	nb := neighbours(pos)
	noClimbing(pos.z, nb)
	noFalling(pos.z, nb)
	checkSeen(nb)
	if *debug {
		println(len(nb), " neighbours from", pos.x, pos.y)
	}

	for k, _ := range nb {
		k.weight = pos.weight + 1
		discovery = append(discovery, k)
	}
}

func genHike() []point {
	pos := point{x: start.x, y: start.y, weight: 10000}
	var hike []point
	for {
		if pos.x == dest.x && pos.y == dest.y {
			hike = append(hike, pos)
			break
		}
		nb := neighbours2(pos)
		bestDir := 0
		lowestWeight := pos.weight
		var bestNext point
		for k, v := range nb {
			if v.weight < lowestWeight {
				lowestWeight = v.weight
				bestDir = k
				bestNext = v
			}
		}
		pos.next = bestDir
		hike = append(hike, pos)
		pos = bestNext
	}
	return hike
}

func neighbours2(pos point) map[int]point {
	var points = make(map[int]point)
	down := position{x: pos.x, y: pos.y + 1}
	up := position{x: pos.x, y: pos.y - 1}
	right := position{x: pos.x + 1, y: pos.y}
	left := position{x: pos.x - 1, y: pos.y}
	wd, ok := weights[down]
	if ok {
		points[DOWN] = point{x: down.x, y: down.y, weight: wd}
	}
	wu, ok := weights[up]
	if ok {
		points[UP] = point{x: up.x, y: up.y, weight: wu}
	}
	wr, ok := weights[right]
	if ok {
		points[RIGHT] = point{x: right.x, y: right.y, weight: wr}
	}
	wl, ok := weights[left]
	if ok {
		points[LEFT] = point{x: left.x, y: left.y, weight: wl}
	}
	return points

}

func checkSeen(nb map[point]bool) {
	for k, _ := range nb {
		if _, ok := seen[position{x: k.x, y: k.y}]; ok {
			delete(nb, k)
		}
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

func printWeights(center position, outFile string) {
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
		startx = center.x - 10
		starty = center.y - 10
		endx = center.x + 10
		endy = center.y + 10
		if startx < 0 {
			startx = 0
		}
		if starty < 0 {
			starty = 0
		}
	}
	for y := starty; y < endy; y++ {
		var line string
		for x := startx; x < endx; x++ {
			if weight, ok := weights[position{x: x, y: y}]; ok {
				line += "	" + fmt.Sprintf("%d", weight)
				continue
			}
			line += "	" + fmt.Sprintf("%d", -1)
		}
		fmt.Fprintln(out, line)
		fmt.Fprintln(out)
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
			if _, ok := seen[position{x: x, y: y}]; ok {
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

func neighbours(pos point) map[point]bool {
	var points = make(map[point]bool)
	down := position{x: pos.x, y: pos.y + 1}
	up := position{x: pos.x, y: pos.y - 1}
	right := position{x: pos.x + 1, y: pos.y}
	left := position{x: pos.x - 1, y: pos.y}
	downE, ok := topomap[down]
	if ok {
		points[point{x: down.x, y: down.y, z: downE}] = true
	}
	upE, ok := topomap[up]
	if ok {
		points[point{x: up.x, y: up.y, z: upE}] = true
	}
	rightE, ok := topomap[right]
	if ok {
		points[point{x: right.x, y: right.y, z: rightE}] = true
	}
	leftE, ok := topomap[left]
	if ok {
		points[point{x: left.x, y: left.y, z: leftE}] = true
	}
	return points
}

func noClimbing(currentZ int, neighbours map[point]bool) {
	for k, _ := range neighbours {
		if k.z > currentZ+1 {
			delete(neighbours, k)
		}
	}
	return
}

func noFalling(currentZ int, neighbours map[point]bool) {
	for k, _ := range neighbours {
		if k.z < currentZ-1 {
			delete(neighbours, k)
		}
	}
	return
}
