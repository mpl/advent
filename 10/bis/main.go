package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

var debug = flag.Bool("debug", false, `debug mode`)

func cycleOnSprite(cycleCount, spriteCenter int) bool {
	return cycleCount == spriteCenter-1 ||
		cycleCount == spriteCenter ||
		cycleCount == spriteCenter+1
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
	register := 1
	pixelCount := 0
	var pixels []bool

	for sc.Scan() {
		if *debug {
			if debugCount > 20 {
				break
			}
		}
		debugCount++

		l := sc.Text()

		duration := 1
		add := 0
		if !strings.HasPrefix(l, "noop") {
			duration = 2
			parts := strings.Split(l, " ")
			add, _ = strconv.Atoi(parts[1])
		}

		for i := 0; i < duration; i++ {
			if cycleOnSprite(pixelCount, register) {
				pixels = append(pixels, true)
			} else {
				pixels = append(pixels, false)
			}
			if *debug {
				println(pixelCount, l, "REGISTER: [", register-1, " - ", register+1, "]", "->", register+add)
				printCRT(pixels)
			}
			if len(pixels)%40 == 0 {
				pixelCount = 0
			} else {
				pixelCount++
			}
		}

		register += add
	}
	printCRT(pixels)
	// ERCREPCJ
}

func printCRT(pixels []bool) {
	var line string
	idx := 0
	for _, v := range pixels {
		if v {
			line += "#"
		} else {
			line += "."
		}
		idx++
		if idx == 40 {
			line += "\n"
			idx = 0
		}
	}
	println(line)
}

/*
##..##..##..##..##..##..##..##..##..##..
###...###...###...###...###...###...###.
####....####....####....####....####....
#####.....#####.....#####.....#####.....
######......######......######......####
#######.......#######.......#######.....
*/
