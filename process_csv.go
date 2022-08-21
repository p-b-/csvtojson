package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var jsonTemplate []string
var jsonEndTemplate []string
var lineIndex int

func ProcessUsingTemplate(headerLineCount uint, jsonTemplateFile, inputFile, outputFile *os.File) error {
	csvScanner := bufio.NewScanner(inputFile)
	jsonScanner := bufio.NewScanner(jsonTemplateFile)
	var headers []string
	var colValsNoQuotes []bool
	if headerLineCount > 0 {
		var err error
		headers, colValsNoQuotes, err = scanHeaders(headerLineCount, csvScanner)
		if err != nil {
			return err
		}
	}
	readTemplate(jsonScanner, outputFile)
	lineIndex = 1 + int(headerLineCount)
	firstRow := true
	for csvScanner.Scan() {
		scannedLine := csvScanner.Text()
		csvRowValues, err := parseLine(scannedLine)
		if err != nil {
			return err
		}
		headerToValue := connectValuesToHeaders(csvRowValues, headers)
		outputTemplate(firstRow, headerToValue, headers, colValsNoQuotes, jsonTemplate, outputFile)

		lineIndex++
		firstRow = false
	}
	outputEndTemplate(jsonEndTemplate, outputFile)
	return nil
}

func ProcessWithoutTemplate(headerLineCount uint, mainArray string, inputFile, outputFile *os.File) error {
	scanner := bufio.NewScanner(inputFile)
	var headers []string
	var colValsNoQuotes []bool
	if headerLineCount > 0 {
		var err error
		headers, colValsNoQuotes, err = scanHeaders(headerLineCount, scanner)
		if err != nil {
			return err
		}
	}

	firstJSONLine := fmt.Sprintf("{\n\t\"%s\": [", mainArray)
	outputFile.WriteString(firstJSONLine)

	tabCount := 3
	firstJSONRow := true
	lineIndex = 1 + int(headerLineCount)
	for scanner.Scan() {
		scannedLine := scanner.Text()
		csvRowValues, err := parseLine(scannedLine)
		if err != nil {
			return err
		}
		writeRowAsJson(firstJSONRow, tabCount, csvRowValues, headers, colValsNoQuotes, outputFile)
		firstJSONRow = false
		lineIndex++
	}
	outputFile.WriteString("\n")
	writeTabs(tabCount-2, outputFile)
	outputFile.WriteString("]\n}")
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return nil
}

func scanHeaders(headerLineCount uint, inputScanner *bufio.Scanner) ([]string, []bool, error) {
	csvLine := 0
	var headers []string
	var colValsNoQuotes []bool

	for inputScanner.Scan() {
		scannedLine := inputScanner.Text()
		scannedLine = removeByteOrderMark(scannedLine)
		if csvLine == int(headerLineCount-1) {
			var err error
			headers, colValsNoQuotes, err = parseHeaderLine(scannedLine)
			if err != nil {
				return nil, nil, err
			}
		}
		csvLine++
		if csvLine == int(headerLineCount) {
			break
		}
	}
	return headers, colValsNoQuotes, nil
}

func removeByteOrderMark(line string) string {
	return strings.TrimLeft(line, "\ufeff")
}

func writeRowAsJson(isFirstJsonRow bool, tabCount int, csvValues []string, headers []string, colValsNoQuotes []bool, outputFile *os.File) {
	if !isFirstJsonRow {
		outputFile.WriteString(",\n")
		writeTabs(tabCount-1, outputFile)
	}
	outputFile.WriteString("{\n")

	firstRow := true
	for index, value := range csvValues {
		if index < len(headers) {
			writeColumnAsJson(firstRow, tabCount, headers[index], value, colValsNoQuotes[index], outputFile)
		}
		firstRow = false
	}

	outputFile.WriteString("\n")
	writeTabs(tabCount-1, outputFile)
	outputFile.WriteString("}")
}

