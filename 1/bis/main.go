package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	f, err := os.Open("../input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	top3 := make([]int, 3)
	filled := false

	debugCount := 0

	currentTotal, max := 0, 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := sc.Text()
		if l == "" {
			if debugCount > 5 {
				// return
			}
			debugCount++

			if currentTotal <= max {
				currentTotal = 0
				continue
			}

			println("NEW MAX: ", currentTotal)

			max = currentTotal
			currentTotal = 0
			if filled {
				top3[2] = max
				sort.Sort(sort.Reverse(sort.IntSlice(top3)))
				max = top3[2]
			} else {
				var maxIndex int
				for i, v := range top3 {
					if v != 0 {
						continue
					}
					top3[i] = max
					if i == 2 {
						filled = true
					}
					maxIndex = i
					break
				}
				sort.Sort(sort.Reverse(sort.IntSlice(top3)))
				max = top3[maxIndex]
			}

			for _, v := range top3 {
				fmt.Printf("%d, ", v)
			}
			println()
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

	grandTotal := 0
	for _, v := range top3 {
		grandTotal += v
	}

	println(grandTotal)

}
