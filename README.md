# SplitCSV

A CLI application that splits a CSV file into multiple files, each containing a header row and a specified number of data rows.

```sh
Usage:
  splitcsv [flags]

Flags:
  -f, --file string     Path to the CSV file (required)
  -h, --help            help for splitcsv
  -r, --rows int        Number of rows per output file (required) (default 10)
  -s, --suffix string   Suffix prepended to the output file name
  -t, --trims3          Trim the s3 paths like s3://file to /file

```

## Download

You can download precompiled binary in [Releases](https://github.com/AugustDev/splitcsv/releases/latest).

## Build

You should have Go installed.

```sh
make build-all
```

which will build binaries in `./bin` directory.

## Example usage

```sh
splitcsv --file file.csv --rows 4 --suffix float --trims3
splitcsv -f file.csv -r 3
```
