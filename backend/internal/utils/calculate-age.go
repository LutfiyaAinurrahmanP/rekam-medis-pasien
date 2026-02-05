package utils

import "time"

func CalculateAge(dateOfBirthStr string) (int, error) {
	birthDate, err := time.Parse("2006-01-02", dateOfBirthStr) // Format YYYY-MM-DD
	if err != nil {
		return 0, err
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// Jika bulan/hari saat ini belum mencapai ulang tahun, kurangi 1 tahun
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age, nil
}
