package utils

import (
	"testing"
	"time"
)

func TestCalculateAge_DateOnly(t *testing.T) {
	age, err := CalculateAge("1990-01-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsedBirthDate, _ := time.Parse(time.DateOnly, "1990-01-01")
	expected := time.Now().Year() - parsedBirthDate.Year()
	if time.Now().Month() < parsedBirthDate.Month() || (time.Now().Month() == parsedBirthDate.Month() && time.Now().Day() < parsedBirthDate.Day()) {
		expected--
	}

	if age != expected {
		t.Fatalf("expected age %d, got %d", expected, age)
	}
}

func TestCalculateAge_RFC3339(t *testing.T) {
	age, err := CalculateAge("1990-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsedBirthDate, _ := time.Parse(time.DateOnly, "1990-01-01")
	expected := time.Now().Year() - parsedBirthDate.Year()
	if time.Now().Month() < parsedBirthDate.Month() || (time.Now().Month() == parsedBirthDate.Month() && time.Now().Day() < parsedBirthDate.Day()) {
		expected--
	}

	if age != expected {
		t.Fatalf("expected age %d, got %d", expected, age)
	}
}

func TestCalculateAge_RejectsEmptyValue(t *testing.T) {
	_, err := CalculateAge("")
	if err == nil {
		t.Fatal("expected error for empty date of birth")
	}
}
