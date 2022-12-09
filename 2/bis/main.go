package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// A Rock
// B Paper
// C Scissors

// X loss
// Y draw
// Z win

var intByPick map[string]int = map[string]int{
	"A": 1,
	"B": 2,
	"C": 3,
	"X": 1,
	"Y": 2,
	"Z": 3,
}

func points(them, outcome int) int {
	// draw
	if outcome == 2 {
		return 3 + them
	}

	// loss
	if outcome == 1 {
		pick := them - 1
		// "underflow"
		if pick == 0 {
			pick = 3
		}
		return 0 + pick
	}

	// win
	// outcome == 3
	pick := them + 1
	if pick == 4 {
		pick = 1
	}
	return 6 + pick
}

func main() {
	f, err := os.Open("../input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	debugCount := 0

	currentTotal, currentScore := 0, 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {

		if debugCount > 5 {
			// break
		}
		debugCount++

		l := sc.Text()
		parts := strings.SplitN(l, " ", 2)
		theirPick := intByPick[parts[0]]
		outcome := intByPick[parts[1]]

		currentScore = points(theirPick, outcome)
		currentTotal += currentScore
		// println(parts[0], parts[1], currentScore)
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	println(currentTotal)
}
