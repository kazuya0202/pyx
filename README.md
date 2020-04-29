# pyx - execute python file in anywhere

## First Setup

### Installation

Download binary from [release]().

### Set path that scripts are putting

```sh
$ pyx --set-path <path>
```

## Basic usage

```sh
$ pyx [flags] [script_name[.py]] [script_option]

# example (using argparse.ArgumentParser)
$ pyx random_cp -n 100
$ pyx random_cp.py --help
```

## flags

### --set-path

Set path that scripts are putting.

```sh
$ pyx --set-path <path>
```

### -f / --find

Select a script using fuzzy-finder, and specify option of script.

```sh
$ pyx -f
$ pyx --find
```

### -s / --search

Display scripts that contains `<string>`  (like `grep`).

```sh
$ pyx -s <string>
$ pyx --search <string>

# example
$ pyx -s a
```

### -l / --list

Display all scripts.

```sh
$ pyx -l
$ pyx --list
```

### -p / --path

Specify target directory path only one time.

When directory path uses every time, execute `--set-path`.

```sh
$ pyx -p <path>
$ pyx --path <path>
```

## Other Options

### -h / --help

Display help message of command.

```sh
$ pyx -h
$ pyx --help
```

**Note:** If python script implement `argparse.ArgumantParser`,  `$ pyx [script_name] -h` and `$ pyx -h [script_name]` take difference actions.

```sh
$ pyx [script_name] -h
# => python {SCIRPT_NAME}.py -h
#    (display help of python script.)
```

```sh
$ pyx -h [script_name]
# => pyx -h
#    (display help of pyx command. ignore python script.)
```

### -v / --version

Display version of command.

```sh
$ pyx -v
$ pyx --version
```

