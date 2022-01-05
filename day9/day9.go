package main

import "fmt"

type marble struct {
	value      int
	prev, next *marble
}

func madness(np, lm int) int {
	scores := make([]int, np)
	cm := &marble{value: 0}
	cm.next = cm
	cm.prev = cm
	for m := 1; m <= lm; m++ {
		if false {
			m0 := cm
			for {
				fmt.Printf("%d ", m0.value)
				m0 = m0.next
				if m0 == cm {
					fmt.Println()
					break
				}
			}
		}
		p := (m - 1) % np
		if m%23 == 0 {
			scores[p] += m
			for i := 0; i < 7; i++ {
				cm = cm.prev
			}
			scores[p] += cm.value
			cm.next.prev = cm.prev
			cm.prev.next = cm.next
			cm = cm.next
		} else {
			nm := &marble{prev: cm.next, next: cm.next.next, value: m}
			cm.next.next.prev = nm
			cm.next.next = nm
			cm = nm
		}
	}
	best := 0
	for _, s := range scores {
		if s > best {
			best = s
		}
	}
	return best
}

func main() {
	fmt.Println(madness(425, 70848))
	fmt.Println(madness(425, 7084800))
}
