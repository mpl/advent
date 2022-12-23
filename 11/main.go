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
	items       []int
	operation   func(int) int
	test        func(int) bool
	monkeyTrue  int
	monkeyFalse int
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
	f, err := os.Open("./input.txt.demo")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	var monkeys []*monkey

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
		kong.test = func(i int) bool {
			return i%divisor == 0
		}

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

	for {
		for idx, mk := range monkeys {
			if *debug {
				println("MONKEY ", idx)
			}
			for k, item := range mk.items {
				if *debug {
					println(k, "OLD WORRY: ", item)
				}
				worry := mk.operation(item)
				worry = worry / 3
				mk.items[k] = worry
				if *debug {
					println(k, "NEW WORRY: ", worry)
				}

				if mk.test(worry) {
					if *debug {
						println(k, worry, "THROWING TO MK ", mk.monkeyTrue)
					}
				} else {
					if *debug {
						println(k, worry, "THROWING TO MK ", mk.monkeyFalse)
					}
				}
			}
		}
		break
	}

}

/*
Monkey 0:
  Monkey inspects an item with a worry level of 79.
    Worry level is multiplied by 19 to 1501.
    Monkey gets bored with item. Worry level is divided by 3 to 500.
    Current worry level is not divisible by 23.
    Item with worry level 500 is thrown to monkey 3.
  Monkey inspects an item with a worry level of 98.
    Worry level is multiplied by 19 to 1862.
    Monkey gets bored with item. Worry level is divided by 3 to 620.
    Current worry level is not divisible by 23.
    Item with worry level 620 is thrown to monkey 3.
Monkey 1:
  Monkey inspects an item with a worry level of 54.
    Worry level increases by 6 to 60.
    Monkey gets bored with item. Worry level is divided by 3 to 20.
    Current worry level is not divisible by 19.
    Item with worry level 20 is thrown to monkey 0.
  Monkey inspects an item with a worry level of 65.
    Worry level increases by 6 to 71.
    Monkey gets bored with item. Worry level is divided by 3 to 23.
    Current worry level is not divisible by 19.
    Item with worry level 23 is thrown to monkey 0.
  Monkey inspects an item with a worry level of 75.
    Worry level increases by 6 to 81.
    Monkey gets bored with item. Worry level is divided by 3 to 27.
    Current worry level is not divisible by 19.
    Item with worry level 27 is thrown to monkey 0.
  Monkey inspects an item with a worry level of 74.
    Worry level increases by 6 to 80.
    Monkey gets bored with item. Worry level is divided by 3 to 26.
    Current worry level is not divisible by 19.
    Item with worry level 26 is thrown to monkey 0.
Monkey 2:
  Monkey inspects an item with a worry level of 79.
    Worry level is multiplied by itself to 6241.
    Monkey gets bored with item. Worry level is divided by 3 to 2080.
    Current worry level is divisible by 13.
    Item with worry level 2080 is thrown to monkey 1.
  Monkey inspects an item with a worry level of 60.
    Worry level is multiplied by itself to 3600.
    Monkey gets bored with item. Worry level is divided by 3 to 1200.
    Current worry level is not divisible by 13.
    Item with worry level 1200 is thrown to monkey 3.
  Monkey inspects an item with a worry level of 97.
    Worry level is multiplied by itself to 9409.
    Monkey gets bored with item. Worry level is divided by 3 to 3136.
    Current worry level is not divisible by 13.
    Item with worry level 3136 is thrown to monkey 3.
Monkey 3:
  Monkey inspects an item with a worry level of 74.
    Worry level increases by 3 to 77.
    Monkey gets bored with item. Worry level is divided by 3 to 25.
    Current worry level is not divisible by 17.
    Item with worry level 25 is thrown to monkey 1.
  Monkey inspects an item with a worry level of 500.
    Worry level increases by 3 to 503.
    Monkey gets bored with item. Worry level is divided by 3 to 167.
    Current worry level is not divisible by 17.
    Item with worry level 167 is thrown to monkey 1.
  Monkey inspects an item with a worry level of 620.
    Worry level increases by 3 to 623.
    Monkey gets bored with item. Worry level is divided by 3 to 207.
    Current worry level is not divisible by 17.
    Item with worry level 207 is thrown to monkey 1.
  Monkey inspects an item with a worry level of 1200.
    Worry level increases by 3 to 1203.
    Monkey gets bored with item. Worry level is divided by 3 to 401.
    Current worry level is not divisible by 17.
    Item with worry level 401 is thrown to monkey 1.
  Monkey inspects an item with a worry level of 3136.
    Worry level increases by 3 to 3139.
    Monkey gets bored with item. Worry level is divided by 3 to 1046.
    Current worry level is not divisible by 17.
    Item with worry level 1046 is thrown to monkey 1.
*/
