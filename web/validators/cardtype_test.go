package validators

import "testing"

func TestIsValidCardType(t *testing.T) {
	tables := []struct {
		value   CardType
		isValid bool
	}{
		{0, false},
		{1, true},
		{2, true},
		{3, false},
		{100, false},
		{7000, false},
	}

	for _, table := range tables {
		result := IsValidCardType(table.value)
		if result != table.isValid {
			t.Errorf("Checking if '%v' is a valid card type returns the wrong result, got: %v, want: %v.",
				table.value, result, table.isValid)
		}
	}
}
