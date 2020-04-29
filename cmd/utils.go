package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
)

// isOption
func isOption(s string) bool {
	if len(s) > 0 {
		return s[:1] == "-"
	}
	return false
}

func selectFzf(array []string) string {
	idx, _ := fuzzyfinder.Find(
		array,
		func(i int) string { return array[i] },
	)
	return array[idx]
}

func search(files []string, target string) {
	println("Searching executable scripts ...")

	isExist := false
	for _, x := range files {
		if strings.Index(x, target) > -1 {
			s := strings.ReplaceAll(x, target, color.RedString(target))
			fmt.Printf(" %s %s\n", color.BlueString("*"), s)
			isExist = true
		}
	}
	if !isExist {
		println("  None.")
	}
}

func displayList(files []string, exts []string) {
	println("Executable scripts ...")

	isExist := false
	for _, x := range files {
		for _, ext := range exts {
			if strings.Contains(x, ext) {
				fmt.Printf(" %s %s\n", color.BlueString("*"), x)
				isExist = true
				break
			}
		}
	}
	if !isExist {
		println("  None.")
	}
}

// MakeErrorMessage ...
func MakeErrorMessage(s string) string {
	return fmt.Sprintf("[%s]: %s", color.RedString("ERROR"), s)
}

// getNullElse ...
func getNullElse(x string, y string) (s string) {
	if x == "" {
		return y
	}
	return x
}

// getNullElseArray ...
func getNullElseArray(x []string, y []string) (s []string) {
	if len(x) == 0 {
		return y
	}
	return x
}

func isNull(s string) bool {
	if s == "" {
		return true
	}
	return false
}
