package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// A X Rock
// B Y Paper
// C Z Scissors

var intByPick map[string]int = map[string]int{
	"A": 1,
	"B": 2,
	"C": 3,
	"X": 1,
	"Y": 2,
	"Z": 3,
}

func pointsForOutcome(them, you int) int {
	if you == them {
		return 3
	}
	if you == them+1 {
		return 6
	}
	if you == them-2 {
		return 6
	}

	return 0
}

func main() {
	f, err := os.Open("./input.txt")
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
		myPick := intByPick[parts[1]]

		currentScore = myPick + pointsForOutcome(theirPick, myPick)
		currentTotal += currentScore
		// println(parts[0], parts[1], currentScore)
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	println(currentTotal)
}
