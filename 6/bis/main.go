package main

import (
	"io/ioutil"
	"log"

//	"github.com/davecgh/go-spew/spew"
)


// bhdhvvsqvqhvvfrfddbpwgvztw

var debug = false

func isMarker(seq []byte) (bool, int) {
	asMap := make(map[byte]int)
	idx := 0
	for k,v := range seq {
		prev, ok := asMap[v]
		if ok {
			idx = prev+1
		}
		asMap[v] = k
	}
	return len(asMap) == 14, idx
}

func main() {
	data, err := ioutil.ReadFile("../input.txt")
	if err != nil {
		log.Fatal(err)
	}

	pos := 0
	seq := data[:14]
	for {
		isM, idx := isMarker(seq)
		if isM {		
			println("FOUND", string(seq))
			pos += 14
			break
		}
		pos += idx
		if pos >= len(data) {
			break
		}
		seq = data[pos:pos+14]
	}

	println(string(data[pos-14:pos]))
	println(pos)
}
