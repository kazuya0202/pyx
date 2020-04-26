package main

import (
	"os"

	"github.com/kazuya0202-dev/pyx/cmd"
)

func main() {
	ret := cmd.Execute()

	if ret != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
