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

func main() {
	flag.Parse()
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	debugCount := 0
	register := 1
	cycleCount := 0
	cycles := []int{20, 60, 100, 140, 180, 220}
	cycleIdx := 0
	nextCycle := cycles[cycleIdx]
	signalStrengths := 0

	for sc.Scan() {
		if *debug {
			if debugCount > 34 {
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
		cycleCount += duration

		if *debug {
			println(debugCount, l, "after cycle", cycleCount, "register AT", register+add)
		}

		if cycleCount < nextCycle {
			register += add
			continue
		}

		println("CYCLECOUNT == ", cycleCount)
		signalStrength := register * nextCycle
		signalStrengths += signalStrength
		println("FOR CYCLE", nextCycle, "REGISTER AT", register, signalStrength)
		register += add
		cycleIdx++
		if cycleIdx == len(cycles) {
			break
		}
		nextCycle = cycles[cycleIdx]

	}
	println("FULL STRENGTH: ", signalStrengths)
}
