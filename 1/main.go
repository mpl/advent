package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	currentTotal, max := 0, 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := sc.Text()
		if l == "" {
			if currentTotal > max {
				println("NEW MAX: ", currentTotal)
				max = currentTotal
			}
			currentTotal = 0
			continue
		}
		currentCalorie, err := strconv.Atoi(l)
		if err != nil {
			log.Fatal(err)
		}
		currentTotal += currentCalorie
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

}
