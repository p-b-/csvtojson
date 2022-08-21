package main

import (
	"flag"
	"fmt"
	"os"
)

var inputFilenameArg = flag.String("i", "", "Input file")
var outputFilenameArg = flag.String("o", "", "Output file")
var mainArrayArg = flag.String("a", "csv", "Main json array identifier that contains the rows of the CSV file")
var headerLineCountArg = flag.Uint("h", 1, "Number of header lines in CSV file")
var jsonTemplateFilenameArg = flag.String("t", "", "File that contains a json template, where CSV column values are surrounded by %'s instead of quotations")

func main() {
	flag.Parse()
	inputFile, outputFile, jsonTemplateFile, err := openFiles(*inputFilenameArg, *outputFilenameArg, *jsonTemplateFilenameArg)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer inputFile.Close()
	defer outputFile.Close()

	if jsonTemplateFile != nil {
		defer jsonTemplateFile.Close()
		err = ProcessUsingTemplate(*headerLineCountArg, jsonTemplateFile, inputFile, outputFile)
	} else {
		err = ProcessWithoutTemplate(*headerLineCountArg, *mainArrayArg, inputFile, outputFile)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func openFiles(inputFilename, outputFilename, jsonTemplateFilename string) (inputFile *os.File, outputFile *os.File, jsonTemplateFile *os.File, err error) {
	inputFile, err = os.Open(inputFilename)
	if err != nil {
		err = fmt.Errorf("could not open input file %s: %s", inputFilename, err)
		return
	}

	outputFile, err = os.Create(outputFilename)
	if err != nil {
		inputFile.Close()
		inputFile = nil
		err = fmt.Errorf("could not create output file %s: %s", outputFilename, err)
		return
	}

	if len(jsonTemplateFilename) > 0 {
		jsonTemplateFile, err = os.Open(jsonTemplateFilename)
		if err != nil {
			inputFile.Close()
			outputFile.Close()
			inputFile = nil
			outputFile = nil

			err = fmt.Errorf("could not open json template file %s: %s", jsonTemplateFilename, err)
			return
		}
	}
	return
}
