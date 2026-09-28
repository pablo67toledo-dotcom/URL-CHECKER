package domain

import "fmt"

// Check realiza a análise de um domínio.
func Check(domainName string) string {
	return fmt.Sprintf("Consultando o domínio: %s", domainName)
}
