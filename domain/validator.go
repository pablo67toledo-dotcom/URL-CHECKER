package domain

import "errors"

// Verifica se o domínio informado é válido.
func Validate(domainName string) error {
	if domainName == "" {
		return errors.New("o domínio não pode ser vazio")
	}
	return nil
}
