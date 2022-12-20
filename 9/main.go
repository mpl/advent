package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var debug = flag.Bool("debug", false, `debug mode`)

type knot struct {
	x int
	y int
}

func (k *knot) move(direction string, distance int) {
	switch direction {
	case "R":
		k.x += distance
	case "L":
		k.x -= distance
	case "U":
		k.y += distance
	case "D":
		k.y -= distance
	}
}

var positions map[knot]bool = map[knot]bool{
	{0, 0}: true,
}

var head, tail = knot{}, knot{}

func main() {
	flag.Parse()
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	debugCount := 0
	for sc.Scan() {
		if *debug {
			if debugCount > 20 {
				break
			}
			debugCount++
		}
		l := sc.Text()
		parts := strings.Split(l, " ")
		direction := parts[0]
		distance, _ := strconv.Atoi(parts[1])
		head.move(direction, distance)
		if *debug {
			println("*******")
			println(l, " -> ", head.x, head.y)
		}
		for {
			positions[tail] = true
			if *debug {
				println(tail.x, tail.y)
			}
			if movedDiag := moveDiag(); movedDiag {
				continue
			}
			if movedX := moveX(); movedX {
				continue
			}
			if movedY := moveY(); movedY {
				continue
			}
			break
		}
	}
	println("POSITIONS: ", len(positions))

	/*
	   ..##..
	   ...##.
	   .####.
	   ....#.
	   s###..
	*/

	if *debug {
		for k, _ := range positions {
			println(k.x, k.y)
		}
	}
	// 5532 too low
}

func moveDiag() (moved bool) {
	if head.x == tail.x || head.y == tail.y {
		return false
	}
	if math.Abs(float64(head.x-tail.x)) < 2 && math.Abs(float64(head.y-tail.y)) < 2 {
		return false
	}
	movedX := moveDiagX()
	movedY := moveDiagY()
	return movedX || movedY
}

func moveDiagX() (moved bool) {
	if head.x > tail.x {
		tail.x++
		return true
	}
	if head.x < tail.x {
		tail.x--
		return true
	}
	return false
}

func moveDiagY() (moved bool) {
	if head.y > tail.y {
		tail.y++
		return true
	}
	if head.y < tail.y {
		tail.y--
		return true
	}
	return false
}

func moveX() (moved bool) {
	if head.x > tail.x+1 {
		tail.x++
		return true
	}
	if head.x < tail.x-1 {
		tail.x--
		return true
	}
	return false
}

func moveY() (moved bool) {
	if head.y > tail.y+1 {
		tail.y++
		return true
	}
	if head.y < tail.y-1 {
		tail.y--
		return true
	}
	return false
}
