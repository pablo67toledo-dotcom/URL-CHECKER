package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
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

func brandHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createBrand(w, r, pool)
		case http.MethodGet:
			getBrands(w, r, pool)
		default:
			w.Header().Set("Allow", "GET, POST")
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{
				Error: "método não permitido; utilize GET ou POST",
			})
		}
	}
}

func createBrand(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
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
		r.Context(), pool, input.Name, input.Slug,
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
		ID: brand.ID, Name: brand.Name, Slug: brand.Slug,
	})
}

func getBrands(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	slug := strings.TrimSpace(r.URL.Query().Get("slug"))

	if slug == "" {
		brands, err := database.ListBrands(r.Context(), pool)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error: "não foi possível listar as marcas",
			})
			return
		}

		response := make([]brandResponse, 0, len(brands))
		for _, brand := range brands {
			response = append(response, brandResponse{
				ID: brand.ID, Name: brand.Name, Slug: brand.Slug,
			})
		}

		writeJSON(w, http.StatusOK, response)
		return
	}

	brand, err := database.GetBrandBySlug(r.Context(), pool, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Error: "marca não encontrada",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: "não foi possível consultar a marca",
		})
		return
	}

	writeJSON(w, http.StatusOK, brandResponse{
		ID: brand.ID, Name: brand.Name, Slug: brand.Slug,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
