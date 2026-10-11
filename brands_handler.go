package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pablo67toledo-dotcom/URL_CHECKER/database"
)

type createBrandRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type brandResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// brandHandler recebe o pool existente e atende às requisições de marcas.
func brandHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{
				Error: "método não permitido; utilize POST",
			})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		var input createBrandRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&input); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Error: "JSON inválido ou requisição muito grande",
			})
			return
		}

		input.Name = strings.TrimSpace(input.Name)
		input.Slug = strings.TrimSpace(input.Slug)

		if input.Name == "" || input.Slug == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Error: "os campos name e slug são obrigatórios",
			})
			return
		}

		brand, err := database.CreateBrand(
			r.Context(),
			pool,
			input.Name,
			input.Slug,
		)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				writeJSON(w, http.StatusConflict, errorResponse{
					Error: "já existe uma marca com esse slug",
				})
				return
			}

			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error: "não foi possível cadastrar a marca",
			})
			return
		}

		writeJSON(w, http.StatusCreated, brandResponse{
			ID:   brand.ID,
			Name: brand.Name,
			Slug: brand.Slug,
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
