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

// JSONDefaultKeys ...
type JSONDefaultKeys struct {
	Cmd  DictString
	Path DictString
	Exts DictStringA
}

// JSONKeys ...
type JSONKeys struct {
	Cmd  string
	Path string
	Exts []string
}

const version string = "0.1.0"

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

	jUtil *kz.JSONUtility // json
	jKey  JSONKeys        // json content

	// default content (it use when key is unset)
	defs = JSONDefaultKeys{
		Cmd:  DictString{"cmd", "python"},
		Path: DictString{"path", kz.GetUserHomeDir()},
		Exts: DictStringA{"extension", []string{".py"}},
	}

	cfgPath string // config path (under $HOME)

	// error fatal
	errFatal  = errors.New("fatal")  // os.Exit(1)
	errNormal = errors.New("normal") // os.Exit(0)
)

func parseArgs(args []string) error {
	args = args[1:] // remove execute file

	if opts.Version.Contains(args) {
		fmt.Printf("pyx version: v%s", version)
		return errNormal
	}

	// set-path option
	if opts.SetPath.Contains(args) {
		p, err := opts.SetPath.GetValueWithOption(args)
		if err != nil || !kz.Exists(p) {
			fmt.Printf("The '%s' is not exist.\nor Unspecified argment.", p)

			for {
				p = kz.GetInput("path")
				if kz.Exists(p) {
					break
				}
				fmt.Printf("The '%s' is not exist.\n", p)
			}
		}
		jUtil.Set(defs.Path.Key, p)
		jUtil.Dump(jUtil.Path)

		fmt.Printf("\nSet path of scripts in '$HOME/%s'", cfgPath)
		println(jUtil.Data)
		return errNormal
	}

	if opts.Path.Contains(args) {
		p, err := opts.Path.GetValueWithOption(args)
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

	files := kz.GetFiles(jKey.Path)

	// show list
	if opts.List.Contains(args) {
		displayList(files, jKey.Exts)
		return errNormal
	}

	// search
	if opts.Search.Contains(args) {
		target, err := opts.Search.GetValueWithOption(args)
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
	scriptPath := path.Join(jKey.Path, scriptName)

	for _, x := range jKey.Exts {
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
	// difine option | Option{Long, Short}
	opts.Help = Option{"--help", "-h"}
	opts.SetPath = Option{"--set-path", ""}
	opts.Search = Option{"--search", "-s"}
	opts.List = Option{"--list", "-l"}
	opts.Find = Option{"--find", "-f"}
	opts.Path = Option{"--path", "-p"}
	opts.Version = Option{"--version", "-v"}

	cfgPath = ".config/kazuya0202/pyx.json"

	// config
	jsonPath := path.Join(kz.GetUserHomeDir(), cfgPath)
	jUtil = kz.NewJSONUtility(jsonPath)

	jKey.Cmd = jUtil.Get(defs.Cmd.Key).String()
	jKey.Path = jUtil.Get(defs.Path.Key).String()
	for _, x := range jUtil.Get(defs.Exts.Key).Array() {
		jKey.Exts = append(jKey.Exts, x.String())
	}

	// null -> write json, set default value
	tmpData := jUtil.Data
	if isNull(jKey.Cmd) {
		jKey.Cmd = defs.Cmd.Val
		jUtil.Set(defs.Cmd.Key, defs.Cmd.Val)
	}
	if isNull(jKey.Path) {
		jKey.Path = defs.Path.Val
		jUtil.Set(defs.Path.Key, defs.Path.Val)
	}
	if len(jKey.Exts) == 0 {
		jKey.Exts = defs.Exts.Val
		jUtil.Set(defs.Exts.Key, defs.Exts.Val)
	}
	if jUtil.Data != tmpData {
		jUtil.Dump(jsonPath)
	}

	// cmd
	cmd.CmdName = jKey.Cmd
	cmd.EnvCmd.DetermineEnvCommand()
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

	var ret error = nil
	if err == nil {
		ret = cmd.execute()
	}
	return ret
}
