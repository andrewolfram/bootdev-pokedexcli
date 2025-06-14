package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  only  LOWERCASE  ",
			expected: []string{"only", "lowercase"},
		},
		{
			input:    "gotta split them all",
			expected: []string{"gotta", "split", "them", "all"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("Input Length did not match expected length %d, but got %d", len(c.expected), len(actual))
			t.Fail()
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Words did not match expected  %s, but got %s", expectedWord, word)
				t.Fail()
			}
		}
	}

}
