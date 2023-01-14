package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	debug      = flag.Bool("debug", false, `debug mode`)
	verbose    = flag.Bool("verbose", false, `verbose mode`)
	demo       = flag.Bool("demo", false, `use demo input`)
	debugStart = flag.Int("dbgstart", -1, `debug start index`)
	debugEnd   = flag.Int("dbgend", -1, `debug end index`)
)

type signal string

type signals []signal

var (
	sigs signals = signals{
		signal("[2]"),
		signal("[6]"),
	}

	debugCount = 2
)

type pair struct {
	left  string
	right string
}

func initSignals() {
	var f *os.File
	var err error
	if *demo {
		f, err = os.Open("../input.txt.demo")
	} else {
		f, err = os.Open("../input.txt")
	}
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := sc.Text()
		sigs = append(sigs, signal(l[1:len(l)-1]))
		sc.Scan()
		l = sc.Text()
		sigs = append(sigs, signal(l[1:len(l)-1]))
		sc.Scan()
	}
}

func (s signals) Len() int      { return len(s) }
func (s signals) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s signals) Less(i, j int) bool {
	return compareLists(string(s[i]), string(s[j])) < 0
}

func main() {
	flag.Parse()

	initSignals()

	sort.Sort(sigs)

	divider1, divider2 := 0, 0
	for k, v := range sigs {
		if *debug {
			println(v)
		}
		if string(v) == "[2]" {
			divider1 = k + 1
			continue
		}
		if string(v) == "[6]" {
			divider2 = k + 1
			continue
		}
	}
	println("DECODER: ", divider1*divider2)
}

func compareLists(leftInput, rightInput string) int {
	lParts := split(leftInput)
	rParts := split(rightInput)

	i := 0
	for {
		if len(lParts) < len(rParts) {
			if i+1 > len(lParts) {
				return -1
			}
		}
		if len(rParts) < len(lParts) {
			if i+1 > len(rParts) {
				return 1
			}
		}
		if len(rParts) == len(lParts) {
			if i+1 > len(lParts) {
				break
			}
		}
		left := lParts[i]
		right := rParts[i]
		i++
		if *debug {
			println(left, "VS", right)
		}

		asLists := false
		leftIsList := strings.Contains(left, ",") || strings.HasPrefix(left, "[")
		rightIsList := strings.Contains(right, ",") || strings.HasPrefix(right, "[")

		if rightIsList || leftIsList {
			asLists = true
		}

		if asLists {
			cmpL := compareLists(left, right)
			if cmpL == 0 {
				continue
			}
			return cmpL
		}

		cmp := compareInts(left, right)
		if cmp == 0 {
			continue
		}
		return cmp
	}
	return 0
}

func split(input string) []string {
	if input == "" {
		return []string{}
	}
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

func compareInts(left, right string) int {
	if left == "" && right == "" {
		return 0
	}
	if left == "" {
		return -1
	}
	if right == "" {
		return 1
	}
	nbLeft, err := strconv.Atoi(left)
	if err != nil {
		panic(err.Error())
	}
	nbRight, err := strconv.Atoi(right)
	if err != nil {
		panic(err.Error())
	}
	if nbLeft < nbRight {
		return -1
	}
	if nbLeft > nbRight {
		return 1
	}
	return 0
}
