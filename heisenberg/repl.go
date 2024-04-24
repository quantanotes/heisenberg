package heisenberg

import (
	"bufio"
	"fmt"
	"heisenberg/parser"
	"os"
	"strings"
)

func (h *Heisenberg) repl() {
	fmt.Println("Welcome to Heisenberg REPL.")
	fmt.Println("Use \\h for help and \\q to quit.")
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}

	for {
		if len(lines) == 0 {
			fmt.Print("> ")
		} else {
			fmt.Print("  ")
		}

		scanner.Scan()
		input := scanner.Text()

		if input == "\\q" || input == "\\quit" {
			break
		}

		if strings.HasPrefix(input, "\\") {
			h.processCmds(input)
			continue
		}

		lines = append(lines, input)
		if strings.HasSuffix(input, ";") || len(input) == 0 {
			h.eval(strings.Join(lines, "\n"))
			lines = []string{}
		}
	}

	fmt.Println("Goodbye!")
}

func (h *Heisenberg) eval(input string) {
	ast, err := parser.Parse(input)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	for _, a := range ast {
		fmt.Printf("%v\n", a)
	}
}

func (h *Heisenberg) processCmds(input string) {

}
