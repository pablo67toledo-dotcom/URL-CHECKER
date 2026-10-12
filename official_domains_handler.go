package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pablo67toledo-dotcom/URL_CHECKER/database"
)

type OfficialDomainsHandler struct {
	Pool *pgxpool.Pool
}

type createOfficialDomainRequest struct {
	Domain    string `json:"domain"`
	SourceURL string `json:"source_url"`
}

func (h *OfficialDomainsHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 || parts[0] != "brands" || parts[2] != "domains" {
		http.NotFound(w, r)
		return
	}

	brandID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || brandID <= 0 {
		http.Error(w, "ID da empresa inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.create(w, r, brandID)
	case http.MethodGet:
		h.list(w, r, brandID)
	}
}

func (h *OfficialDomainsHandler) create(
	w http.ResponseWriter,
	r *http.Request,
	brandID int64,
) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var req createOfficialDomainRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Domain) == "" ||
		strings.TrimSpace(req.SourceURL) == "" {
		http.Error(
			w,
			"Os campos domain e source_url são obrigatórios",
			http.StatusBadRequest,
		)
		return
	}

	domain, err := database.CreateOfficialDomain(
		r.Context(),
		h.Pool,
		brandID,
		req.Domain,
		req.SourceURL,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				http.Error(
					w,
					"Esse domínio já está cadastrado",
					http.StatusConflict,
				)
				return
			case "23503":
				http.Error(
					w,
					"Empresa não encontrada",
					http.StatusNotFound,
				)
				return
			}
		}

		http.Error(
			w,
			"Erro ao cadastrar domínio oficial",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(domain)
}

func (h *OfficialDomainsHandler) list(
	w http.ResponseWriter,
	r *http.Request,
	brandID int64,
) {
	domains, err := database.ListOfficialDomainsByBrandID(
		r.Context(),
		h.Pool,
		brandID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Empresa não encontrada", http.StatusNotFound)
			return
		}

		http.Error(
			w,
			"Erro ao consultar domínios oficiais",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(domains)
}
