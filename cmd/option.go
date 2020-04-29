package cmd

import (
	"fmt"
	"strings"
)

// Option is opiton template.
type Option struct {
	Long  string
	Short string
}

// Contains is ...
func (o *Option) Contains(args []string) bool {
	return o.IndexOf(args) > -1
}

// IndexOf is ...
func (o *Option) IndexOf(args []string) (idx int) {
	for i, x := range args {
		if o.Long == x || o.Short == x {
			return i
		}
	}
	return -1
}

// GetValue ...
func (o *Option) GetValueWithOption(args []string) (string, error) {
	idx := o.IndexOf(args)
	if idx > -1 && len(args)-1 > idx {
		return args[idx+1], nil
	}
	return "", errFatal
}

// CheckError ...
func (o *Option) CheckError(err error) error {
	if err != nil {
		x := fmt.Sprintf("%s %s", o.Short, o.Long)
		x = strings.Trim(x, " ")
		x = strings.ReplaceAll(x, " ", ", ")

		s := fmt.Sprintf("Option `%s` is required argument.", x)
		fmt.Println(MakeErrorMessage(s))

		return errFatal
	}
	return nil
}

// showUsage ...
func showUsage() {
	usage := `Description:
  This application is a tool to exec python scripts.
  If you put scripts in one directory and specify path of it,
  you can exec python scripts in any directory like 'npx' command.

Usage:
  pyx [flags] [script_name] [script_option]

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
