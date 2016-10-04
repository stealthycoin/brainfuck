package main

import (
	"os"
	"fmt"
	"flag"
	"io/ioutil"
)

var debug = flag.Bool("d", false, "Enables web debugger on port 8888")


//
// Parse flags and kick off interpreter
//
func main() {
	flag.Parse()

	program, err := ioutil.ReadFile(flag.Args()[0])
//	beautiful := beautify(string(program))

	if err != nil {
		fmt.Println("Error reading file", err)
	} else {
		if *debug {
			// TODO: something
		}

		brainfuck(string(program))
	}
}


//
// Takes a program string and beautifies it
// TODO: clean this up, ugly
//
func beautify(program string) string {
	tab := "    "
	out := ""
	tab_level := ""

	for _,char := range program {
		if char == '[' {
			out += "\n" + tab_level + "[\n"
			tab_level += tab
			out += tab_level
		} else if char == ']' {
			tab_level = tab_level[:len(tab_level)-len(tab)]
			out += "\n" + tab_level + "]\n"
		} else {
			out += string(char)
		}
	}

	return out
}


//
// Interpret and run a brainfuck program
//
func brainfuck(program string) {
	// Housekeeping variables
	jump_forward := make(map[int]int)
	jump_backward := make(map[int]int)
	stack := make([]int, 0)

	// Preprocess program
	for idx, symbol := range program {
		switch symbol {
		case '[':
			stack = append(stack, idx)
		case ']':
			if len(stack) == 0 {
				fmt.Printf("Mismatched ] at index: %d\n", idx)
				os.Exit(1)
			}

			// Adjust stack and tables
			end := len(stack)
			open := stack[end-1]
			stack = stack[:end-1]

			jump_forward[open] = idx
			jump_backward[idx] = open-1
		}
	}

	failed := false
	for len(stack) > 0 {
		fmt.Printf("Mismatching [ at index: %d\n", stack[0])
		stack = stack[1:]
		failed = true
	}
	if failed {
		os.Exit(1)
	}


	// Run actual progarm
	mem := make([]uint8, 30000)
	ptr := 0
	program_len := len(program)


	for program_counter := 0 ; program_counter < program_len ; program_counter++ {

		switch program[program_counter] {
		case '+':
			mem[ptr] += 1

		case '-':
			mem[ptr] -= 1

		case '<':
			ptr -= 1
			if ptr < 0 {
				fmt.Println("pointer out of bounds to the left")
			}

		case '>':
			ptr += 1
			if ptr >= len(mem) {
				mem = append(mem, 0)
			}

		case '.':
			fmt.Printf("%s", string(mem[ptr]))

		case ',':
			b:= make([]byte, 1)
			os.Stdin.Read(b)
			mem[ptr] = b[0]

		case '[':
			if mem[ptr] == 0 {
				program_counter = jump_forward[program_counter]
			}

		case ']':
			if mem[ptr] != 0 {
				program_counter = jump_backward[program_counter]
			}
		}
	}
}
