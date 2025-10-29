package utils_test

import (
	"testing"
	"wilin.com/api/src/server/utils"
)

func TestValidEmails(t *testing.T) {
	expected := true
	validEmails := []string{
		"test@gmail.com",
		"person@mail.ca",
		"allo@olla.gov",
	}
	for _, email := range validEmails {
		got := utils.IsValidEmail(email)
		if got != expected {
			t.Errorf("IsValidEmail(%v), got %v: expected %v", email, got, expected)
		}
	}
}

func TestInvalidEmails(t *testing.T) {
	expected := false
	invalidEmails := []string{
		"failme",
		"sfksdfhui",
		"893jlkdh",
		"fhsdifhu@",
		"@hkjsfhk",
		"kawa.com",
	}
	for _, email := range invalidEmails {
		got := utils.IsValidEmail(email)
		if got != expected {
			t.Errorf("IsValidEmail(%v), got %v: expected %v", email, got, expected)
		}
	}
}

func TestSamePasswordHasSameHash(t *testing.T) {
	expected := true
	passwords := []string{
		"test",
		"123",
		"boo",
		"securepassword",
		"sigmarizz",
		"noOneWillGuessThisPasswordEver42069!",
	}
	for _, password := range passwords {
		hash, err := utils.GeneratePasswordHash(password)
		if err != nil {
			t.Errorf("GeneratePasswordHash(%s) failed to generate hash: %v", password, err)
			continue
		}
		got := utils.IsPasswordAndHashSame(password, hash)
		if got != expected {
			t.Errorf("IsPasswordAndHashSame(%s, %s), got %v: expected %v", password, hash, got, expected)
		}
	}
}

func TestDifferentPasswordHaveDifferentHash(t *testing.T) {
	expected := false
	falsePasswords := []string{
		"fake",
		"false",
		"thisPasswordShouldntWork",
		"tryingAnotherFalsePassword",
		"ThisPasswordIsSecureButNotTheSameAsTheRealOnes42069!",
	}
	passwords := []string{
		"test",
		"123",
		"boo",
		"securepassword",
		"sigmarizz",
		"noOneWillGuessThisPasswordEver42069!",
	}
	for _, password := range passwords {
		hash, err := utils.GeneratePasswordHash(password)
		if err != nil {
			t.Errorf("GeneratePasswordHash(%s) failed to generate hash: %v", password, err)
			continue
		}
		for _, falsePassword := range falsePasswords {
			got := utils.IsPasswordAndHashSame(falsePassword, hash)
			if got != expected {
				t.Errorf("IsPasswordAndHashSame(%s, %s), got %v: expected %v", password, hash, got, expected)
			}
		}
	}
}
