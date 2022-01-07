package main

import (
	"fmt"
	"sort"
	"strings"
)

type cart struct {
	x, y      int
	dx, dy    int
	turnstate int
}

func cartspeed(b byte) (int, int) {
	if b == '>' {
		return 1, 0
	} else if b == '<' {
		return -1, 0
	} else if b == '^' {
		return 0, -1
	} else if b == 'v' {
		return 0, 1
	}
	return 0, 0
}

// dostep moves the carts one step forward, returning a slice of the new carts,
// and all collision points this tick.
func dostep(grid [][]byte, cs []cart) ([]cart, [][2]int) {
	var colls [][2]int
	sort.Slice(cs, func(i, j int) bool {
		if cs[i].y != cs[j].y {
			return cs[i].y < cs[j].y
		}
		return cs[i].x < cs[j].x
	})
	var rm uint32
	for i := range cs {
		c := &cs[i]
		if (rm>>i)&1 == 1 {
			continue
		}
		c.x += c.dx
		c.y += c.dy
		collided := false
		for j, d := range cs {
			if i == j || (rm>>j)&1 == 1 {
				continue
			}
			if c.x == d.x && c.y == d.y {
				collided = true
				colls = append(colls, [2]int{c.x, c.y})
				rm |= (1 << i)
				rm |= (1 << j)
			}
		}
		if collided {
			continue
		}
		g := grid[c.y][c.x]
		if g == '/' {
			c.dx, c.dy = -c.dy, -c.dx
		} else if g == '\\' {
			c.dx, c.dy = c.dy, c.dx
		} else if g == '+' {
			if c.turnstate == 0 {
				c.dx, c.dy = c.dy, -c.dx
			} else if c.turnstate == 2 {
				c.dx, c.dy = -c.dy, c.dx
			}
			c.turnstate = (c.turnstate + 1) % 3
		} else if g == '-' {
			if c.dy != 0 {
				panic("-")
			}
		} else if g == '|' {
			if c.dx != 0 {
				panic("|")
			}
		} else {
			panic(fmt.Sprintf("offgrid? %c", g))
		}
	}
	var ncs []cart
	for i, c := range cs {
		if (rm>>i)&1 == 0 {
			ncs = append(ncs, c)
		}
	}
	return ncs, colls
}

func (c *cart) byte() byte {
	if c.dx > 0 {
		return '>'
	} else if c.dx < 0 {
		return '<'
	} else if c.dy > 0 {
		return 'v'
	} else if c.dy < 0 {
		return '^'
	} else {
		panic("bad cart speed")
	}
}

func findcart(cs []cart, i, j int) (byte, bool) {
	for _, c := range cs {
		if c.x == i && c.y == j {
			return c.byte(), true
		}
	}
	return 0, false
}

