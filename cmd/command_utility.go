package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	kz "github.com/kazuya0202/kazuya0202"
)

// CommandUtility is struct.
type CommandUtility struct {
	CmdName string
	Option  string
	Arg     string
	Command *exec.Cmd
	EnvCmd  kz.EnvCommand
}

func (c *CommandUtility) execute() error {
	c.setCommand()
	fmt.Printf("[%s]: %s\n", color.BlueString("EXECUTE"), c.shapeCommandString())
	println("---------")

	err := kz.ExecCmdInRealTime(c.Command)
	return err
}

func (c *CommandUtility) setCommand() {
	c.Command = exec.Command(c.EnvCmd.Cmd, c.EnvCmd.Option, c.Arg)
}

func (c *CommandUtility) shapeCommandString() string {
	// <appName> <opt> <cmdName> ... -> <cmdName> ...
	str := c.Command.String()
	str = str[strings.Index(str, c.CmdName):]
	str = strings.ReplaceAll(str, "  ", " ") // "  " -> " "
	return str
}
