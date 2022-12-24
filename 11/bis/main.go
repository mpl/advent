package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var debug = flag.Bool("debug", false, `debug mode`)

type monkey struct {
	items        []int
	operation    func(int) int
	test         func(int) bool
	monkeyTrue   int
	monkeyFalse  int
	inspectCount int
}

func parseOperation(operation string) func(int) int {
	parts := strings.Split(operation, " ")
	self := false
	if parts[1] == "old" {
		self = true
	}
	operand, _ := strconv.Atoi(parts[1])
	sign := string(operation[0])
	if sign == "+" {
		return func(i int) int {
			if self {
				return i + i
			}
			return i + operand
		}
	}
	if sign == "*" {
		return func(i int) int {
			if self {
				return i * i
			}
			return i * operand
		}
	}
	panic(fmt.Sprintf("unsupported operation: %s", operation))
}

func main() {
	flag.Parse()
	// f, err := os.Open("../input.txt.demo")
	f, err := os.Open("../input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	var monkeys []*monkey

	commonDivisor := 1
	debugCount := 0
	for sc.Scan() {
		if *debug {
			if debugCount > 34 {
				break
			}
		}
		debugCount++

		l := sc.Text()
		if !strings.HasPrefix(l, "Monkey") {
			panic("BIM")
		}
		kong := monkey{}

		sc.Scan()
		l = sc.Text()
		items := strings.Split(strings.TrimPrefix(l, "  Starting items: "), ", ")
		for _, v := range items {
			item, _ := strconv.Atoi(v)
			kong.items = append(kong.items, item)
		}

		sc.Scan()
		l = sc.Text()
		operation := strings.TrimPrefix(l, "  Operation: new = old ")
		kong.operation = parseOperation(operation)

		sc.Scan()
		l = sc.Text()
		if !strings.HasPrefix(l, "  Test: divisible by ") {
			panic(fmt.Sprintf("unsupported test: %s", l))
		}
		divisor, _ := strconv.Atoi(strings.TrimPrefix(l, "  Test: divisible by "))
		if *debug {
			println(len(monkeys), "DIVISOR: ", divisor)
		}
		// TODO: we should defer, and make it directly return the monkey destination
		kong.test = func(i int) bool {
			return i%divisor == 0
		}
		commonDivisor *= divisor

		sc.Scan()
		l = sc.Text()
		kong.monkeyTrue, _ = strconv.Atoi(strings.TrimPrefix(l, "    If true: throw to monkey "))
		if *debug {
			println(len(monkeys), "MONKEYTRUE: ", kong.monkeyTrue)
		}

		sc.Scan()
		l = sc.Text()
		kong.monkeyFalse, _ = strconv.Atoi(strings.TrimPrefix(l, "    If false: throw to monkey "))
		if *debug {
			println(len(monkeys), "MONKEYFALSE: ", kong.monkeyFalse)
		}

		monkeys = append(monkeys, &kong)

		sc.Scan()
	}

	for i := 0; i < 10000; i++ {
		for idx, mk := range monkeys {
			if *debug {
				println("MONKEY ", idx)
			}
			for k, item := range mk.items {
				if *debug {
					println(k, "OLD WORRY: ", item)
				}
				worry := mk.operation(item)
				worry %= commonDivisor
				// mk.items[k] = worry
				if *debug {
					println(k, "NEW WORRY: ", worry)
				}

				monkeyDest := 0
				if mk.test(worry) {
					monkeyDest = mk.monkeyTrue
				} else {
					monkeyDest = mk.monkeyFalse
				}
				if *debug {
					println(k, worry, "THROWING TO MK ", monkeyDest)
				}
				monkeys[monkeyDest].items = append(monkeys[monkeyDest].items, worry)
				mk.inspectCount++
			}
			mk.items = []int{}
		}
	}

	for i, mk := range monkeys {
		println("Monkey", i, "inspected items", mk.inspectCount, "times.")
	}

	/*
	   Monkey 0 inspected items 60002 times.
	   Monkey 1 inspected items 56692 times.
	   Monkey 2 inspected items 63353 times.
	   Monkey 3 inspected items 116668 times.
	   Monkey 4 inspected items 60034 times.
	   Monkey 5 inspected items 116628 times.
	   Monkey 6 inspected items 60010 times.
	   Monkey 7 inspected items 63343 times.
	*/

	println(116668 * 116628)
	println(uint64(116668) * uint64(116628))

	// Answer: 13606755504
}
