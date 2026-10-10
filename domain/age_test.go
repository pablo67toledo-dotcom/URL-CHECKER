package domain

import (
	"testing"
	"time"
)

func TestGetDomainCreationDate(t *testing.T) {
	data, err := GetDomainCreationDate("example.com")
	if err != nil {
		t.Fatalf("erro ao consultar a data de criação: %v", err)
	}

	t.Logf("Data de criação: %s", data.Format("2006-01-02"))
}

func TestCalculateDomainAge(t *testing.T) {
	creationDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	age := CalculateDomainAge(creationDate)

	expectedAge := time.Now().Year() - 2000

	if age != expectedAge {
		t.Errorf("idade incorreta: esperada %d, recebida %d", expectedAge, age)
	}
}
