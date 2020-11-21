package utils

import (
	"testing"
)

func TestArgon2IDHashCompare(t *testing.T) {
	tables := []struct {
		hash  string
		s     string
		equal bool
	}{
		{"$argon2id$v=19$m=65536,t=1,p=4$UWi4RSQYGK1n/UfHRx1NQQ$femAm63seJzibGja1UXkzyEx1jmZZQlZMuZ0ng0aAyo", "Cardsity", true},
		{"$argon2id$v=19$m=65536,t=1,p=4$UWi4RSQYGK1n/UfHRx1NQQ$femAm63seJzibGja1UXkzyEx1jmZZQlZMuZ0ng0aAyo", "CarsCity", false},
		{"$argon2id$v=19$m=65536,t=1,p=4$V+4kAiucJM5MDzv8OC/rQg$qHuPtmWCyq8MqX/l19yaTeNL0VVzzNUtvvjCVpNM5B4", "🙂", true},
		{"$argon2id$v=19$m=65536,t=1,p=4$V+4kAiucJM5MDzv8OC/rQg$qHuPtmWCyq8MqX/l19yaTeNL0VVzzNUtvvjCVpNM5B4", "🙁", false},
	}

	for _, table := range tables {
		result, err := Argon2IDHashCompare(table.s, table.hash)
		if err != nil {
			t.Errorf("Comparing the hash '%v' with '%v' produced an error: %v.", table.hash, table.s, err)
		}
		if result != table.equal {
			t.Errorf("Comparing the hash '%v' with '%v' produced the wrong result, got: %v, want: %v.", table.hash, table.s, result, table.equal)
		}
	}
}

func TestArgon2IDHash(t *testing.T) {
	config := GetDefaultArgon2IDConfig()
	hash, err := Argon2IDHashString("Cardsity", config)
	if err != nil {
		t.Errorf("Hashing 'Cardsity' produces an error: %v.", err)
	}

	result, err := Argon2IDHashCompare("Cardsity", hash)
	if err != nil {
		t.Errorf("Comparing the hash of 'Cardsity' with 'Cardsity' produced an error: %v.", err)
	}
	if !result {
		t.Errorf("Comparing the recently generated hash of 'Cardsity' with 'Cardsity' does not produce the right result, got: %v, want: true.", result)
	}
}
