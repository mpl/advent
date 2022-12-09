package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func contains(l1, h1, l2, h2 int) bool {
	if l2 > h1 {
		return false
	}
	if l1 > h2 {
		return false
	}
	return true
	/*
	   	if l2 > l1 {
	   		return l2 <= h1
	   	}

	   	if l2 < l1 {
	   		return h2 >= l1
	   	}

	   return l2 == l1 || h2 == h1
	*/
}

var debug = false

func main() {
	f, err := os.Open("../input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	debugCount := 0

	currentTotal := 0

	sc := bufio.NewScanner(f)
	for sc.Scan() {

		if debug {
			if debugCount > 20 {
				break
			}
			debugCount++
		}

		l := sc.Text()

		parts := strings.Split(l, ",")
		one := strings.Split(parts[0], "-")
		two := strings.Split(parts[1], "-")
		l1, _ := strconv.Atoi(one[0])
		h1, _ := strconv.Atoi(one[1])
		l2, _ := strconv.Atoi(two[0])
		h2, _ := strconv.Atoi(two[1])

		if debug {
			println(l)
		}
		if contains(l1, h1, l2, h2) {
			if debug {
				println("MATCH")
			}
			currentTotal++
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	println(currentTotal)

}
