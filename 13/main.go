package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	debug   = flag.Bool("debug", false, `debug mode`)
	verbose = flag.Bool("verbose", false, `verbose mode`)
	demo    = flag.Bool("demo", false, `use demo input`)
)

var (
	pairs []pair

	debugCount = 2
)

type pair struct {
	left  string
	right string
}

func initPairs() {
	var f *os.File
	var err error
	if *demo {
		f, err = os.Open("./input.txt.demo")
	} else {
		f, err = os.Open("./input.txt")
	}
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		p := pair{}
		l := sc.Text()
		p.left = l
		sc.Scan()
		l = sc.Text()
		p.right = l
		pairs = append(pairs, p)
		sc.Scan()
	}
}

func main() {
	flag.Parse()

	initPairs()

	if *debug {
	}

	count := 0
	for _, v := range pairs {
		if *debug {
			if count == debugCount {
				break
			}
		}
		count++
		lParts := split(v.left)
		rParts := split(v.right)

		for k, left := range lParts {
			right := rParts[k]
			if *debug {
				println(left, "VS", right)
			}

			asLists := false
			leftIsList := strings.Contains(left, ",")
			rightIsList := strings.Contains(right, ",")

			if leftIsList && !rightIsList {
				// TODO: fixright
				asLists = true
			}

			if rightIsList && !leftIsList {
				// TODO: fixleft
				asLists = true
			}

			if rightIsList || leftIsList {
				asLists = true
			}

			if asLists {
				// TODO: recurse
				continue
			}

			cmp := compareInts(left, right)
			if cmp == 0 {
				continue
			}
			if cmp < 0 {
				println("correct order")
				break
			}
			if cmp > 0 {
				println("wrong order")
				break
			}
		}
	}

}

func compareInts(left, right) int {
	nbLeft, err := strconv.Atoi(left)
	if err != nil {
		panic(err.Error())
	}
	nbRight, err := strconv.Atoi(right)
	if err != nil {
		panic(err.Error())
	}
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func split(input string) []string {
	input = input[1 : len(input)-1]
	if !strings.Contains(input, "[") {
		return strings.Split(input, ",")
	}
	var parts []string
	for {
		if *debug {
			msh, _ := json.Marshal(parts)
			println("PARTS: ", string(msh))
		}
		if len(input) == 0 {
			break
		}
		if input[0] != '[' {
			// that part is not a list, so it should be a number
			idx := strings.Index(input, ",")
			if idx == -1 {
				// that last part consumed all the input left
				parts = append(parts, input)
				break
			}
			// number consumed and added to the parts, keep going
			parts = append(parts, input[:idx])
			input = input[idx+1:]
			continue
		}
		// we should be dealing with a list here
		idx := indexClosing(input)
		if idx == -1 {
			panic("BOOM")
		}
		if *debug {
			println("IDX", idx)
		}
		parts = append(parts, input[1:idx])
		input = strings.TrimPrefix(input[idx+1:], ",")
	}
	return parts
}

func indexClosing(input string) int {
	if input[0] != '[' {
		// redundant with the caller, but oh well.
		panic("BIM")
	}
	n := 0
	for k, v := range input {
		if v == '[' {
			n++
			continue
		}
		if v == ']' {
			n--
			if n == 0 {
				return k
			}
		}
	}
	return -1
}
