package main

import (
	"fmt"
	"io/ioutil"
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

func get(isReg bool, r [4]int, v int) (int, bool) {
	if isReg && v < 0 || v > 3 {
		return 0, false
	}
	if isReg {
		return r[v], true
	}
	return v, true
}

func (o *opinfo) execute(ins [4]int, r [4]int) ([4]int, bool) {
	a, aok := get(o.regA, r, ins[1])
	b, bok := get(o.regB, r, ins[2])
	if !aok || !bok {
		return r, false
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
		panic("unknown op?")
	}
	return r, true
}

type example struct {
	before [4]int
	op     [4]int
	after  [4]int
}

func readExamples() ([]example, error) {
	bs, err := ioutil.ReadFile("./day16/day16.txt")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bs), "\n")
	var examples []example
	for i := 0; i < len(lines); i += 4 {
		var ex example
		var a, b, c, d int
		_, err := fmt.Sscanf(lines[i], "Before: [%d, %d, %d, %d]", &a, &b, &c, &d)
		if err != nil {
			return nil, fmt.Errorf("before: failed to parse %q: %v", lines[i], err)
		}
		ex.before = [4]int{a, b, c, d}

		_, err = fmt.Sscanf(lines[i+1], "%d %d %d %d", &a, &b, &c, &d)
		if err != nil {
			return nil, fmt.Errorf("op: failed to parse %q: %v", lines[i+1], err)
		}
		ex.op = [4]int{a, b, c, d}

		_, err = fmt.Sscanf(lines[i+2], "After: [%d, %d, %d, %d]", &a, &b, &c, &d)
		if err != nil {
			return nil, fmt.Errorf("after: failed to parse %q: %v", lines[i+2], err)
		}
		ex.after = [4]int{a, b, c, d}
		examples = append(examples, ex)
	}
	return examples, nil
}

func readProgram() ([][4]int, error) {
	bs, err := ioutil.ReadFile("./day16/day16_prog.txt")
	if err != nil {
		return nil, err
	}
	var r [][4]int
	for _, line := range strings.Split(string(bs), "\n") {
		var a, b, c, d int
		if _, err := fmt.Sscanf(line, "%d %d %d %d", &a, &b, &c, &d); err != nil {
			return nil, err
		}
		r = append(r, [4]int{a, b, c, d})
	}
	return r, nil
}

func getOpMap(canBe []map[*opinfo]bool, used uint32, i int) ([]*opinfo, bool) {
	if i == 16 {
		return nil, true
	}
	count := 0
	var r []*opinfo
	for j, op := range allops {
		if (used>>j)&1 == 1 {
			continue
		}
		if !canBe[i][op] {
			continue
		}
		om, ok := getOpMap(canBe, used|(1<<j), i+1)
		if ok {
			count++
			r = append([]*opinfo{op}, om...)
		}
	}
	return r, count == 1
}

func main() {
	exs, err := readExamples()
	if err != nil {
		panic(err)
	}
	sum := 0
	canBe := make([]map[*opinfo]bool, 16)
	for i := range canBe {
		canBe[i] = map[*opinfo]bool{}
		for _, op := range allops {
			canBe[i][op] = true
		}
	}
	for _, ex := range exs {
		valid := 0
		for _, ops := range allops {
			got, ok := ops.execute(ex.op, ex.before)
			if ok && got == ex.after {
				valid++
			} else {
				delete(canBe[ex.op[0]], ops)
			}
		}
		if valid >= 3 {
			sum++
		}
	}
	fmt.Println(sum)
	opmap, ok := getOpMap(canBe, 0, 0)
	if !ok {
		panic("opmap not found")
	}
	prog, err := readProgram()
	if err != nil {
		panic(fmt.Sprintf("program read failed: %v", err))
	}
	var r [4]int
	for _, ins := range prog {
		var ok bool
		r, ok = opmap[ins[0]].execute(ins, r)
		if !ok {
			panic("not ok in execution")
		}
	}
	fmt.Println(r[0])
}