func writeColumnAsJson(isFirstCSVColumn bool, tabCount int, propertyName string, propertyValue string, withoutQuotes bool, outputFile *os.File) {
	if !isFirstCSVColumn {
		outputFile.WriteString(",\n")
	}
	writeTabs(tabCount, outputFile)
	var outputString string
	if withoutQuotes {
		outputString = fmt.Sprintf("\"%s\": %s", propertyName, propertyValue)
	} else {
		outputString = fmt.Sprintf("\"%s\": \"%s\"", propertyName, propertyValue)
	}
	outputFile.WriteString(outputString)
}

func writeTabs(tabCount int, outputFile *os.File) {
	outputFile.WriteString(strings.Repeat("\t", tabCount))
}

var lineAsChars []rune
var lineLength int
var runeIndex int

func parseHeaderLine(line string) ([]string, []bool, error) {
	headers := make([]string, 0)
	var columnValuesNoQuotes []bool

	columnValuesNoQuotes = make([]bool, 0)
	line = strings.TrimSpace(line)
	lineAsChars = []rune(line)
	lineLength = len(lineAsChars)
	runeIndex = 0

	for runeIndex < lineLength {
		currentHeader, err := getNextLineValue()
		if err != nil {
			return nil, nil, err
		}
		currentHeader = strings.TrimSpace(currentHeader)
		jsonFieldDoesNotNeedQuotes := false
		if strings.Contains(currentHeader, ":") {
			lineSplit := strings.Split(currentHeader, ":")
			currentHeader = lineSplit[0]

			if len(lineSplit) > 1 {
				if strings.ToLower(lineSplit[1]) == "bool" ||
					strings.ToLower(lineSplit[1]) == "number" {
					jsonFieldDoesNotNeedQuotes = true
				}
			}
		}
		headers = append(headers, currentHeader)
		columnValuesNoQuotes = append(columnValuesNoQuotes, jsonFieldDoesNotNeedQuotes)
	}
	return headers, columnValuesNoQuotes, nil
}

func parseLine(line string) ([]string, error) {
	toReturn := make([]string, 0)
	line = strings.TrimSpace(line)
	lineAsChars = []rune(line)
	lineLength = len(lineAsChars)
	runeIndex = 0

	for runeIndex < lineLength {
		nextValue, err := getNextLineValue()
		if err != nil {
			return toReturn, err
		}
		nextValue = strings.TrimSpace(nextValue)
		toReturn = append(toReturn, nextValue)
	}
	return toReturn, nil
}

func getNextLineValue() (string, error) {
	toReturn := ""
	startsWithQuotes := false
	nextChar := ""

	nextChar, startsWithQuotes, err := processValueStart()
	if err != nil {
		return "", err
	}

	for nextChar != "" {
		if nextChar == "\\" {
			nextChar, toReturn = processEscapeCharacter(toReturn)
		} else if startsWithQuotes && nextChar == "\"" {
			nextChar = getNextRune()
			if nextChar == "" || nextChar == "," {
				return toReturn, nil
			} else {
				err := fmt.Errorf("unexpected character '%s' on line %d", nextChar, lineIndex)
				return "", err
			}
		} else if !startsWithQuotes && nextChar == "," {
			return toReturn, nil
		} else {
			toReturn += nextChar
			nextChar = getNextRune()
		}
	}
	if startsWithQuotes {
		err := fmt.Errorf("value should end with a quote on line %d", lineIndex)
		return toReturn, err
	} else {
		return toReturn, nil
	}
}

func processValueStart() (nextChar string, startsWithQuotes bool, err error) {
	for runeIndex < lineLength {
		nextChar = getNextRune()
		if nextChar == "" {
			err = fmt.Errorf("cannot process line %d", lineIndex)
			return
		}
		if nextChar == "\"" {
			if startsWithQuotes {
				return
			}
			startsWithQuotes = true
		} else if nextChar == " " || nextChar == "\t" {
			continue
		} else {
			break
		}
	}
	return
}

