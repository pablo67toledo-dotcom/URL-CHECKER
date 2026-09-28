package domain

import (
	"net"
)

// Check realiza a análise de um domínio.
func Check(domainName string) error {
	_, err := net.LookupHost(domainName)

	if err != nil {
		return err
	}

	return nil
}
