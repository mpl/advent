package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

//[T]             [P]     [J]
//[F]     [S]     [T]     [R]     [B]
//[V]     [M] [H] [S]     [F]     [R]
//[Z]     [P] [Q] [B]     [S] [W] [P]
//[C]     [Q] [R] [D] [Z] [N] [H] [Q]
//[W] [B] [T] [F] [L] [T] [M] [F] [T]
//[S] [R] [Z] [V] [G] [R] [Q] [N] [Z]
//[Q] [Q] [B] [D] [J] [W] [H] [R] [J]
// 1   2   3   4   5   6   7   8   9

var debug = false

var stacks [][]string

var indexMax = 7

func initStacks() {
	sizes := []int{8, 3, 7, 6, 8, 4, 8, 5, 7}
	for i := 0; i < 9; i++ {
		stacks = append(stacks, make([]string, sizes[i]))
	}
}

func main() {
	initStacks()
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	initCount := 0
	for sc.Scan() {
		if initCount == indexMax+1 {
			break
		}

		l := sc.Text()
		// println(initCount)
		for k, v := range l {
			colIdx := k / 4
			s := string(v)
			switch s {
			case "[", "]", " ":
				continue
			}
			stack := stacks[colIdx]
			// println(s, len(stack))
			stack[indexMax-initCount] = s
			stacks[colIdx] = stack
		}

		initCount++

	}

	sc.Scan()

	// move 3 from 8 to 2
	// move 3 from 1 to 5

	debugCount := 0
	for sc.Scan() {
		if debug {
			if debugCount > 0 {
				break
			}
			debugCount++
		}

		l := sc.Text()
		parts := strings.Split(l, " ")
		N, _ := strconv.Atoi(parts[1])
		from, _ := strconv.Atoi(parts[3])
		from--
		to, _ := strconv.Atoi(parts[5])
		to--

		move(from, to, N)

		if debug {
			spew.Dump(stacks)
		}

	}

	var topStacks string
	for _,v := range stacks {
		topStacks += v[len(v)-1]
	}
	println(topStacks)
}

func move(from, to, N int) {
	stackFrom := stacks[from]
	stackTo := stacks[to]

	popped := stackFrom[len(stackFrom)-N:]
	truncated := stackFrom[:len(stackFrom)-N]

	for i := len(popped) - 1; i >= 0; i-- {
		stackTo = append(stackTo, popped[i])
	}
	stacks[from] = truncated
	stacks[to] = stackTo
}
