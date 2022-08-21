package main

import (
	"strings"
	"testing"
)

func Test_ParseHeaderLine(t *testing.T) {
	testCases := []struct {
		expectedSuccess    bool
		columnHeaders      []string
		columnBoolOrNumber []bool
	}{
		{true, []string{"Header1", "Header2", "Header3"}, []bool{false, false, false}},
		{true, []string{"\"Header1\"", "\"Header2\"", "\"Header3\""}, []bool{false, false, false}},
		{true, []string{"Header1"}, []bool{false, false, false}},
		{true, []string{"\"Header1\""}, []bool{false}},
		{false, []string{"\"Header1\"", "\"Header2"}, []bool{false, false}},
		{true, []string{"Header1:bool", "Header2:number", "Header3"}, []bool{true, true, false}},
		{true, []string{"Header1:bool", "Header2:Number", "Header3"}, []bool{true, true, false}},
	}

	for testCaseIndex, testCase := range testCases {
		colHeaderCount := len(testCase.columnHeaders)
		headerLine := ""
		for _, colHeader := range testCase.columnHeaders {
			if len(headerLine) > 0 {
				headerLine += ","
			}
			headerLine += colHeader
		}
		parsedHeaders, parsedColValsNoQuotes, err := parseHeaderLine(headerLine)
		if testCase.expectedSuccess {
			if err != nil {
				t.Errorf("Test case %d: Parsing header return error %s", testCaseIndex, err)
			} else {
				if len(parsedHeaders) != colHeaderCount {
					t.Errorf("Test case index %d: Insufficent headers parsed, was expecting %d and got %d", testCaseIndex, colHeaderCount, len(parsedHeaders))
				} else if len(parsedColValsNoQuotes) != colHeaderCount {
					t.Errorf("Test case index %d: Insufficent header types parsed, was expecting %d and got %d", testCaseIndex, colHeaderCount, len(parsedColValsNoQuotes))
				} else {
					for index, parsedHeader := range parsedHeaders {
						expecting := strings.Trim(testCase.columnHeaders[index], "\"")
						if testCase.columnBoolOrNumber[index] {
							indexOfColon := strings.Index(expecting, ":")
							expecting = expecting[:indexOfColon]
						}
						if parsedHeader != expecting {
							t.Errorf("Test case index %d: Parsed header %d is %s, expected %s", testCaseIndex, index, parsedHeader, expecting)
						}
					}
					for index, parsedColVal := range parsedColValsNoQuotes {
						expecting := testCase.columnBoolOrNumber[index]
						if parsedColVal != expecting {
							t.Errorf("Test case index %d: Parsed header %d is marked as bool or number %t, expected %t", testCaseIndex, index, parsedColVal, expecting)
						}
					}
				}
			}
		} else {
			if err == nil {
				t.Errorf("Test case was expected to fail, but passed")
			}
		}
	}
}
