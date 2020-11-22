package utils

import "testing"

func TestMeetsPasswordRequirements(t *testing.T) {
	tables := []struct {
		password          string
		meetsRequirements bool
	}{
		{"Hello", false},
		{"H3ll0", false},
		{"henlo", false},
		{"helloworld", true},
		{"\000\000\000\000\000\000\000\000", true},
		{"Cardsity", true},
		{"Test123!", true},
		{"test1234", true},
	}

	for _, table := range tables {
		result := MeetsPasswordRequirements(table.password)
		if result != table.meetsRequirements {
			t.Errorf("Checking if '%v' meets the password requirements produces the wrong result, got: %v, want: %v.",
				table.password, result, table.meetsRequirements)
		}
	}
}
