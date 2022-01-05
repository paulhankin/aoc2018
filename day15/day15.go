package main

import (
	"fmt"
	"sort"
	"strings"
)

var input = `################################
#########################.G.####
#########################....###
##################.G.........###
##################.##.......####
#################...#.........##
################..............##
######..########...G...#.#....##
#####....######.G.GG..G..##.####
#######.#####G............#.####
#####.........G..G......#...####
#####..G......G..........G....##
######GG......#####........E.###
#######......#######..........##
######...G.G#########........###
######......#########.....E..###
#####.......#########........###
#####....G..#########........###
######.##.#.#########......#####
#######......#######.......#####
#######.......#####....E...#####
##.G..#.##............##.....###
#.....#........###..#.#.....####
#.........E.E...#####.#.#....###
######......#.....###...#.#.E###
#####........##...###..####..###
####...G#.##....E####E.####...##
####.#########....###E.####....#
###...#######.....###E.####....#
####..#######.##.##########...##
####..######################.###
################################`

type creature struct {
	t    byte
	x, y int
	hp   int
}

func mklocs(cs []creature) map[[2]int]byte {
	r := map[[2]int]byte{}
	for _, c := range cs {
		r[[2]int{c.x, c.y}] = c.t
	}
	return r
}

func gridEmpty(g [][]byte, locs map[[2]int]byte, x, y int) bool {
	if x < 0 || x >= len(g[0]) || y < 0 || y >= len(g) {
		return false
	}
	if g[y][x] != '.' {
		return false
	}
	if locs[[2]int{x, y}] != 0 {
		return false
	}
	return true
}

// allPathreturns a map from location to the length of the shortest path
// to that location from x, y, not passing through any creatures.
func allPath(grid [][]byte, x, y int, locs map[[2]int]byte) map[[2]int]int {
	qd := map[[2]int]bool{{x, y}: true}
	todo := [][3]int{{x, y, 0}}
	dists := map[[2]int]int{}
	for len(todo) > 0 {
		tx, ty, td := todo[0][0], todo[0][1], todo[0][2]
		dists[[2]int{tx, ty}] = td
		todo = todo[1:]
		ds := [][2]int{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
		for _, dd := range ds {
			nx := tx + dd[0]
			ny := ty + dd[1]
			k := [2]int{nx, ny}
			if !qd[k] && gridEmpty(grid, locs, nx, ny) {
				todo = append(todo, [3]int{nx, ny, td + 1})
				qd[k] = true
			}
		}
	}
	return dists
}

func step(grid [][]byte, cs []creature) []creature {
	sort.Slice(cs, func(i, j int) bool {
		if cs[i].y != cs[j].y {
			return cs[i].y < cs[j].y
		}
		return cs[i].x < cs[j].x
	})
	locs := mklocs(cs)
	for ci, c := range cs {
		if c.hp <= 0 {
			continue
		}
		dists := allPath(grid, c.x, c.y, locs)
		bestx, besty := 0, 0
		best := 9999
		for _, d := range cs {
			if d.hp <= 0 || d.t == c.t {
				continue
			}
			ds := [][2]int{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
			for _, dd := range ds {
				nx := d.x + dd[0]
				ny := d.y + dd[1]
				if d, ok := dists[[2]int{nx, ny}]; ok && d <= best {
					if d == best {
						if ny > besty || ny == besty && nx >= bestx {
							continue
						}
					}
					best = d
					bestx = nx
					besty = ny
				}
			}
		}
		if bestx != c.x || besty != c.y {
			fmt.Printf("%d: moving to %d %d\n", ci, bestx, besty)
		}
	}
	return nil
}

func showGrid(grid [][]byte, cs []creature, dists map[[2]int]int) {
	locs := mklocs(cs)
	for i, g := range grid {
		for j := range g {
			if b := locs[[2]int{j, i}]; b != 0 {
				fmt.Printf(" %c ", b)
			} else if d := dists[[2]int{j, i}]; d > 0 && d < 100 {
				fmt.Printf("%2d ", d)
			} else {
				fmt.Printf("%c%c%c", g[j], g[j], g[j])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	var grid [][]byte
	var cs []creature
	for i, line := range strings.Split(input, "\n") {
		g := []byte(line)
		for j := range g {
			b := g[j]
			if b == 'G' {
				g[j] = '.'
				cs = append(cs, creature{t: 'G', x: j, y: i, hp: 200})
			} else if b == 'E' {
				g[j] = '.'
				cs = append(cs, creature{t: 'E', x: j, y: i, hp: 200})
			} else if b != '.' && b != '#' {
				panic(b)
			}
		}
		grid = append(grid, g)
	}
	d := allPath(grid, cs[0].x, cs[0].y, mklocs(cs))
	showGrid(grid, cs, d)
	step(grid, cs)
}
