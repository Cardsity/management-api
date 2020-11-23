package utils

import "testing"

func TestGetBlankCount(t *testing.T) {
	tables := []struct {
		text       string
		blankCount int
	}{
		{"Hello", 1},
		{"____", 1},
		{"____________", 1},
		{"___", 1},
		{"x___z", 1},
		{"x____z", 1},
		{" _________ ", 1},
		{" ___ ", 1},
		{"Hello ____, you are ____.", 2},
		{"Hello ___, you are ____.", 1},
		{"Hello ____, you are ___.", 1},
		{"Hello ______________, you are ___.", 1},
		{"Hello ______________, you are _____.", 2},
	}

	for _, table := range tables {
		result := GetBlankCount(table.text)
		if result != table.blankCount {
			t.Errorf("Checking '%v' for blanks does not return the right result, got: %v, want: %v.",
				table.text, result, table.blankCount)
		}
	}
}
