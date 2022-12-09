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
			if debugCount > 1 {
				break
			}
			debugCount++
		}

		l := sc.Text()
		if debug {
			println(l)
		}

		prevElf := make(map[rune]int)
		for _, v := range l {
			prevElf[v] = 1
		}

		elf := make(map[rune]int)
		sc.Scan()
		l = sc.Text()
		if debug {
			println(l)
		}

		for _, v := range l {
			if _, ok := prevElf[v]; ok {
				elf[v] = 1
			}
		}

		prevElf = elf
		elf = make(map[rune]int)
		sc.Scan()
		l = sc.Text()
		if debug {
			println(l)
		}
		for k, v := range l {
			if _, ok := prevElf[v]; !ok {
				continue
			}
			currentTotal += priority(v)
			if debug {
				println(k, v, string(v), priority(v))
			}
			break
		}

	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	println(currentTotal)

}
