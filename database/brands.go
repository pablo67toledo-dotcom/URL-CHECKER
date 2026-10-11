package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Brand representa uma marca cadastrada no banco de dados.
type Brand struct {
	ID        int64
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateBrand cadastra uma marca e retorna o registro criado.
func CreateBrand(
	ctx context.Context,
	pool *pgxpool.Pool,
	name string,
	slug string,
) (Brand, error) {
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)

	if name == "" {
		return Brand{}, fmt.Errorf("o nome da marca é obrigatório")
	}

	if slug == "" {
		return Brand{}, fmt.Errorf("o slug da marca é obrigatório")
	}

	var brand Brand

	err := pool.QueryRow(
		ctx,
		`INSERT INTO brands (name, slug)
		 VALUES ($1, $2)
		 RETURNING id, name, slug, created_at, updated_at`,
		name,
		slug,
	).Scan(
		&brand.ID,
		&brand.Name,
		&brand.Slug,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)

	if err != nil {
		return Brand{}, fmt.Errorf("não foi possível cadastrar a marca: %w", err)
	}

	return brand, nil
}

// GetBrandBySlug busca uma marca pelo seu identificador textual.
func GetBrandBySlug(
	ctx context.Context,
	pool *pgxpool.Pool,
	slug string,
) (Brand, error) {
	slug = strings.TrimSpace(slug)

	if slug == "" {
		return Brand{}, fmt.Errorf("o slug da marca é obrigatório")
	}

	var brand Brand

	err := pool.QueryRow(
		ctx,
		`SELECT id, name, slug, created_at, updated_at
		 FROM brands
		 WHERE slug = $1`,
		slug,
	).Scan(
		&brand.ID,
		&brand.Name,
		&brand.Slug,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Brand{}, fmt.Errorf("marca não encontrada: %s: %w", slug, pgx.ErrNoRows)
		}

		return Brand{}, fmt.Errorf("não foi possível consultar a marca: %w", err)
	}

	return brand, nil
}

// ListBrands retorna todas as marcas ordenadas pelo ID.
func ListBrands(
	ctx context.Context,
	pool *pgxpool.Pool,
) ([]Brand, error) {
	rows, err := pool.Query(
		ctx,
		`SELECT id, name, slug, created_at, updated_at
		 FROM brands
		 ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("não foi possível listar as marcas: %w", err)
	}
	defer rows.Close()

	brands := make([]Brand, 0)

	for rows.Next() {
		var brand Brand

		if err := rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.Slug,
			&brand.CreatedAt,
			&brand.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("não foi possível ler uma marca: %w", err)
		}

		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a listagem das marcas: %w", err)
	}

	return brands, nil
}
