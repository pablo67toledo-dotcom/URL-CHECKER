package domain

import (
	"net"
	"net/http"
)

// Check da análise de um domínio e redirecionamento de uma URL.
//
//	Retorna um erro caso o domínio seja inválido ou não possa ser encontrado.
func Check(domainName string) error {
	_, err := net.LookupHost(domainName)

	if err != nil {
		return err
	}

	return nil
}

func CheckRedirect(urlRecebida string) (string, error) {
	redirectURL := ""
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resposta, err := client.Get(urlRecebida)
	if err != nil {
		return "", err
	}

	defer resposta.Body.Close()

	if resposta.StatusCode >= 300 && resposta.StatusCode < 400 {
		redirectURL = resposta.Header.Get("Location")
	}

	return redirectURL, nil

}