func getNextRune() string {
	if runeIndex == lineLength {
		return ""
	}
	char := string(lineAsChars[runeIndex])
	runeIndex++
	return char
}

func processEscapeCharacter(toReturn string) (string, string) {
	nextChar := getNextRune()
	if nextChar == "" {
		return "", toReturn
	}
	toReturn += "\""
	toReturn += nextChar

	nextChar = getNextRune()
	return nextChar, toReturn
}

func connectValuesToHeaders(values []string, headers []string) *MultiMap[string, string] {
	returnValue := NewMultiMap[string, string]()

	for index, value := range values {
		if index >= len(headers) {
			break
		}
		columnHeader := headers[index]
		returnValue.AddKeyValue(strings.ToLower(columnHeader), value)
	}

	return returnValue
}

func readTemplate(templateScanner *bufio.Scanner, outputFile *os.File) {
	scanStatus := 0
	jsonTemplate = make([]string, 0)
	jsonEndTemplate = make([]string, 0)
	for templateScanner.Scan() {
		scannedLine := templateScanner.Text()
		switch scanStatus {
		case 0:
			if strings.ToLower(scannedLine) == "%row start%" {
				scanStatus = 1
			} else {
				outputFile.WriteString(scannedLine)
				outputFile.WriteString("\n")
			}
		case 1:
			if strings.ToLower(scannedLine) == "%row end%" {
				scanStatus = 2
			} else {
				jsonTemplate = append(jsonTemplate, scannedLine)
			}
		case 2:
			jsonEndTemplate = append(jsonEndTemplate, scannedLine)
		}
	}
}

func outputTemplate(firstRow bool, headerToValue *MultiMap[string, string], headers []string, colValsNoQuotes []bool, template []string, outputFile *os.File) {
	if !firstRow {
		outputFile.WriteString(",\n")
	}
	lineCount := len(template)
	for index, originalLine := range template {
		line := originalLine
		loop := true
		outputLine := true
		for loop {
			loop = false
			if strings.Contains(line, "%") {
				alteredLine, replaced := replaceFromCSVValues(line, headerToValue, headers, colValsNoQuotes)
				if replaced {
					line = alteredLine
					loop = true
				} else {
					outputLine = false
				}
			}
		}
		if outputLine {
			outputFile.WriteString(line)
			if index < lineCount-1 {
				outputFile.WriteString("\n")
			}
		}
	}
}

func replaceFromCSVValues(line string, headerToValue *MultiMap[string, string], headers []string, colValsNoQuotes []bool) (string, bool) {
	firstMarker := strings.Index(line, "%")
	subString := line[firstMarker+1:]
	secondMarker := strings.Index(subString, "%")

	replaced := false
	if secondMarker != -1 {
		subString = strings.ToLower(subString[:secondMarker])
		replaceWith, _ := headerToValue.GetFirstValueIfKeyExists(subString)

		headerToValue.RemoveKeyValue(subString, replaceWith)
		if columnValueNeedsSurroundingQuotes(subString, headers, colValsNoQuotes) {
			line = line[:firstMarker] + "\"" + replaceWith + "\"" + line[firstMarker+1+secondMarker+1:]
		} else {
			line = line[:firstMarker] + replaceWith + line[firstMarker+1+secondMarker+1:]
		}
		replaced = true
	}
	return line, replaced
}

func outputEndTemplate(template []string, outputFile *os.File) {
	if len(template) > 0 {
		outputFile.WriteString(("\n"))
	}
	for _, line := range template {
		outputFile.WriteString(line)
		outputFile.WriteString("\n")
	}
}

func columnValueNeedsSurroundingQuotes(column string, headers []string, colValsNoQuotes []bool) bool {
	column = strings.ToLower(column)
	for index, quote := range headers {
		if strings.ToLower(quote) == column {
			return !colValsNoQuotes[index]
		}
	}
	return false
}
