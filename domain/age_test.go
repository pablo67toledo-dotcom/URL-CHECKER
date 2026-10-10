package domain

import "testing"

func TestGetDomainCreationDate(t *testing.T) {
	data, err := GetDomainCreationDate("example.com")
	if err != nil {
		t.Fatalf("erro ao consultar a data de criação: %v", err)
	}

	t.Logf("Data de criação: %s", data.Format("2006-01-02"))

}
