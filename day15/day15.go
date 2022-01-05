package main

import (
	"fmt"
	"sort"
	"strings"
)

var verbose = false

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

var ex = `#########
#G..G..G#
#.......#
#.......#
#G..E..G#
#.......#
#.......#
#G..G..G#
#########`

var ex2 = `#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######`

var ex3 = `#######
#G..#E#
#E#E.E#
#G.##.#
#...#E#
#...E.#
#######`

var ex4 = `#######
#.E...#
#.#..G#
#.###.#
#E#G#G#
#...#G#
#######`

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

var dirs = [][2]int{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

func findCreatureAt(cs []creature, x, y int) *creature {
	for i, c := range cs {
		if c.hp > 0 && c.x == x && c.y == y {
			return &cs[i]
		}
	}
	return nil
}

func step(grid [][]byte, cs []creature, elfpow int) ([]creature, bool) {
	sort.Slice(cs, func(i, j int) bool {
		if cs[i].y != cs[j].y {
			return cs[i].y < cs[j].y
		}
		return cs[i].x < cs[j].x
	})
	locs := mklocs(cs)
	end := false
	for ci := range cs {
		c := &cs[ci]
		if c.hp <= 0 {
			continue
		}
		dists := allPath(grid, c.x, c.y, locs)
		bestx, besty := -1, -1
		best := 9999
		notargs := true
		for _, d := range cs {
			if d.hp <= 0 || d.t == c.t {
				continue
			}
			notargs = false
			for _, dd := range dirs {
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
		if notargs {
			end = true
		}
		if bestx != c.x || besty != c.y {
			if bestx == -1 || besty == -1 {
				// We'd like to move, but can't
				continue
			}
			rd := allPath(grid, bestx, besty, locs)
			nx, ny := -1, -1
			bestd := 99999
			for _, dd := range dirs {
				nx0 := c.x + dd[0]
				ny0 := c.y + dd[1]
				d, ok := rd[[2]int{nx0, ny0}]
				if !ok {
					continue
				}
				if d < bestd || d == bestd && (ny0 < ny || ny0 == ny && nx0 < nx) {
					nx, ny = nx0, ny0
					bestd = d
				}
			}
			if nx == -1 || ny == -1 {
				showGrid(grid, cs, rd)
				panic("failed to find step!")
			}
			delete(locs, [2]int{c.x, c.y})
			cs[ci].x = nx
			cs[ci].y = ny
			locs[[2]int{nx, ny}] = c.t
		}
		// We either didn't move because we're in range, or we moved to our target square.
		// HIT!
		if bestx == c.x && besty == c.y {
			var bt *creature
			for _, dd := range dirs {
				nx0 := c.x + dd[0]
				ny0 := c.y + dd[1]
				if v := locs[[2]int{nx0, ny0}]; v == 0 || v == c.t {
					continue
				}
				if targ := findCreatureAt(cs, nx0, ny0); targ != nil {
					if targ.hp > 0 && bt == nil || (targ.hp < bt.hp || targ.hp == bt.hp && (targ.y < bt.y || targ.y == bt.y && targ.x < bt.x)) {
						bt = targ
					}
				}
			}
			if bt == nil {
				panic("wanted to hit, but didn't find a target?")
			}
			if c.t == 'E' {
				bt.hp -= elfpow
			} else {
				bt.hp -= 3
			}
			if bt.hp <= 0 {
				delete(locs, [2]int{bt.x, bt.y})
			}
		}
	}
	var ncs []creature
	for _, c := range cs {
		if c.hp > 0 {
			ncs = append(ncs, c)
		}
	}
	return ncs, end
}

func showGrid(grid [][]byte, cs []creature, dists map[[2]int]int) {
	locs := mklocs(cs)
	for i, g := range grid {
		for j := range g {
			if dists == nil {
				if b := locs[[2]int{j, i}]; b != 0 {
					fmt.Printf("%c", b)
				} else {
					fmt.Printf("%c", g[j])
				}
			} else {
				if b := locs[[2]int{j, i}]; b != 0 {
					fmt.Printf(".%c.", b)
				} else if d := dists[[2]int{j, i}]; d > 0 && d < 100 {
					fmt.Printf("%2d ", d)
				} else {
					fmt.Printf("%c%c%c", g[j], g[j], g[j])
				}
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

	dosim := func(elfpow int, grid [][]byte, cs []creature) ([]creature, int) {
		steps := 0
		for ; ; steps++ {
			if verbose {
				plural := "s"
				if steps == 1 {
					plural = ""
				}
				fmt.Printf("After %d round%s:\n", steps, plural)
				for _, c := range cs {
					fmt.Printf("%c(%d), ", c.t, c.hp)
				}
				fmt.Println()
				showGrid(grid, cs, nil)
			}
			end := false
			cs, end = step(grid, cs, elfpow)
			if end {
				break
			}
		}
		sum := 0
		for _, c := range cs {
			if c.hp > 0 {
				sum += c.hp
			}
		}
		return cs, sum * steps
	}
	_, val := dosim(3, grid, append([]creature{}, cs...))
	fmt.Println(val)
	countelf := func(cs []creature) int {
		sum := 0
		for _, c := range cs {
			if c.t == 'E' {
				sum++
			}
		}
		return sum
	}
	nelf0 := countelf(cs)
	for elfpow := 3; elfpow < 200; elfpow++ {
		ncs, val := dosim(elfpow, grid, append([]creature{}, cs...))
		nelf := countelf(ncs)
		if nelf == nelf0 {
			fmt.Println(val)
			break
		}
	}
}
