package database

import (
	"context"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestCreateBrandAndGetBrandBySlug(t *testing.T) {
	if testing.Short() {
		t.Skip("teste de integração ignorado no modo short")
	}

	// Carrega .env quando as variáveis não estão definidas no ambiente.
	if err := loadTestEnvironment(); err != nil {
		t.Fatalf("não foi possível carregar o ambiente: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := NewPool(ctx)
	if err != nil {
		t.Fatalf("não foi possível conectar ao PostgreSQL: %v", err)
	}
	defer pool.Close()

	slug := "teste-url-checker-integration"

	// Remove somente o registro reservado para este teste.
	_, err = pool.Exec(ctx, "DELETE FROM brands WHERE slug = $1", slug)
	if err != nil {
		t.Fatalf("não foi possível preparar o teste: %v", err)
	}

	created, err := CreateBrand(ctx, pool, "Marca de Teste", slug)
	if err != nil {
		t.Fatalf("não foi possível criar a marca: %v", err)
	}

	found, err := GetBrandBySlug(ctx, pool, slug)
	if err != nil {
		t.Fatalf("não foi possível consultar a marca: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("ID obtido = %d; esperado = %d", found.ID, created.ID)
	}

	if found.Name != "Marca de Teste" {
		t.Errorf("nome obtido = %q; esperado = %q", found.Name, "Marca de Teste")
	}

	if found.Slug != slug {
		t.Errorf("slug obtido = %q; esperado = %q", found.Slug, slug)
	}
}

func loadTestEnvironment() error {
	if err := godotenv.Load(".env"); err == nil {
		return nil
	}

	return godotenv.Load("../.env")
}
