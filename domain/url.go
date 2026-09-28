package domain

import (
	"errors"
	"net/url"
)

// Extrai o domínio de uma URL válida.
func ExtractDomain(urlRecebida string) (string, error) {
	urlAnalisada, err := url.Parse(urlRecebida)
	if err != nil {
		return "", errors.New("não foi possível analisar a URL")
	}

	dominio := urlAnalisada.Hostname()

	if dominio == "" {
		return "", errors.New("não foi possível identificar o domínio da URL")
	}

	return dominio, nil
}
