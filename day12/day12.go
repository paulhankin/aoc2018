package main

import (
	"fmt"
	"strings"
)

var input = `#.#.#....##...##...##...#.##.#.###...#.##...#....#.#...#.##.........#.#...#..##.#.....#..#.###
####. => #
..#.. => .
#.#.. => .
.##.. => .
##... => .
#.##. => #
##.#. => .
##..# => .
.###. => .
.#.## => .
.#..# => #
..... => .
###.. => #
#..## => .
##.## => .
#.... => .
...## => #
....# => .
#.#.# => #
###.# => .
.#### => #
.#... => #
#.### => .
..### => .
.#.#. => #
.##.# => .
#..#. => #
...#. => .
#...# => #
..##. => .
##### => #
..#.# => #`

var ex = `#..#.#..##......###...###
...## => #
..#.. => #
.#... => #
.#.#. => #
.#.## => #
.##.. => #
.#### => #
#.#.# => #
#.### => #
##.#. => #
##.## => #
###.. => #
###.# => #
####. => #`

func plantGet(s string, i int) byte {
	if i < 0 || i >= len(s) {
		return '.'
	}
	return s[i]
}

func transform(rules map[string]bool, left int, s string) (int, string) {
	var b strings.Builder
	key := ""
	for i := -4; i < len(s)+5; i++ {
		if len(key) == 5 {
			key = key[1:]
		}
		key = key + string(plantGet(s, i))
		if rules[key] {
			b.WriteByte('#')
		} else {
			b.WriteByte('.')
		}
	}
	r := b.String()
	left -= 6
	for i := 0; i < len(s); i++ {
		if r[i] != '.' {
			break
		}
		left++
	}
	r = strings.Trim(r, ".")
	return left, r
}

func main() {
	lines := strings.Split(input, "\n")
	is := lines[0]
	rules := map[string]bool{}
	for _, line := range lines[1:] {
		rules[line[:5]] = (line[9] == '#')
	}

	doiters := func(left int, x string, iters int) int {
		for step := 0; step < iters; step++ {
			left, x = transform(rules, left, x)
		}
		sum := 0
		for i, s := range x {
			if s == '#' {
				sum += i + left
			}
		}
		return sum
	}
	fmt.Println(doiters(0, is, 20))

	x := is
	left := 0
	var step int
	got := map[string][2]int{}
	const N = 10000
	for step = 0; step < N; step++ {
		if _, ok := got[x]; ok {
			break
		}
		got[x] = [2]int{left, step}
		left, x = transform(rules, left, x)
	}
	if step == N {
		panic("not found")
	}

	left0, step0 := got[x][0], got[x][1]
	iters := (50000000000 - step0)
	loops := iters / (step - step0)
	iters %= step - step0
	left1 := left0 + (left-left0)*loops

	fmt.Println(doiters(left1, x, iters))
}
