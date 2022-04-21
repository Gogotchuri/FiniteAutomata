package main

import (
	"fmt"
	"github.com/gogotchuri/FiniteAutomata/simulation"
	"os"
)

func main() {
	var word string
	if _, err := fmt.Scanf("%s\n", &word); err != nil {
		panic(err)
	}
	fa, err := simulation.BuildAutomataFromDescriptionSTDIn(os.Stdin)
	if err != nil {
		panic(err)
	}
	fa.Simulate(word)
}
