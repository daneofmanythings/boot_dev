package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const PROMPT string = "pokedex > "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		args := strings.Split(line, " ")

		COMMANDS["help"] = help
		if command, ok := COMMANDS[args[0]]; !ok {
			fmt.Printf("Command '%s' not recognized\n", line)
		} else {
			err := command.Run(args)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
