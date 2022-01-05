package main

import (
	"fmt"
	"strings"
)

var input = `83, 153
201, 74
291, 245
269, 271
222, 337
291, 271
173, 346
189, 184
170, 240
127, 96
76, 46
92, 182
107, 160
311, 142
247, 321
303, 295
141, 310
147, 70
48, 41
40, 276
46, 313
175, 279
149, 177
181, 189
347, 163
215, 135
103, 159
222, 304
201, 184
272, 354
113, 74
59, 231
302, 251
127, 312
259, 259
41, 244
43, 238
193, 172
147, 353
332, 316
353, 218
100, 115
111, 58
210, 108
101, 175
185, 98
256, 311
142, 41
68, 228
327, 194`

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	return -min(-x, -y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func findNearest(cs [][2]int, x, y int) int {
	best := 0
	besti := -1
	ok := false
	for ci, c := range cs {
		d := abs(x-c[0]) + abs(y-c[1])
		if besti == -1 || d < best {
			besti = ci
			best = d
			ok = true
		} else if d == best {
			ok = false
		}
	}
	if !ok {
		return -1
	}
	return besti
}

func totalDistance(cs [][2]int, x, y int) int {
	var sum int
	for _, c := range cs {
		sum += abs(x-c[0]) + abs(y-c[1])
	}
	return sum
}

func main() {
	lines := strings.Split(input, "\n")
	var coords [][2]int
	var mx, my, Mx, My int
	for _, line := range lines {
		var x, y int
		if _, err := fmt.Sscanf(line, "%d, %d", &x, &y); err != nil {
			panic(err)
		}
		coords = append(coords, [2]int{x, y})
		mx = min(mx, x)
		Mx = max(Mx, x)
		my = min(my, y)
		My = max(My, y)
	}
	WIDTH := Mx - mx + 1
	HEIGHT := My - my + 1
	grid := make([]int, WIDTH*HEIGHT)
	get := func(grid []int, i, j int) int {
		return grid[(i-mx)+WIDTH*(j-my)]
	}
	set := func(grid []int, i, j, v int) {
		grid[(i-mx)+WIDTH*(j-my)] = v
	}
	for i := mx; i <= Mx; i++ {
		for j := my; j <= My; j++ {
			n := findNearest(coords, i, j)
			if n != -1 {
				set(grid, i, j, n)
			}
		}
	}
	edges := map[int]bool{}
	for i := mx; i <= Mx; i++ {
		edges[get(grid, i, my)] = true
		edges[get(grid, i, My)] = true
	}
	for i := my; i <= My; i++ {
		edges[get(grid, mx, i)] = true
		edges[get(grid, Mx, i)] = true
	}
	counts := map[int]int{}
	td := make([]int, WIDTH*HEIGHT)
	for i := mx; i <= Mx; i++ {
		for j := my; j <= My; j++ {
			counts[get(grid, i, j)]++
			set(td, i, j, totalDistance(coords, i, j))
		}
	}
	best := 0
	for k, v := range counts {
		if !edges[k] && v >= best {
			best = v
		}
	}
	fmt.Println(best)

	sum := 0
	for i := mx; i <= Mx; i++ {
		for j := my; j <= My; j++ {
			t := get(td, i, j)
			if t < 10000 {
				sum += 1
				if i == mx || i == Mx || j == my || j == My {
					// check that we don't have any points on the edge of our grid
					// that are within our region
					fmt.Println(i, j, t)
				}
			}
		}
	}
	fmt.Println(sum)
}
