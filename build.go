package main

import (
	"bufio"
	"fmt"
	"github.com/gogotchuri/FiniteAutomata/simulation"
	"os"
	"strings"
)

func main() {
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.Trim(input, " \n")
	fa := simulation.BuildAutomataFromRegex(input)
	fmt.Println(fa)
}
