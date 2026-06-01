package utils

import (
	"errors"
	"strings"
	"time"
)

func CalculateAge(dateOfBirthStr string) (int, error) {
	birthDate, err := parseDateOfBirth(dateOfBirthStr)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// Jika bulan/hari saat ini belum mencapai ulang tahun, kurangi 1 tahun.
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age, nil
}

func parseDateOfBirth(dateOfBirthStr string) (time.Time, error) {
	value := strings.TrimSpace(dateOfBirthStr)
	if value == "" {
		return time.Time{}, errors.New("date of birth is empty")
	}

	layouts := []string{
		time.DateOnly,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z0700",
		"2006-01-02T15:04:05Z07:00",
	}

	for _, layout := range layouts {
		if parsed, parseErr := time.Parse(layout, value); parseErr == nil {
			return parsed, nil
		}
	}

	if len(value) >= len(time.DateOnly) {
		if parsed, parseErr := time.Parse(time.DateOnly, value[:len(time.DateOnly)]); parseErr == nil {
			return parsed, nil
		}
	}

	return time.Time{}, errors.New("unsupported date of birth format")
}
