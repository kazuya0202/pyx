package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
)

// envCommand is struct.
type envCommand struct {
	Cmd    string // execute command
	Option string // command option
}

// getEnvCommand determines command depend in os environment.
func getEnvCommand() envCommand {
	// windows
	if runtime.GOOS == "windows" {
		return envCommand{"cmd", "/c"}
	}
	// other than windows
	return envCommand{"sh", "-c"}
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

func showUsage() {
	usage := `Description:
  This application is a tool to exec python scripts.
  If you put scripts in one directory and specify path of it,
  you can exec python scripts in any directory like 'npx' command.

Usage:
  pyx [script name] [flags] [script option]

Flags:
  -f, --find              find script using fuzzy-finder
  -l, --list              path of target directory
  -p, --path string       set path that target directory
  -s, --search string     search scripts that contains string (like grep)
      --set-path string   display version
  -h, --help              display help of command
  -v, --version           display version of command
`
	fmt.Print(usage)
}
