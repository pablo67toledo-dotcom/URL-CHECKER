package domain

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type rdapEvent struct {
	EventAction string `json:"eventAction"`
	EventDate   string `json:"eventDate"`
}

type rdapResponse struct {
	Events []rdapEvent `json:"events"`
}

// GetDomainCreationDate consulta o serviço RDAP para obter a data de criação de um domínio.
func GetDomainCreationDate(domainName string) (time.Time, error) {
	if domainName == "" {
		return time.Time{}, errors.New("o domínio não pode ser vazio")
	}

	rdapURL := "https://rdap.org/domain/" + domainName

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(rdapURL)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return time.Time{}, errors.New("não foi possível consultar os dados do domínio")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, err
	}

	var rdapData rdapResponse
	if err := json.Unmarshal(body, &rdapData); err != nil {
		return time.Time{}, err
	}

	for _, event := range rdapData.Events {
		if event.EventAction == "registration" {
			return time.Parse(time.RFC3339, event.EventDate)
		}
	}

	return time.Time{}, errors.New("data de criação do domínio não encontrada")
}

func CalculateDomainAge(creationDate time.Time) int {
	today := time.Now()

	age := today.Year() - creationDate.Year()

	if today.YearDay() < creationDate.YearDay() {
		age--
	}

	return age
}
