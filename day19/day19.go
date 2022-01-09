package main

import (
	"fmt"
	"strings"
)

type opinfo struct {
	name       string
	regA, regB bool
}

type operation int

var (
	addr = opinfo{"addr", true, true}
	addi = opinfo{"addi", true, false}

	mulr = opinfo{"mulr", true, true}
	muli = opinfo{"muli", true, false}

	banr = opinfo{"banr", true, true}
	bani = opinfo{"bani", true, false}

	borr = opinfo{"borr", true, true}
	bori = opinfo{"bori", true, false}

	setr = opinfo{"setr", true, false}
	seti = opinfo{"seti", false, false}

	gtir = opinfo{"gtir", false, true}
	gtri = opinfo{"gtri", true, false}
	gtrr = opinfo{"gtrr", true, true}

	eqir = opinfo{"eqir", false, true}
	eqri = opinfo{"eqri", true, false}
	eqrr = opinfo{"eqrr", true, true}
)

var allops = []*opinfo{
	&addr, &addi, &mulr, &muli,
	&banr, &bani, &borr, &bori,
	&setr, &seti, &gtir, &gtri, &gtrr,
	&eqir, &eqri, &eqrr,
}

func get(isReg bool, r [6]int, v int) (int, error) {
	if isReg && (v < 0 || v >= len(r)) {
		return 0, fmt.Errorf("bad register %d", v)
	}
	if isReg {
		return r[v], nil
	}
	return v, nil
}

func (o *opinfo) execute(ins [4]int, r [6]int) ([6]int, error) {
	a, aok := get(o.regA, r, ins[1])
	b, bok := get(o.regB, r, ins[2])
	if aok != nil || bok != nil {
		return r, fmt.Errorf("%s(%d): error %v and %v", o.name, ins, aok, bok)
	}
	if bok != nil {
		return r, bok
	}
	switch o.name[0] {
	case 'a':
		r[ins[3]] = a + b
	case 'm':
		r[ins[3]] = a * b
	case 'b':
		if o.name[1] == 'a' {
			r[ins[3]] = a & b
		} else {
			r[ins[3]] = a | b
		}
	case 's':
		r[ins[3]] = a
	case 'g':
		if a > b {
			r[ins[3]] = 1
		} else {
			r[ins[3]] = 0
		}
	case 'e':
		if a == b {
			r[ins[3]] = 1
		} else {
			r[ins[3]] = 0
		}
	default:
		return r, fmt.Errorf("unknown op %s", o.name)
	}
	return r, nil
}

type program struct {
	ip  int
	ops [][4]int
}

func readProgram(s string) (program, error) {
	var prog program

	for _, line := range strings.Split(s, "\n") {
		var a, b, c int
		var opn string
		if _, err := fmt.Sscanf(line, "#ip %d", &a); err == nil {
			prog.ip = a
			continue
		}
		if _, err := fmt.Sscanf(line, "%s %d %d %d", &opn, &a, &b, &c); err != nil {
			return program{}, err
		}
		d := -1
		for i, o := range allops {
			if o.name == opn {
				d = i
				break
			}
		}
		if d == -1 {
			return program{}, fmt.Errorf("op %q not found", opn)
		}
		prog.ops = append(prog.ops, [4]int{d, a, b, c})
	}
	return prog, nil
}

func main() {
	prog, err := readProgram(input)
	if err != nil {
		panic(err)
	}

	for part := 1; part <= 2; part++ {
		pc := 0
		var regs [6]int
		if part == 2 {
			regs[0] = 1
		}
		for pc >= 0 && pc < len(prog.ops) {
			// fmt.Println(pc, regs, allops[prog.ops[pc][0]].name)
			regs[prog.ip] = pc
			regs, err = allops[prog.ops[pc][0]].execute(prog.ops[pc], regs)
			if err != nil {
				panic(err)
			}

			if regs[prog.ip] != pc {
				if pc == 11 {
					// I note that after setup, the program sums
					// the factors of reg[4] very slowly.
					// Just do it.
					x := regs[4]
					sum := 0
					for i := 1; i <= x; i++ {
						if x%i == 0 {
							sum += i
						}
					}
					regs[0] = sum
					break
				}
				pc = regs[prog.ip]
			}
			pc++
		}
		fmt.Println(regs[0])
	}
}

var ex = `#ip 0
seti 5 0 1
seti 6 0 2
addi 0 1 0
addr 1 2 3
setr 1 0 0
seti 8 0 4
seti 9 0 5`

var input = `#ip 3
addi 3 16 3
seti 1 6 1
seti 1 9 5
mulr 1 5 2
eqrr 2 4 2
addr 2 3 3
addi 3 1 3
addr 1 0 0
addi 5 1 5
gtrr 5 4 2
addr 3 2 3
seti 2 4 3
addi 1 1 1
gtrr 1 4 2
addr 2 3 3
seti 1 0 3
mulr 3 3 3
addi 4 2 4
mulr 4 4 4
mulr 3 4 4
muli 4 11 4
addi 2 5 2
mulr 2 3 2
addi 2 1 2
addr 4 2 4
addr 3 0 3
seti 0 3 3
setr 3 6 2
mulr 2 3 2
addr 3 2 2
mulr 3 2 2
muli 2 14 2
mulr 2 3 2
addr 4 2 4
seti 0 8 0
seti 0 8 3`
