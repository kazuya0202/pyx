package main

import (
	"fmt"
	"os"

	"github.com/kazuya0202-dev/pyx/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// os.Exit(0)
}
