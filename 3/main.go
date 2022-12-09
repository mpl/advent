package main

import (
	"bufio"
	"log"
	"os"
)

// A 65 -> 27
// a 97 -> 1

func priority(item rune) int {
	bVal := int(item)
	if bVal < 97 {
		return bVal - 38
	}
	return bVal - 96
}

var debug = false

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	debugCount := 0

	currentTotal := 0

	sc := bufio.NewScanner(f)
	for sc.Scan() {

		if debug {
			if debugCount > 5 {
				break
			}
			debugCount++
		}

		l := sc.Text()
		ll := len(l)

		firstHalf := make(map[rune]bool)

		halfPoint := ll / 2
		if debug {
			println(ll)
			println(halfPoint)
		}
		for k, v := range l {
			if k == halfPoint {
				break
			}
			firstHalf[v] = true
		}
		if debug {
			println(l[:halfPoint], l[halfPoint:])
		}
		for k, v := range l[halfPoint:] {
			if _, ok := firstHalf[v]; ok {
				if debug {
					println("DUPE AT: ", k, string(v), v, priority(v))
				}
				currentTotal += priority(v)
				break
			}
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	println(currentTotal)

}
