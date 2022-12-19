package main

import (
	"flag"
	"log"
	"os"
	"bufio"
	"strconv"
)

var debug = flag.Bool("debug", false, `debug mode`)

type tree struct {
	height int
	top bool
	bottom bool
	left bool
	right bool
}

func (t tree) visible() bool {
	return t.top || t.bottom || t.left || t.right
}

var forest [][]*tree

func main() {
	flag.Parse()
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	debugCount := 0
	lineIdx := 0
	forestWidth, forestHeight := 0, 0
	maxTopByCol := make(map[int]int)
	for sc.Scan() {
		forestHeight++
		if *debug {
			if debugCount > 5 {
				break
			}
			debugCount++
		}
		l := sc.Text()
		var line []*tree
		forestWidth = len(l)
		maxLeft := 0
		for k,v := range l {
			height, _ := strconv.Atoi(string(v))
			tr := tree{height: height}
			
			if lineIdx == 0 {
				tr.top = true
			}
			if k == 0 {
				tr.left = true
			}
			if k == len(l) - 1 {
				tr.right = true
			}
			// we can't do bottom row now because we don't know height of forest

			if height > maxLeft {
				tr.left = true
				maxLeft = height
			}
			maxTopInCol, _ := maxTopByCol[k]
			if height > maxTopInCol {
				maxTopByCol[k] = height
				tr.top = true
			}
			line = append(line, &tr)
		}
		forest = append(forest, line)
		lineIdx++
	}

	visibleCount := 0
	maxBottomByCol := make(map[int]int)
	for y:= forestHeight-1; y>=0; y-- {
		maxRight := 0
		for x:= forestWidth-1; x>=0; x-- {
			tr := forest[y][x]
			if y == forestHeight-1 {
				tr.bottom = true
			}
			if tr.height > maxRight {
				tr.right = true
				maxRight = tr.height
			}
			maxBottomInCol, _ := maxBottomByCol[x]
			if tr.height > maxBottomInCol {
				maxBottomByCol[x] = tr.height
				tr.bottom = true
			}
			if tr.visible() {
				visibleCount++
			}
		}
	}

	println("FOUND ", visibleCount, " / ", len(forest) * len(forest[0]))

}
