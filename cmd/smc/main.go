package main

import (
	"fmt"
	"os"

	"github.com/geisonbiazus/smc/internal/smc"
)

func main() {
	compiler := smc.NewCompiler(os.Stdin, os.Stdout)
	err := compiler.Compile()

	if err != nil {
		for _, e := range compiler.Errors {
			fmt.Println(e.String())
		}
	}
}
