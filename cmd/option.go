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
func (o *Option) GetValue(args []string) (string, error) {
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

// isOption
func isOption(s string) bool {
	if len(s) > 0 {
		return s[:1] == "-"
	}
	return false
}
