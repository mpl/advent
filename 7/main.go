package main

import (
	"flag"
	"log"
	"os"
	"bufio"
	"strings"
	"strconv"
)

var debug = flag.Bool("debug", false, `debug mode`)

type node struct {
	name string
	isDir bool
	size int
	children map[string]*node
	parent *node
}

func (n node) print(level int) {
	ident := ""
	for i:=0; i<level; i++ {
		ident += "	"
	}
	if !n.isDir {
		println(ident, n.name, n.size)
		return
	}
	println(ident, n.name, n.size, "->")
	for _,v := range n.children {
		v.print(level+1)
	}
}

var cwd *node

var root *node = &node {
	name: "/",
	children: make(map[string]*node),
	isDir: true,
}

func moveTo(dest string) {
	if dest == "/" {
		cwd = root
		return
	}
	if dest == ".." {
		cwd = cwd.parent
		return
	}

	from := cwd
	if from.children == nil {
		from.children = make(map[string]*node)
	}
	child, ok := from.children[dest]
	if ok {
		cwd = child
		return
	}
	nd := &node {
		name: dest,
		parent: from,
		isDir: true,
		children: make(map[string]*node),
	}
	from.children[dest] = nd
	cwd = nd
}

func populateSize(n *node) int {
	if !n.isDir {
		return n.size
	}
	size := 0
	for _,v := range n.children {
		if !v.isDir {
			size += v.size
			continue
		}
		size += populateSize(v)
	}
	n.size = size
	return n.size
}

func walk(root *node, fn func(n *node) int) int {
	res := fn(root)
	for _, v := range root.children {
		res += walk(v, fn)
	}
	return res
}

func main() {
	flag.Parse()
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	debugCount := 0
	for sc.Scan() {
	if *debug {
			if debugCount > 100 {
				// break
			}
			debugCount++
		}
		l := sc.Text()
		// println(l)
		if strings.HasPrefix(l, "$") {
			command := strings.TrimPrefix(l, "$ ")
			if strings.HasPrefix(command, "cd") {
				moveTo(strings.TrimPrefix(command, "cd "))
				continue
			}
			if !strings.HasPrefix(command, "ls") {
				panic("OIOIOI")
			}
			continue
		}

		// output of an ls
		parts := strings.Split(l, " ")
		name := parts[1]
		if cwd.children == nil {
			cwd.children = make(map[string]*node)
		}
		child, ok := cwd.children[name]
		if !ok {
			child = &node {
				name: name,
				parent: cwd,
			}
		}
		if parts[0] == "dir" {
			child.isDir = true
			if child.children == nil {
				child.children = make(map[string]*node)
			}
		} else {
			size, _ := strconv.Atoi(parts[0])
			child.size = size
		}
		cwd.children[name] = child
		continue
	}

	println(populateSize(root))

	if *debug {
		root.print(0)
	}

	totalSize := walk(root, func(n *node) int {
		if !n.isDir {
			return 0
		}
		if n.size > 100000 {
			return 0
		}
		return n.size
	})

	println(totalSize)


}
