package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	filePath    string
	rowsPerFile int
	suffix      string
	trimS3Paths bool
)

var rootCmd = &cobra.Command{
	Use:   "splitcsv",
	Short: "Split a CSV file into multiple files",
	Long:  `A CLI application that splits a CSV file into multiple files, each containing a header row and a specified number of data rows.`,
	Run:   run,
}

func init() {
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the CSV file (required)")
	rootCmd.Flags().IntVarP(&rowsPerFile, "rows", "r", 10, "Number of rows per output file (required)")
	rootCmd.Flags().StringVarP(&suffix, "suffix", "s", "", "Suffix prepended to the output file name")
	rootCmd.Flags().BoolVarP(&trimS3Paths, "trims3", "t", false, "Trim the s3 paths like s3://file to /file")
	rootCmd.MarkFlagRequired("file")
	rootCmd.MarkFlagRequired("number")
}

func run(cmd *cobra.Command, args []string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		fmt.Printf("Error reading header: %v\n", err)
		os.Exit(1)
	}

	fileCounter := 1
	rowCounter := 0

	var records [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading record: %v\n", err)
			os.Exit(1)
		}

		records = append(records, record)
		rowCounter++

		// trim S3 paths
		if trimS3Paths {
			trimS3RecordPaths(record)
		}

		if rowCounter == rowsPerFile {
			writeFile(fileCounter, header, records, suffix)
			records = [][]string{}
			rowCounter = 0
			fileCounter++
		}
	}

	if len(records) > 0 {
		writeFile(fileCounter, header, records, suffix)
		fileCounter++
	}

	fmt.Printf("CSV file split into %d files successfully.\n", fileCounter-1)
}

func writeFile(fileCounter int, header []string, records [][]string, suffix string) {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	filenameWithoutExt := filename[:len(filename)-len(filepath.Ext(filename))]

	outputFileName := ""
	if suffix != "" {
		outputFileName = fmt.Sprintf("%s_%s_part_%d.csv", filenameWithoutExt, suffix, fileCounter)
	} else {
		outputFileName = fmt.Sprintf("%s_part_%d.csv", filenameWithoutExt, fileCounter)
	}

	fullPath := filepath.Join(dir, outputFileName)

	outputFile, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	writer.Write(header)
	writer.WriteAll(records)
}

func trimS3RecordPaths(record []string) {
	for i := range record {
		record[i] = strings.ReplaceAll(record[i], "s3://", "/")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
