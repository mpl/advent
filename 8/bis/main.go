package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
)

var debug = flag.Bool("debug", false, `debug mode`)

type tree struct {
	height int
	top    int
	bottom int
	left   int
	right  int
}

var forest [][]*tree

/*
30373
25512
65332
33549
35390

313213123212200312011243203120214010202554420335045116203101005212525131015305511140012431022113113

*/

func main() {
	flag.Parse()
	f, err := os.Open("../input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	debugCount := 0
	lineIdx := 0
	forestWidth, forestHeight := 0, 0
	//	maxTopByCol := make(map[int]int)
	for sc.Scan() {
		forestHeight++
		if *debug {
			if debugCount > 20 {
				// break
			}
			debugCount++
		}
		l := sc.Text()

		// left to right pass
		var line []*tree
		forestWidth = len(l)
		closestHeights := make([]int, 10)
		for k, v := range l {
			height, _ := strconv.Atoi(string(v))
			tr := tree{height: height}

			closestBlocker := 0
			// TODO: maybe better algo: just look at the tree heights starting from our position.
			for _, v := range closestHeights[height:] {
				if v != 0 {
					if v > closestBlocker {
						closestBlocker = v
					}
				}
			}
			tr.left = k - closestBlocker
			closestHeights[height] = k
			line = append(line, &tr)
		}

		// right to left pass
		closestHeights = make([]int, 10)
		for i := forestWidth - 1; i >= 0; i-- {
			tr := line[i]
			height := tr.height
			closestBlocker := forestWidth - 1
			for _, v := range closestHeights[height:] {
				if v != 0 {
					if v < closestBlocker {
						closestBlocker = v
					}
				}
			}
			tr.right = closestBlocker - i
			closestHeights[height] = i
		}

		forest = append(forest, line)
		lineIdx++
	}

	bestScenic := 0
	scenicX, scenicY := 0, 0
	for x := 0; x < forestWidth; x++ {
		// top to bottom pass
		closestHeights := make([]int, 10)
		for y := 0; y < forestHeight; y++ {
			tr := forest[y][x]
			height := tr.height
			closestBlocker := 0
			for _, v := range closestHeights[height:] {
				if v != 0 {
					if v > closestBlocker {
						closestBlocker = v
					}
				}
			}
			tr.top = y - closestBlocker
			closestHeights[height] = y
		}

		// bottom to top pass
		closestHeights = make([]int, 10)
		for y := forestHeight - 1; y >= 0; y-- {
			tr := forest[y][x]
			height := tr.height
			closestBlocker := forestHeight - 1
			for _, v := range closestHeights[height:] {
				if v != 0 {
					if v < closestBlocker {
						closestBlocker = v
					}
				}
			}
			tr.bottom = closestBlocker - y
			closestHeights[height] = y
			scenic := tr.left * tr.top * tr.right * tr.bottom
			if scenic > bestScenic {
				bestScenic = scenic
				scenicX, scenicY = x, y
				if *debug {
					println("NEW BEST SCENIC: ", bestScenic, "(", tr.left, tr.top, tr.right, tr.bottom, ")", "at", scenicY, ",", scenicX, "with height", tr.height)
				}
			}
		}
	}

	println(forestWidth, forestHeight)

	bestTree := forest[scenicY][scenicX]
	println("BEST SCENIC: ", bestScenic, "(", bestTree.left, bestTree.top, bestTree.right, bestTree.bottom, ")", "at", scenicY, ",", scenicX, "with height", bestTree.height)


	return

	if *debug {

		for x := 0; x < forestWidth; x++ {
			for y := 0; y < forestHeight; y++ {
				tr := forest[y][x]
				println(y, tr.top, "^", tr.height, "v", tr.bottom)
			}
			break
		}

		return

		for k, v := range forest[0] {
			println(v.top)
			println(k, v.left, "<-", v.height, "->", v.right)
			println(v.bottom)
		}

	}

}
