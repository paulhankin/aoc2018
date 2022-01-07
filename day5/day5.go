package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

type chain struct {
	b          byte
	next, prev *chain
}

func toggleCase(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 'a' + 'A'
	} else {
		return b - 'A' + 'a'
	}
}

func react(input []byte) int {
	head := &chain{}
	head.next = head
	head.prev = head
	c := head
	for i := range input {
		b := input[i]
		if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' {
			c2 := chain{b: b, next: c.next, prev: c}
			c.next = &c2
			c = &c2
		}
	}

	x := c.next
	for {
		if x.next == head {
			break
		}
		if toggleCase(x.b) == x.next.b {
			x.prev.next = x.next.next
			x.next.next.prev = x.prev
			if x.prev == head {
				x = x.next.next
			} else {
				x = x.prev
			}
		} else {
			x = x.next
		}
	}

	x = head
	sum := 0
	for {
		x = x.next
		if x == head {
			break
		}
		sum++
	}
	return sum
}

func main() {
	input, err := ioutil.ReadFile("./day5/day5.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(react(input))
	best := len(input) + 1
	for i := 0; i < 26; i++ {
		var ni bytes.Buffer
		for _, b := range input {
			if b-'a' == byte(i) || b-'A' == byte(i) {
				continue
			}
			ni.WriteByte(b)
		}
		if x := react(ni.Bytes()); x < best {
			best = x
		}
	}
	fmt.Println(best)

}
