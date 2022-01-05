package main

import (
	"fmt"
	"strconv"
)

func matches(target string, x []int) bool {
	if len(x) != len(target) {
		panic("matches")
	}
	for i, d := range x {
		if int(target[i]-'0') != d {
			return false
		}
	}
	return true
}

func dorecipes(x int, target string) (string, int) {
	r := []int{3, 7}
	e0 := 0
	e1 := 1
	gottarget := 0
	for steps := 0; len(r) <= x+11 || gottarget == 0; steps++ {
		n := r[e0] + r[e1]
		if n > 9 {
			r = append(r, n/10, n%10)
		} else {
			r = append(r, n)
		}
		if gottarget == 0 && len(r) > len(target) {
			if matches(target, r[len(r)-len(target):]) {
				gottarget = len(r) - len(target)
			}
		}
		if gottarget == 0 && len(r) > len(target)+1 {
			if matches(target, r[len(r)-len(target)-1:len(r)-1]) {
				gottarget = len(r) - len(target) - 1
			}
		}
		e0 = (e0 + r[e0] + 1) % len(r)
		e1 = (e1 + r[e1] + 1) % len(r)
	}
	R := ""
	for i := 0; i < 10; i++ {
		R += strconv.Itoa(r[x+i])
	}
	return R, gottarget
}

func main() {
	cases := map[int]string{
		// 9:      "51589",
		// 5:      "01245",
		// 18:     "92510",
		// 2018:   "59414",
		409551: "409551",
	}
	for i, j := range cases {
		x, y := dorecipes(i, j)
		fmt.Println(x)
		fmt.Println(y)
	}
}
