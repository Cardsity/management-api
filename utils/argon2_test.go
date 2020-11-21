package utils

import (
	"testing"
)

func TestArgon2idHashStringWithSalt(t *testing.T) {
	input := "Hello world!"
	salt := []byte{240, 244, 57, 235, 142, 133, 0, 160, 193, 40, 70, 30, 235, 51, 79, 73}
	config := GetDefaultArgon2IDConfig()

	hash, err := argon2idHashStringWithSalt(input, salt, config)
	if err != nil {
		t.Errorf("Hashing 'Hello world!' produces an error: %v.", err)
	}
	if hash != "$argon2id$v=19$m=65536,t=1,p=4$8PQ5646FAKDBKEYe6zNPSQ$VavIZwx2xsUcvFDOtTpIfgQdwmyf7aUSywn8W8Sueho" {
		t.Error("Hashing 'Hello world!' does not return the right result.")
	}
}

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

func TestArgon2IDHashAndCompare(t *testing.T) {
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
