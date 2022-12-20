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

var head = &knot{}

var knots []*knot = make([]*knot, 9)

func init() {
	for i := 0; i < 9; i++ {
		knots[i] = &knot{}
	}
}

func main() {
	flag.Parse()
	f, err := os.Open("../input.txt")
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
			leader := head
			idx := 0
			for k, v := range knots {
				idx = k
				follower := v
				if movedKnot := moveKnot(leader, follower); !movedKnot {
					break
				}
				if k == 8 {
					positions[*follower] = true
					if *debug {
						println("TAIL: ", follower.x, follower.y)
					}
				}
				leader = follower
			}
			// we did not even move the first knot, which means there's nothing to move anymore
			if idx == 0 {
				break
			}
		}
	}
	println("POSITIONS: ", len(positions))

	/*
	..........................
	..........................
	..........................
	..........................
	..........................
	..........................
	..........................
	..........................
	..........................
	#.........................
	#.............###.........
	#............#...#........
	.#..........#.....#.......
	..#..........#.....#......
	...#........#.......#.....
	....#......s.........#....
	.....#..............#.....
	......#............#......
	.......#..........#.......
	........#........#........
	.........########.........
	*/

	return

	if *debug {
		for k, _ := range positions {
			println(k.x, k.y)
		}
	}
	// 6524 too high
	// 2244 too low
}

func moveKnot(leader, follower *knot) (moved bool) {
	if movedDiag := moveDiag(leader, follower); movedDiag {
		return true
	}
	if movedX := moveX(leader, follower); movedX {
		return true
	}
	if movedY := moveY(leader, follower); movedY {
		return true
	}
	return false
}

func moveDiag(head, tail *knot) (moved bool) {
	if head.x == tail.x || head.y == tail.y {
		return false
	}
	if math.Abs(float64(head.x-tail.x)) < 2 && math.Abs(float64(head.y-tail.y)) < 2 {
		return false
	}
	movedX := moveDiagX(head, tail)
	movedY := moveDiagY(head, tail)
	return movedX || movedY
}

func moveDiagX(head, tail *knot) (moved bool) {
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

func moveDiagY(head, tail *knot) (moved bool) {
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

func moveX(head, tail *knot) (moved bool) {
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

func moveY(head, tail *knot) (moved bool) {
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
