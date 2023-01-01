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
	hike                  []point
	debugX, debugY        = 41, 15
	depthDebug            = 30
	pause                 = 10 * time.Millisecond
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

	zstart, _ := topomap[start]
	pos := point{x: start.x, y: start.y, z: zstart}

	visit(pos)

	for _, v := range discovery {
		weights[position{x: v.x, y: v.y}] = v.weight
	}

	printSeen("./seen.txt")
	printWeights("")

	// printHike(hike, point{})
	// println("STEPS: ", len(hike))
}

func visit(pos point) {
	if destinationFucked {
		return
	}
	seen[position{x: pos.x, y: pos.y}] = true
	// println("VISITING ", pos.x, pos.y)

	if *debug && pos.weight > depthDebug {
		// time.Sleep(pause)
		// printMap(pos, "map.txt")
		// printMap(pos, "")
	}

	if pos.x == dest.x && pos.y == dest.y {
		if *verbose || *debug {
			println("Found dest ", dest.x, dest.y)
		}
		destinationFucked = true
		return
	}

	nb := neighbours(pos)

	noClimbing(pos.z, nb)
	noFalling(pos.z, nb)
	checkSeen(nb)

	newEntriesIndex := len(discovery)
	for k, _ := range nb {
		k.weight = pos.weight + 1
		discovery = append(discovery, k)
	}

	for _, v := range discovery[newEntriesIndex:] {
		visit(v)
	}
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

func printWeights(outFile string) {
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
		var line string
		for x := startx; x < endx; x++ {
			if weight, ok := weights[position{x: x, y: y}]; ok {
				line += "	" + fmt.Sprintf("%d", weight)
				continue
			}
			line += "	" + fmt.Sprintf("%d", -1)
		}
		fmt.Fprintln(out, line)
		println()
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

/*
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
*/

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
		if k.z < currentZ {
			delete(neighbours, k)
		}
	}
	return
}
