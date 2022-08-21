package main

import "testing"

func Test_GetKeyValues(t *testing.T) {
	testCases := []struct {
		key string
		val string
	}{
		{"Hello", "World"},
		{"Goodbye", "Moon"},
		{"Everybody's free", "To Wear Sunscreen"},
	}
	mm := NewMultiMap[string, string]()
	for _, tc := range testCases {
		mm.AddKeyValue(tc.key, tc.val)
	}

	for _, tc := range testCases {
		value, exists := mm.GetFirstValueIfKeyExists(tc.key)
		if !exists || value != tc.val {
			t.Errorf("GetFirstValueIfKeyExists[string,string](\"%s\") returned (%s,%t), wanted (\"%s\",true)", tc.key, value, exists, tc.val)
		}
	}
}

func Test_GetFirstKeyValue(t *testing.T) {
	testCases := []struct {
		key string
		val string
	}{
		{"Hello", "World"},
		{"Hello", "Moon"},
		{"Hello", "Let's dance"},
	}
	mm := NewMultiMap[string, string]()
	for _, tc := range testCases {
		mm.AddKeyValue(tc.key, tc.val)
	}

	for _, tc := range testCases {
		value, exists := mm.GetFirstValueIfKeyExists(tc.key)
		if !exists || value != "World" {
			t.Errorf("GetFirstValueIfKeyExists[string,string](\"%s\") returned (%s,%t), wanted (\"World\",true)", tc.key, value, exists)
		}
	}
}

func Test_ContainsKey(t *testing.T) {
	testCases := []struct {
		key string
		val string
	}{
		{"Hello", "World"},
		{"Goodbye", "Moon"},
		{"Everybody's free", "To Wear Sunscreen"},
	}
	nonExistingTestCases := []struct {
		key string
		val string
	}{
		{"abcdefg", "World"},
		{"11111", "World"},
		{"Watch out for the wolf", "World"},
	}
	mm := NewMultiMap[string, string]()
	for _, tc := range testCases {
		mm.AddKeyValue(tc.key, tc.val)
	}

	for _, tc := range testCases {
		exists := mm.ContainsKey(tc.key)
		if !exists {
			t.Errorf("ContainsKey[string,string](\"%s\") returned (%t), wanted (true)", tc.key, exists)
		}
	}
	for _, tc := range nonExistingTestCases {
		exists := mm.ContainsKey(tc.key)
		if exists {
			t.Errorf("ContainsKey[string,string](\"%s\") returned (%t), wanted (false)", tc.key, exists)
		}
	}
}

func Test_RemoveKeyValue(t *testing.T) {
	testCases := []struct {
		key string
		val string
	}{
		{"Hello", "World"},
		{"Goodbye", "Moon"},
		{"Everybody's free", "To Wear Sunscreen"},
		{"Hello", "Moon"},
		{"Hello", "Let's dance"},
	}
	mm := NewMultiMap[string, string]()
	for _, tc := range testCases {
		mm.AddKeyValue(tc.key, tc.val)
	}

	for _, tc := range testCases {
		value, exists := mm.GetFirstValueIfKeyExists(tc.key)
		if !exists || value != tc.val {
			// Rest of tests may fail if state of map isn't as expected
			t.Fatalf("GetFirstValueIfKeyExists[string,string](\"%s\") returned (%s,%t), wanted (\"%s\",true)", tc.key, value, exists, tc.val)
		}

		mm.RemoveKeyValue(tc.key, tc.val)

		value, exists = mm.GetFirstValueIfKeyExists(tc.key)
		if exists && value == tc.val {
			t.Errorf("GetFirstValueIfKeyExists[string,string](\"%s\") returned (%s,%t), wanted (\"%s\",false)", tc.key, value, exists, tc.val)
		}
	}
}
