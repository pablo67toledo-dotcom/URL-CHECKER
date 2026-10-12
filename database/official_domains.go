package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OfficialDomain representa um domínio oficial cadastrado no banco.
type OfficialDomain struct {
	ID         int64
	BrandID    int64
	Domain     string
	SourceURL  string
	VerifiedAt *time.Time
	IsActive   bool
	CreatedAt  time.Time
}

// CreateOfficialDomain cadastra um domínio oficial para uma empresa.
func CreateOfficialDomain(
	ctx context.Context,
	pool *pgxpool.Pool,
	brandID int64,
	domain string,
	sourceURL string,
) (OfficialDomain, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	domain = strings.TrimSuffix(domain, ".")

	sourceURL = strings.TrimSpace(sourceURL)

	if brandID <= 0 {
		return OfficialDomain{}, fmt.Errorf("brandID deve ser maior que zero")
	}

	if domain == "" {
		return OfficialDomain{}, fmt.Errorf("domínio não pode ser vazio")
	}

	if sourceURL == "" {
		return OfficialDomain{}, fmt.Errorf("sourceURL não pode ser vazio")
	}

	var officialDomain OfficialDomain

	err := pool.QueryRow(
		ctx,
		`
		INSERT INTO official_domains (
			brand_id,
			domain,
			source_url
		)
		VALUES ($1, $2, $3)
		RETURNING id, brand_id, domain, source_url,
		          verified_at, is_active, created_at
		`,
		brandID,
		domain,
		sourceURL,
	).Scan(
		&officialDomain.ID,
		&officialDomain.BrandID,
		&officialDomain.Domain,
		&officialDomain.SourceURL,
		&officialDomain.VerifiedAt,
		&officialDomain.IsActive,
		&officialDomain.CreatedAt,
	)

	if err != nil {
		return OfficialDomain{}, fmt.Errorf(
			"erro ao cadastrar domínio oficial: %w", err,
		)
	}

	return officialDomain, nil
}

// ListOfficialDomainsByBrandID lista os domínios de uma empresa.
func ListOfficialDomainsByBrandID(
	ctx context.Context,
	pool *pgxpool.Pool,
	brandID int64,
) ([]OfficialDomain, error) {
	if brandID <= 0 {
		return nil, fmt.Errorf("brandID deve ser maior que zero")
	}

	rows, err := pool.Query(
		ctx,
		`
		SELECT id, brand_id, domain, source_url,
		       verified_at, is_active, created_at
		FROM official_domains
		WHERE brand_id = $1
		ORDER BY id
		`,
		brandID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao consultar domínios oficiais: %w", err,
		)
	}
	defer rows.Close()

	domains := make([]OfficialDomain, 0)

	for rows.Next() {
		var officialDomain OfficialDomain
		var verifiedAt pgtype.Timestamptz

		err := rows.Scan(
			&officialDomain.ID,
			&officialDomain.BrandID,
			&officialDomain.Domain,
			&officialDomain.SourceURL,
			&verifiedAt,
			&officialDomain.IsActive,
			&officialDomain.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"erro ao ler domínio oficial: %w", err,
			)
		}

		if verifiedAt.Valid {
			officialDomain.VerifiedAt = &verifiedAt.Time
		}

		domains = append(domains, officialDomain)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"erro durante a leitura dos domínios oficiais: %w", err,
		)
	}

	return domains, nil
}

// ListActiveOfficialDomains lista todos os domínios oficiais ativos.
func ListActiveOfficialDomains(
	ctx context.Context,
	pool *pgxpool.Pool,
) ([]OfficialDomain, error) {
	rows, err := pool.Query(
		ctx,
		`
		SELECT id, brand_id, domain, source_url,
		       verified_at, is_active, created_at
		FROM official_domains
		WHERE is_active = TRUE
		ORDER BY brand_id, domain
		`,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao consultar domínios oficiais ativos: %w", err,
		)
	}
	defer rows.Close()

	domains := make([]OfficialDomain, 0)

	for rows.Next() {
		var officialDomain OfficialDomain
		var verifiedAt pgtype.Timestamptz

		err := rows.Scan(
			&officialDomain.ID,
			&officialDomain.BrandID,
			&officialDomain.Domain,
			&officialDomain.SourceURL,
			&verifiedAt,
			&officialDomain.IsActive,
			&officialDomain.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"erro ao ler domínio oficial: %w", err,
			)
		}

		if verifiedAt.Valid {
			officialDomain.VerifiedAt = &verifiedAt.Time
		}

		domains = append(domains, officialDomain)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"erro durante a leitura dos domínios oficiais: %w", err,
		)
	}

	return domains, nil
}

// IsNoRows indica se a consulta não encontrou um registro.
func IsNoRows(err error) bool {
	return err == pgx.ErrNoRows
}
