package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	kz "github.com/kazuya0202/kazuya0202"
)

// DictString ...
type DictString struct {
	Key string
	Val string
}

// DictStringA ...
type DictStringA struct {
	Key string
	Val []string
}

// JSONDefaults ...
type JSONDefaults struct {
	Cmd  DictString
	Path DictString
	Ext  DictStringA
}

// JSONKeys ...
type JSONKeys struct {
	Cmd  string
	Path string
	Exts []string
}

var (
	// command
	cmd CommandUtility

	// options
	opts struct {
		Find    Option // find script using fuzzy-finder
		Help    Option // display help of command
		List    Option // display list of executable scripts
		Path    Option // path of target directory
		SetPath Option // set path that target directory
		Search  Option // search scripts that contains string (like grep)
		Version Option // display version
	}

	// json
	ju *kz.JSONUtility

	// json content
	jc JSONKeys

	// default content (it use when key is unset)
	defs = JSONDefaults{
		Cmd:  DictString{"cmd", "python"},
		Path: DictString{"path", kz.GetUserHomeDir()},
		Ext:  DictStringA{"extension", []string{".py"}},
	}

	// error fatal
	errFatal  = errors.New("fatal")  // os.Exit(1)
	errNormal = errors.New("normal") // os.Exit(0)
)

func parseArgs(args []string) error {
	args = args[1:] // remove execute file

	if opts.Version.Contains(args) {
		println("pyx version: v0.1.0")
		return errNormal
	}

	// set-path option
	if opts.SetPath.Contains(args) {
		p, err := opts.SetPath.GetValue(args)
		if err != nil || !kz.Exists(p) {
			fmt.Printf("The '%s' is not exist.\n", p)
			println("or Unspecified argment.")

			for {
				p = kz.GetInput("path")
				if kz.Exists(p) {
					break
				}
				fmt.Printf("The '%s' is not exist.\n", p)
			}
		}
		ju.Set(defs.Path.Key, p)
		ju.Dump(ju.Path)

		println("\nSet path of scripts in '$HOME/.config/kazuya0202/pyx.json'")
		println(ju.Data)
		return errNormal
	}

	if opts.Path.Contains(args) {
		p, err := opts.Path.GetValue(args)
		if opts.Path.CheckError(err) != nil {
			return errFatal
		}
		if !kz.Exists(p) {
			fmt.Printf("The '%s' is not exist.\n", p)
			return errFatal
		}
	}

	// help option
	idx := opts.Help.IndexOf(args)
	// help of pyx command
	if idx == 0 {
		showUsage()
		return errNormal
	}

	if idx > -1 {
		// help of python script
		cmd.Option = args[idx]
	} else if idx == -1 {
		// not help
		if len(args) >= 1 {
			cmd.Option = strings.Join(args[1:], " ")
		}
	}

	files := kz.GetFiles(jc.Path)

	// show list
	if opts.List.Contains(args) {
		displayList(files, jc.Exts)
		return errNormal
	}

	// search
	if opts.Search.Contains(args) {
		target, err := opts.Search.GetValue(args)
		if opts.Search.CheckError(err) != nil {
			return errFatal
		}
		search(files, target)
		return errNormal
	}

	if opts.Find.Contains(args) {
		selected := selectFzf(files)
		println(selected)
		return errNormal
	}

	// input
	if len(args) == 0 {
		args = append(args, kz.GetInput("script"))
		cmd.Option = kz.GetInput("option")
	}

	isExist := false
	scriptName := args[0]
	scriptPath := path.Join(jc.Path, scriptName)

	for _, x := range jc.Exts {
		p := scriptPath + x
		if kz.Exists(p) {
			scriptPath = p
			isExist = true
			break
		}
	}

	if !isExist {
		fmt.Printf("The script name '%s' is not exist.\n", scriptName)
		return errFatal
	}

	cmd.Arg = fmt.Sprintf("%s %s %s", cmd.CmdName, scriptPath, cmd.Option)
	return nil
}

// initialize
func initialize() {
	// option | Option{Long, Short}
	opts.Help = Option{"--help", "-h"}
	opts.SetPath = Option{"--set-path", ""}
	opts.Search = Option{"--search", "-s"}
	opts.List = Option{"--list", "-l"}
	opts.Find = Option{"--find", "-f"}
	opts.Path = Option{"--path", "-p"}
	opts.Version = Option{"--version", "-v"}

	// config
	jsonPath := path.Join(kz.GetUserHomeDir(), ".config/kazuya0202/pyx.json")
	ju = kz.NewJSONUtility(jsonPath)

	jc.Cmd = ju.Get(defs.Cmd.Key).String()
	jc.Path = ju.Get(defs.Path.Key).String()
	for _, x := range ju.Get(defs.Ext.Key).Array() {
		jc.Exts = append(jc.Exts, x.String())
	}

	// null -> write json, set default value
	any := false
	if isNull(jc.Cmd) {
		any = true
		jc.Cmd = defs.Cmd.Val
		ju.Set(defs.Cmd.Key, defs.Cmd.Val)
	}
	if isNull(jc.Path) {
		any = true
		jc.Path = defs.Path.Val
		ju.Set(defs.Path.Key, defs.Path.Val)
	}
	if len(jc.Exts) == 0 {
		any = true
		jc.Exts = defs.Ext.Val
		ju.Set(defs.Ext.Key, defs.Ext.Val)
	}
	if any {
		ju.Dump(jsonPath)
	}

	// cmd
	cmd.CmdName = jc.Cmd
	cmd.EnvCmd = getEnvCommand()
}

// Execute ...
func Execute() error {
	initialize()
	err := parseArgs(os.Args)

	// after processing option
	if err == errNormal {
		os.Exit(0)
	} else if err == errFatal {
		os.Exit(1)
	}

	if err == nil {
		return cmd.execute()
	}
	return nil
}
