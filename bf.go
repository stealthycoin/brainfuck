package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"sync"
	"net/http"
	"io/ioutil"
)

var debug = flag.Bool("d", false, "Enables web debugger on port 8888")
var wg sync.WaitGroup

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
			debug_server()
		}

		brainfuck(string(program))
	}

}


//
// Run debug server
//
func debug_server() {
	wg.Add(1)
	fmt.Println("Debugger launched, visit http://localhost:8888/ to begin.")

	home := func(w http.ResponseWriter, r *http.Request) {
		b, err := Asset("data/html/home.html")
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(b))
	}

	static := func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		b, err := Asset(r.URL.Path[1:])
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(b))
	}

	go func() {
		http.HandleFunc("/", home)
		http.HandleFunc("/data/", static)
		http.ListenAndServe(":8888", nil)
	}()
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
	// TODO: Generalize and pull out of this function
	// needs to be usable in a debugger, also its a little ugly
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

		// Wait for debugger
		wg.Wait()

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
				mem = append(mem, make([]byte, 30000)...)
			}

		case '.':
			fmt.Printf("%s", string(mem[ptr]))

		case ',':
			b := make([]byte, 1)
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
