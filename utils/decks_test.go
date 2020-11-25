package utils

import "testing"

func TestDeckCardAmountValid(t *testing.T) {
	tables := []struct {
		amountWhiteCards int
		amountBlackCards int
		valid            bool
	}{
		{-100, 37, false},
		{-1000, 100000, false},
		{0, 0, false},
		{2, 4, false},
		{4, 5, false},
		{5, 4, false},
		{5, 5, true},
		{6, 5, true},
		{5, 6, true},
		{6, 6, true},
		{10, 10, true},
		{50, 10, true},
		{10, 50, true},
		{250, 250, true},
		{251, 250, false},
		{250, 251, false},
		{251, 251, false},
		{500, 0, false},
		{0, 500, false},
		{500, 5, false},
		{5, 500, false},
		{495, 5, true},
		{5, 495, true},
		{495, 6, false},
		{6, 495, false},
	}

	for _, table := range tables {
		result := DeckCardAmountValid(table.amountWhiteCards, table.amountBlackCards)
		if result != table.valid {
			t.Errorf("Checking if %v white and %v black cards are a valid deck does not return the right result, got: %v, want: %v.",
				table.amountWhiteCards, table.amountBlackCards, result, table.valid)
		}
	}
}
