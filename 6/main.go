package main

import (
	"io/ioutil"
	"log"

//	"github.com/davecgh/go-spew/spew"
)

var debug = false

func isMarker(seq []byte) bool {
	var prev byte
	asMap := make(map[byte]bool)
	for _,v := range seq {
		// fastPath in case of two consecutive identical
		if v == prev {
			return false
		}
		prev = v
		asMap[v] = true
	}
	return len(asMap) == 4
}

func main() {
	data, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	seq := data[:4]

	pos := 4
	for k,v := range data[4:] {
		if isMarker(seq) {
			println("FOUND AT ", k, string(seq))
			break
		}
		seq = append(seq[1:], v)
		pos++
	}

	println(string(data[pos-4:pos]))
	println(pos)

}