func showBoard(grid [][]byte, cs []cart) {
	for j, row := range grid {
		for i, b := range row {
			if cb, ok := findcart(cs, i, j); ok {
				fmt.Printf("%c", cb)
			} else {
				fmt.Printf("%c", b)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	lines := strings.Split(input, "\n")
	W := len(lines[0])
	H := len(lines[1])
	_, _ = W, H
	grid := [][]byte{}
	carts := []cart{}
	for row, line := range lines {
		g := []byte(line)
		for col := range g {
			b := g[col]
			dx, dy := cartspeed(b)
			if dx+dy != 0 {
				carts = append(carts, cart{
					x:  col,
					y:  row,
					dx: dx,
					dy: dy,
				})
				if dx != 0 {
					g[col] = '-'
				} else {
					g[col] = '|'
				}
			}
		}
		grid = append(grid, g)
	}
	for _, row := range grid {
		fmt.Println(string(row))
	}
	fmt.Println(carts)
	p1 := false
	for step := 0; step < 200000; step++ {
		// showBoard(grid, carts)
		var colls [][2]int
		carts, colls = dostep(grid, carts)
		if !p1 && len(colls) > 0 {
			fmt.Printf("%d,%d\n", colls[0][0], colls[0][1])
			p1 = true
		}
		if len(carts) == 1 {
			fmt.Printf("%d,%d\n", carts[0].x, carts[0].y)
			break
		}
	}
}

var input = `            /-------------------------------------------------------------------------------------\               /-------------------------------\   
  /---------+---\                                                                   /-------------+---------------+---\              /-----------\|   
  |       /-+---+-----------------------------------------------\                 /-+-------------+---------------+---+----\         |           ||   
  |       | |   |                /----------\      /------------+------------\    | |/------------+---------------+---+----+---------+------\    ||   
 /+-------+-+---+----------------+----------+--\   |  /---------+------------+----+-++------------+---------------+---+----+---------+--\   |    ||   
 ||       | |   |                |          |  |   |  |         |           /+----+-++------------+---------------+---+----+---------+--+---+\   ||   
 ||       | |   | /--------------+----------+--+---+--+---------+-----------++----+\||            |         /-----+---+----+---\     |  |   ||   ||   
 ||       | |/--+-+------\       |          |  | /-+--+---------+-----------++----++++------------+---------+-----+---+----+---+-----+--+---++---++\  
 ||       | ||  | |      |       |   /------+--+-+-+--+---------+-----------++----++++------------+------\  |     |   |    |   |     |  |   ||   |||  
 ||       | ||  | |      |       |/--+-----\|  | | |  |         |           ||    ||||            |      |  |     |   |    |   |     |  |   ||   |||  
 ||       | ||  | |      |       ||  |     ||  | | |  |         |           ||    ||||            |      |  |     |   | /--+---+-----+--+---++---+++\ 
 ||       | ||  | |      |       ||  |     ||  | |/+--+---------+-----------++----++++------------+------+--+-----+---+-+--+---+-----+--+---++-\ |||| 
 ||   /---+-++--+-+------+-------++--+-----++--+-+++--+---------+---\       ||    ||||            |      |  |     |   | |  |   |     \--+---++-+-/||| 
 ||   |   \-++--+-+------+-------++--+-----++--+-+++--+---------/   |       ||    ||||            |      |  |     |   | |  |   |        |   || |  ||| 
 ||  /+-----++--+-+------+-------++--+-----++--+-+++--+-------------+-------++----++++-----------\|      |  |     |   | |  |   |        |   || |  ||| 
 ||  ||     ||  | |      |   /---++--+-----++--+-+++--+-------------+-------++----++++-----------++-----\|  |     |   | |  | /-+-------\|   || |  ||| 
 ||  ||     ||  | |   /--+---+---++--+-----++--+-+++--+-------------+-------++----++++-----------++-----++--+-----+\  | |  | | |       ||   || |  ||| 
 ||  || /---++--+-+---+--+---+---++--+-----++--+-+++--+-------------+-------++--\ ||||           ||     ||  |     ||  | |  | | |       ||   || |  ||| 
 ||  || |/--++--+-+---+--+---+---++--+-----++--+-+++-\|    /--------+-------++--+-++++-----------++->---++--+-----++--+-+--+-+-+---\   ||   || |  ||| 
 ||  || ||  ||  | |   |  |   |   \+--+-----+/  | ||| ||    |        |       ||  | ||||       /---++-----++--+-----++--+-+--+-+\|   |   ||   || |  ||| 
 \+--++-++--++--+-+---+--+---+----+--+-----+---/ ||| ||    |        |       ||/-+-++++-------+---++-----++--+-----++--+-+--+-+++---+-\ ||   || |  ||| 
  | /++-++--++--+-+---+--+---+----+--+-----+-----+++-++----+-------\|       ||| | ||||       |   ||     ||  |     ||  | |  | |||   | | ||   || |  ||| 
/-+-+++-++--++--+-+---+--+---+----+--+-----+-----+++-++----+-------++-------+++-+-++++------\|   ||     ||  |     ||  | \--+-+++---+-+-++---++-+--++/ 
| | ||| ||  |\--+-+---+--/   |  /-+--+---\ |     ||| ||    |       ||      /+++-+-++++------++---++-----++--+-----++--+---\| |||   | | ||   || |  ||  
| | ||| ||  |   | |   |     /+--+-+--+---+-+-----+++-++----+--\  /-++------++++-+-++++------++---++-----++--+-----++--+---++-+++---+>+-++--\|| |  ||  
| | ||| || /+---+-+---+-----++--+-+--+---+-+-----+++-++----+--+--+-++---->-++++-+-++++---\  ||   ||     ||  |     ||  |   || |||   | | ||  ||| |  ||  
| | ||| || ||   | |   |     ||  | |  |   | |   /-+++-++----+--+--+-++------++++-+-++++---+--++---++-----++--+-----++--+---++-+++\  | | ||  ||| |  ||  
| | ||| || ||   | |   |     ||  | |  |   | |   | ||| ||    |  |  | ||      |||| | ||||   |  ||   ||     ||  |     ||  |   || ||||  | | ||  ||| |  ||  
|/+-+++-++-++---+-+---+-----++--+-+--+---+-+---+-+++-++----+--+--+-++-\    |||| | ||||   |  ||   ||     ||  |     ||  |   || ||||  | | ||  ||| |  ||  
||| ||| || ||   |/+---+-----++--+-+--+---+-+---+-+++-++----+--+--+-++-+----++++-+-++++---+--++---++-----++--+-----++--+-\ || ||||  | | ||  ||| |  ||  
||| ||| || || /-+++---+-----++--+-+-\|   | |  /+-+++-++----+--+--+-++-+----++++-+-++++---+--++---++-----++--+-----++--+\| || ||||  | | ||  ||| |  ||  
||| ||| || || | |||   \-----++--+-+-++---+-+--++-+++-++----+--+--+-++-+----++++-+-++++---+--++---++-----++--+-----+/  ||| || ||||  | | ||  ||| |  ||  
||| ||| || || | |||         ||  | \-++---+-/  || ||\-++----+--+--+-++-+----++/| | |||| /-+--++---++-----++--+-----+---+++-++-++++--+-+-++--+++-+--++-\
||| ||| || || | |||         ||  |/--++---+----++-++--++----+--+--+-++-+----++-+-+\|||| | |  ||   ||     ||  |     |   ||| || ||||  | | ||  ||| |  || |
||| ||| ||/++-+-+++---------++-\||  ||   |    || ||/-++----+--+--+-++-+----++-+-++++++-+-+--++---++-----++--+-----+---+++-++-++++--+-+-++--+++\|  || |
||| ||| ||||| | |||         || |||  ||   |    || ||| ||    |  |  | || |    \+-+-++++++-+-+--++---++-----++--+-----+---+++-/| ||||  | | ||  |||||  || |
||| ||| ||||| | |||       /-++-+++--++---+----++-+++-++----+--+--+-++-+-----+-+-++++++-+-+--++---++-----++--+-----+---+++--+-++++--+-+\||  |||||  || |
||| ||| ||||| | |||       | || |||  ||/--+----++-+++-++----+--+--+-++-+-----+-+-++++++-+-+--++---++-----++--+-----+---+++--+-++++--+-++++\ |||||  || |
||| ||| ||||| | |||       | || |||  |||  |    |\-+++-++----+--+--+-++-+-----+-+-++++++-+-+--++---++-----++--+-----+---+++--+-+++/  | ||||| |||||  || |
||| ||| ||||| | |||       | |\-+++--+++--+----+--+++-++----+--+--+-++-+-----+-+-++++++-+-+--++---++-----/|  |     |   |||  | |||   | ||||| |||||  || |
|||/+++-+++++-+-+++-----\ | |  |||  ||| /+----+--+++-++----+--+--+-++-+-----+-+-++++++-+-+--++---++------+--+-\   |   |||  | |||   | ||||| |||||  || |
||||||| ||||| | |||  /--+-+-+--+++--+++-++----+--+++-++----+--+--+-++-+-----+-+-++++++-+-+--++---++------+--+-+---+---+++--+-+++-\ | ||||| |||||  || |
||||||| ||||| | |||  |  | | |  |||  ||| ||    |  ||| ||  /-+--+--+-++-+-----+-+-++++++-+-+--++---++------+--+-+\  |   |||  | ||| | | ||||| |||||  || |
||||||| \++++-+-+++--+--+-+-+--+++--+++-++----+--+++-++--+-+--+--+-++-+-----+-+-/||||| v |  ||/--++------+--+-++--+---+++--+-+++-+-+-+++++\|||||  || |
|||||||  |||| | |||  |  | | |  |||  ||| ||    |  ||| ||  | |  |  | || |     | |  ||||| | |  |||  ||      |  | ||  |   |||  | \++-+-+-++/||||||||  || |
|||||||  |||| | |||  |  | | |  |||  ||| ||    |  ||| ||  | |  |  | || |     | |  ||||| | |  |||  ||      |  | ||  |   |||  |  || | | || ||||||||  || |
|||||||  |||| | |||  |  | | |  |||  ||| ||    |  ||| ||  | |  |  | || |     | |/-+++++-+-+--+++--++------+--+-++--+---+++--+--++-+-+-++\||||||||  || |
|||||||  |||| | |||  |  | | |  |||  ||| ||    |  |||/++--+-+--+--+-++-+-----+-++-+++++-+-+--+++--++------+--+-++--+---+++--+\ || | | |||||||||||  || |
|||||||  |||| | |||/-+--+-+-+--+++--+++-++----+--++++++--+-+--+--+-++-+-----+-++-+++++-+-+--+++--++------+--+-++--+---+++--++\|| | | ||||||||^||  || |
|||||||  |||| | |||| |  | | |  |||  ||| ||    |  ||||||  | |  |  | || |     | || ||||| | |  |||  ||      |  | ||  |   |||  ||||| | | |||||||||||  || |
|||||||  |||| | |||| |  | | |  ||| /+++-++----+--++++++--+-+--+--+-++-+-----+-++-+++++-+-+--+++--++------+--+-++--+\  |||  ||||| | | |||||||||||  || |
|||||||  |||| | |||| |  | | |  ||| |||| ||   /+--++++++--+-+--+--+-++-+-----+-++-+++++-+-+--+++--++-----\|  | ||  ||  ||| /+++++-+-+\|||||||||||  || |
|||||||  |||| | |||| |  | | |  ||| |||| ||   ||  ||||||  | |  |  | || |     | || ||||| | |  ||\--++-----++--+-++--++--+++-++++++-+-+++++++/|||||  || |
|||||||  |||| | |||| |  | | |  ||| |||| ||   ||  ||||||  | |  |  | || |     | || ||||| | |  ||   ||     ||  | ||  ||  ||| |||||| | ||||||| |||||  || |
|||||||  |||| | |||| |  | | |  ||| |||| ||   ||  ||||||  | \--+--+-++-+-----+-++-+++++-+-+--++---++-----++--+-++--++--+++-++++++-+-/|||||| |||||  || |
|||||||/-++++-+-++++-+--+-+-+--+++-++++-++---++-\||||||  |    |  | || |     | || ||||| | |  ||   ||     ||  | ||  ||  ||| |||||| |  |||||| |||||  || |
|||||||| |||| | |||| |  | | |  ||| ||\+-++---++-+++++++--+----+--+-++-+-----+-++-+++++-+-+--++---++-----+/  | ||  ||  ||| |||||| |  |||||| |||||  || |
|||||||| |||| | |||| |  | | \--+++-++-+-++---++-+++++++--+----/  | || |     | || |\+++-+-+--++---++-----+---+-++--++--+++-+/|||| |  |||||| |||||  || |
|||||||| |||| | |||| |  | |    |||/++-+-++---++-+++++++--+-------+-++-+-----+-++-+-+++-+-+--++---++\    |   | ||  ||  ||| | |||| |  |||||| |||||  || |
|||||||| |||| | |||| |  | |    |||||| | ||   || |||||||  | /-----+-++\|     | || | ||| | |  ||   |||    |   | ||  ||  ||| | |||| |  ^||||| |||||  || |
|||||||| |||| | |||| |  | |    |||||| | ||   \+-+++++++--+-+-----+-++++-----+-++-+-+++-+-+--++---+++----/   | ||  ||  ||| | |||| |  |||||| |||||  || |
|||||||| |||| | |||| |  | |  /-++++++-+-++----+-+++++++--+-+-----+-++++-----+\|| | ||| | |  |\---+++--------+-++--++--+++-+-++/| |  |||||| |||||  || |
|||||||| |||| | |||| |  | |  | |||||| | ||    | |||||||  | |     | ||||     |||| | ||| | |  |    |||        | ||  ^|  ||| | || | |  |||||| |||||  || |
|||||||| |||| | |||| |  | \--+-++++++-+-++----+-+++++++--+-+-----+-++++-----++++-+-+++-+-+--+----+++--------+-++--++--+++-+-++-+-+--++/||| |||||  || |
|||||||| |||| | |||| |  |    |/++++++-+-++----+-+++++++--+-+-----+-++++-----++++-+-+++-+-+--+----+++--------+-++--++--+++-+-++-+-+--++-+++\|||||  || |
||||||||/++++-+-++++-+--+----++++++++-+-++----+-+++++++--+-+-----+-++++-----++++-+-+++-+-+--+----+++--------+-++\/++--+++-+-++-+-+--++-+++++++++-\|| |
||||||||||||| | |||| |/-+----++++++++-+-++----+-+++++++--+-+-----+-++++-----++++-+-+++-+-+--+----+++--------+-++++++--+++-+-++-+-+\ || ||||||||| ||| |
|||||||\+++++-+-++++-++-+----++++++++-+-++----+-/||||||  | |     | ||||   /-++++-+-+++-+-+\ |    |||        | ||||||  ||| | || | || || ||||||||| ||| |
|||||||/+++++-+-++++-++-+----++++++++-+-++----+--++++++--+-+-----+-++++--\| |||| | ||| | || |    |||        | ||||||  ||| | || | || || ||||||||| ||| |
||||||||||||| | |||| || |    ||||||\+-+-++----+--++++++--+-+-----+-++++--++-++++-+-+++-+-++-+----+++--------+-+++++/  ||| | || | || || ||||||||| ||| |
|||||||||||||/+-++++-++-+----++++++-+-+-++----+--++++++--+-+-----+-++++--++-++++-+-+++-+-++\|    |||        | |||||   ||| | || | || || ||||||||| ||| |
||||||||||||||| |||| || |    |||||| | | ||    |  ||||||  | |     |/++++--++-++++-+-+++-+-++++----+++\   /---+-+++++---+++-+-++-+\|| || ||||||||| ||| |
||||||||||||||| |||| || |    ||||\+-+-+-++----+--++++++--+-+-----++++++--++-++++-/ ||| | ||||    ||||   |   | |||||   ||| | || |||| || ||||||||| ||| |
||||||||||||||| |||| ||/+--<-++++-+-+-+-++----+--++++++--+-+-----++++++--++-++++--\||| | ||||    ||||   |   | |||||   ||| | || |||| || ||||||||| ||| |
||\++++++++++++-/||| ||||    \+++-+-+-+-++----+--++++++--+-+-----++++++--++-+/||/-++++-+-++++----++++---+---+-+++++---+++-+-++-++++-++-+++++++++\||| |
|| ||||||||||||  |||/++++-----+++-+-+-+-++----+--++++++--+-+-----++++++--++-+-+++-++++-+-++++----++++---+---+-+++++-\ ||| | || |||| || ||||||||||||| |
|\-++++++++++++--++++++++-----+++-+-+-+-++----+--++++++--+-+-----+++++/  || | ||| |||| | ||||    ||||   |   \-+++++-+-+++-+-++-/||| || ||||||||||||| |
|/-++++++++++++--++++++++-----+++-+-+-+-++----+--++++++--+-+-----+++++---++-+-+++-++++-+-++++----++++---+----\||||| | ||| | ||  ||| || ||||||||||||| |
|| ||||||||||||  \+++++++-----+++-+-+-+-++----+--++++++--+-+-----+++++---++-+-+++-++++-+-++++----++++---+----++++++-+-++/ | ||  ||| || ||||||||||||| |
|| |||||||||\++---+++++++-----+++-+-+-+-++----+--++++++--+-+-----+++++---++-+-+++-++++-+-++++----+/||   |    |||||| | ||  | ||  ||| || ||||||||||||| |
|| ||||||||| ||  /+++++++-----+++-+-+\| ||    |  ||||||  | |     |||||   || | ||| |||| | ||||    | ||   |    ^||||| | ||  | ||  ||| || ||||||||||||| |
|| ||||||||| ||  ||||||||     ||| \-+++-++----+--++++++--+-+-----+++++---++-+-+++-++++-+-++++----+-/|   |    |||||| | ||  | ||  ||| || ||||||||||||| |
|| ||||||||| ||  ||||||||     ||\---+++-+/    |  ||||||  | |     |||||   || | |||/++++-+-++++----+--+---+----++++++-+-++--+-++--+++-++\||||||||||||| |
|| ||||||||| ||  ||||||||     ||    ||| |     |  |||||\--+-+-----+++++---++-+-++++++++-+-++++----+--+---+----++++++-+-++--+-++--+++-++++/||||||||||| |
|| ||||||||| ||  ||||||||     ||    ||| |/----+--+++++---+-+-----+++++-\ || | |||||||| | ||||    |  |   |    |||||| | ||  | ||  ||| |||| ||||||||||| |
|| ||||||||| |v  ||||||\+-----++----+++-++----+--+++++---+-+-----+++++-+-++-+-++++/||| |/++++----+--+---+----++++++-+-++--+-++-\||| |||| ||||||||||| |
|| |||\+++++-++--++++++-+-----++----+++-++----+--+++++---+-+-----+++/| | || | |||| ||| ||||||    |  |   |    |||||| | ||  | || |||| |||| ||||||||||| |
|| ||| ||||| ||  |||||| |     ||    ||| ||    |  |||||   | |     ||| | | || | |\++-+++-++++++----+--+---+----++++++-+-++--+-++-++++-+++/ ||||||||||| |
|| ||| ||||| ||  |||||| |     ||    ||| ||    |  |||||   | |     ||| | | || | | || ||| ||||||    |  |   |    |||||| | ||  | || |||| |||  ||||||||||| |
|| ||| ||||| ||  |||||| | /---++----+++-++---\|  |||||   | |     ||| | | || | | || ||\>++++++----+--+---+----++++++-+-++--+-++-++++-+++--+++/||||||| |
|| ||| ||||| ||  |||||| | |   ||    ||| ||   ||  |||||   | |     ||| | | || | | || |\-<++++++----+--+---+----++++++-+-/|  | || |||| |||  ||| ||||||| |
|| ||| ||||| ||  |||||| | |   ||    ||| ||   ||  |||||   | |     ||| | | || | | \+-+---++++++----+--+---+----++++++-+--+--+-++-++++-+++--+++-+++/||| |
||/+++-+++++-++--++++++-+-+---++-\  ||| ||   ||  |||\+---+-+-----+++-+-+-++-+-+--+-+---++++++----+--+---+----++++++-+>-+--+-/| |||| |||  ||| ||| ||| |
|\++++-+++++-++--++++++-+-+---++-+--+++-++---++--+++-+---+-+-----+++-+-+-++-+-+--+-+---++++++----+--+---+----/||||| |  |  |  | |||| |||  ||| ||| ||| |
| |||| ||||| ||  |\++++-+-+---++-+--+++-++---++--+++-+---+-+-----+++-+-+-++-+-+--+-/   ||||||    |  |  /+-----+++++-+--+--+--+-++++-+++\ ||| ||| ||| |
| |||| ||||| ||/-+-++++-+-+---++\|  ||| ||   ||  ||| |   | |     ||| | | || | |  |     ||||||    |  |  ||     ||||| |  |/-+--+-++++-++++-+++-+++-+++\|
| |||| ||||| ||| | |||| | |   ||||/-+++-++---++--+++-+---+-+----\\++-+-+-++-+-+--+-----++++++----+--+--++-----+++++-+--++-+--+-++++-++++-++/ ||| |||||
| |||| ||||| ||| | |||| | |   ||||| ||| ||   ||  \++-+---+-+----+-++-+-+-++-+-+--+-----++++++----+--+--++-----+++++-+--++-+--+-++++-++++-++--+++-++/||
| |||| ||||| ||| | |||\-+-+---+++++-+++-++---++---++-+---+-+----+-++-+-+-++-+-+--+-----++++++----+--+--++-----+++++-+--++-+--+-+++/ |||| ||  ||| || ||
| ||\+-+++++-+++-+-+++--+-+---+++++-+++-++---++---++-+---+-+----+-+/ | | ||/+-+--+-----++++++----+--+--++-----+++++-+-\|| |  | |||  |||| ||  ||| || ||
| || | ||||| ||| | |||  | |/--+++++-+++-++---++---++-+---+-+--\ | |  | | |||| |  |     ||||||   /+--+--++-----+++++-+-+++-+--+-+++--++++-++--+++-++\||
| || | ||||| ||| | |||  | ||  ||||| ||| ||   ||   || | /-+-+--+-+-+--+-+-++++-+--+-----++++++---++--+--++-----+++++-+-+++-+\ | |||  |||| ||  ||| |||||
| || | |||||/+++-+-+++--+-++--+++++-+++-++---++---++-+-+-+-+--+-+-+--+-+-++++-+--+-----++++++---++--+--++-\   ||||| | ^|| || |/+++--++++-++-\||| |||||
| || | ||||||||| | |||  | ||  ||||| ||\-++---++---++-+-+-+-+--+-+-+--+-+-++++-+--+-----++++++---++--+--++-+---+++++-+-+++-++-+++++--++++-/| |||| |||||
| || | ||||||||| \-+++--+-++--+++++-+/  ||   ||   || | | | |  | | |  | | |||\-+--+-----++++++---++--+--++-+---+++++-+-+++-++-+++++--++++--+-+/|| |||||
| || | \++++++++---+++--+-++--+++++-+---++---++---++-+-+-+-+--+-+-+--+-+-/||  |  |     ||||||   ||  |  || |   ||||| | ||| || |||||  ||||  | | || |||||
| || |  ||||||||   |||  | ||  ||||| |   ||   ||   || | | | |  | | |  | |  ||  |  |     ||||||   ||  |  || |   ||||| | ||| || |||||  ||||  | | || |||||
\-++-+--++++++++---+++--+-++--+++++-+---++---++---++-+-+-+-+--+-+-+--+-+--++--+--+-----+++++/   ||  |  || |   ||||\-+-+++-++-+++++--++++--+-+-++-+/|||
  || |  ||||||||   |||  | ||  ||||| |   ||   ||   || | | | |  | | |  | |  ||  |  |     |\+++----++--+--++-+---++++--+-+++-++-++/|v  ||||  | | || | |||
  || |  ||||||||   |||  | ||  ||||| |   ||   ||   |\-+-+-+-+--+-+-+--+-+--++--+--+-----+-+++----++--+--++-+---++++--+-+++-++-++-++--++++--+-+-/| | |||
  || |  ||||||||   |||  | ||  ||||| |   \+---++---+--+-+-+-+--+-+-+--+-+--++--+--+-----+-+++----++--+--++-+---/|||  | ||| || || ||  ||||  | |  | | |||
  || |  ||||\+++---+++--+-++--+++++-+----+---++---+--+-+-+-+--+-+-+--+-+--++--+--+-----+-+++----++--+--++-/    |||  | ||| || || ||  ||||  | |  | | |||
  || |  |||| |||   ||\--+-++--+++++-+----+--<++---+--+-+-+-+--+-+-+--+-+--++--+--+-----+-+++----++--+--++------+++--+-+++-++-++-+/  ||||  | |  | | |||
  || |  |||\-+++---++---+-++--+++++-+----+---++---+--+-+-+-+--+-+-+--+-+--++--+--+-----+-/||    \+--+--++------+++--+-+++-++-++-+---++++--+-+--+-+-/||
  || |  |||  |||   ||   | ||  ||||| |    |   ||   |  | | | |  | | |  | |  ||  |  |     |  ||     |  |  ||      |||  | ||| || |\-+---++++--+-/  | |  ||
  || |  |||  |||   ||   | ||  \++++-+----+---++---+--+-+-+-+--+-+-+--+-+--++--+--+-----+--++-----+--+--++------+++--+-+++-++-+--+---++++--/    | |  ||
  |\-+--+++--+++---++---/ ||   |||| |    |   ||   |  | |/+-+--+-+-+--+-+--++--+--+-----+\ ||     |  |  ||      |||  | ||| || |  |   ||||       | |  ||
  |  |  |||  |||   ||     ||   |||| |    |   ||   \--+-+++-+--+-+-+--+-+--++--+--+-----++-++-----+--+--++------+++--+-+++-++-+--+---++++-------/ |  ||
  |  |  |||  |||   ||     ||   |||| |    |/--++------+-+++-+--+-+-+--+>+--++--+--+-----++-++-----+--+--++------+++--+-+++-++-+--+--\||||         |  ||
  |/-+--+++--+++---++-----++---++++-+----++--++\     | ||| |  | | |  | |  ||  |  |     || ||     |  |  \+------+++--+-+++-++-+--+--++++/         |  ||
  || \--+++--+++---++-----++---++++-+----++--+++-----+-+++-+--+-+-+--+-+--++--+--+-----++-++-----/  |   |      |||  | ||| || |  |  ||||          |  ||
  ||  /-+++--+++---++-----++---++++-+----++--+++-----+-+++-+--+-+-+--+-+--++--+--+---\ || ||        |   |      ||\--+-+++-++-+--+--++++----------/  ||
  ||  | |||  |||   ||     ||   |||| |/---++--+++-----+-+++-+--+-+-+--+-+--++--+--+---+-++-++--------+---+------++---+-+++-++-+--+-\||||             ||
  ||  | |\+--+++---++-----++---++++-++---++--+++-----/ ||| |  | | |  | |  ||  |  |   | || ||        |   |      ||   | ||| || |  | |||||             ||
  ||  | | |  |||   ||     ||   |||\-++---++--+++-------+++-+--+-/ |  | |  ||  |  \---+-++-++--------+---+------++---+-+++-++-+--+-++++/             ||
  ||  | | |  |||   ||     \+---+++--++---++--/||       ||\-+--+---+--+-+--++--+------+-++-++--------+---+------/|   | ||| || |  | ||||              ||
  ||  | | |  |||   ||      |   |||  ||   ||   ||       ||  |  |   |  | |  ||  |      | \+-++--------+---+-------+---+-+++-++-+--+-++++--------------+/
  ||  | | |  |||   ||      |   |||  ||   ||   ||       ||  |  |   |  | |  ||  |      |  | ||       /+---+-------+---+-+++\|| |  | ||||              | 
  ||  | | |  ||\---++------+---+/|  ||   ||   ||       ||  |  |   |  | |  ||  \------+--+-++-------++---+-------+---+-++++++-+--+-+++/              | 
  ||  | | |  ||    ||     /+---+-+--++---++---++-------++--+--+---+--+-+--++---------+--+-++-------++---+--\    |   | |||||| |  | |||               | 
  ||  | | |  ||    \+-----++---+-+--++---++---++-------++--+--+---+--+-+--++---------+--+-++-------++---+--+----+---+-++++++-/  | |||               | 
  ||  | | |  ||     |     ||   | |  ||   ||   ||       |\--+--+---+--+-+--++---------+--/ ||       ||   |  |    |   | ||||||    | |||               | 
 /++--+-+-+--++-----+-----++---+-+--++---++---++-------+--\|  |   |  | |  ||         |    ||       ||   |  |    |   | ||||||    | |||               | 
 |||  | | |  |\-----+-----++---+-+--/|   ||   ||       |  ||  |   |  | |  \+---------+----/|       ||   |  |    |   | ||||||    | |||               | 
 |||  | | |  |      |     ||   | |   |   |\---++-------+--++--+---+--+-+---+---------+-----+-------++---+--+----+---+-++++++----+-+/|               | 
 |||  | | |  |      |     ||   | |   |   |    \+-------+--++--+---+--+-+---+---------+-----+-------++---+--+----+---+-+/||||    | | |               | 
 |\+--+-+-+--+------+-----++---+-/   \---+-----+-------+--++--+---+--+-+---+---------+-----+-------++---+--+----+---+-+-++++----+-/ |               | 
 | |  | | \--+------+-----++---/         |     |       |  ||  |   \--+-+---+---------+-----+-------+/   |  |    |   | | ||||    |   |               | 
 | |  | |    |      \-----++-------------+-----+-------+--++--+------+-+---+---------+-----+-------+----+--+----+---/ | ||||    |   |               | 
 | |  | |    |            ||             |     |       \--++--+------+-+---+---------+-----+-------+----+--+----+-----+-+++/    |   |               | 
 | |  | |    |            ||             |     |          |\--+------/ |   |         |     |       |    |  |    |     | |||     |   |               | 
 | |  \-+----+------------++-------------+-----+----------+---+--------+---+---------/     |       |    |  |    |     | |||     |   |               | 
 | |    |    \------------++-------------+-----+----------+---+--------+---+---------------/       |    |  |    |     | \++-----+---+---------------/ 
 | |    |                 ||             \-----+----------+---+--------/   \-----------------------+----+--+----+-----/  ||     |   |                 
 | |    |                 \+-------------------+----------+---+------------------------------------+----+--/    |        ||     |   |                 
 | \----+------------------+-------------------/          |   |                                    \----+-------+--------/\-----+---/                 
 |      |                  \------------------------------+---/                                         |       |               |                     
 |      \-------------------------------------------------+---------------------------------------------+-------/               |                     
 \--------------------------------------------------------/                                             \-----------------------/                    `

var ex = `/->-\        
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/   `