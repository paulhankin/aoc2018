package main

import "fmt"

func power(sn, sz int) (string, int) {
	pow := func(x, y int) int {
		p := (x + 10) * y + sn
		p *= (x + 10)
		p /= 100
		p %= 10
		return p - 5
	}
	best := 0
	bestij := ""
	var acc [301][301]int
	for i := 0; i < 301; i++ {
		acc[i][0] = 0
		acc[0][i] = 0
	}

	for i := 1; i < 301; i++ {
		for j := 1; j < 301; j++ {
			acc[i][j] = acc[i-1][j] + acc[i][j-1] + pow(i, j) - acc[i-1][j-1]
		}
	}
	for i := 1; i < 300 - sz; i++ {
		for j := 1; j < 300 - sz; j++ {
			sum := acc[i+sz-1][j+sz-1] - acc[i-1][j+sz-1] - acc[i+sz-1][j-1] + acc[i-1][j-1]
			if sum > best {
				best = sum
				bestij = fmt.Sprintf("%d,%d", i, j)
			}
		}
	}
	return bestij, best
}

func main() {
	ij, _ := power(5034, 3)
	fmt.Println(ij)
	best := 0
	bestr := ""
	for sz := 3; sz <= 300; sz++ {
		if x, v := power(5034, sz); v > best {
			bestr = fmt.Sprintf("%s,%d", x, sz)
			best = v
		}
	}
	fmt.Println(bestr)
}
