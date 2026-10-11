# URL Checker

API desenvolvida em Go para analisar URLs, verificar informações relacionadas ao domínio e calcular sua idade. O projeto também possui uma estrutura inicial de persistência de dados utilizando PostgreSQL.

## Tecnologias

- **Go:** desenvolvimento da API.
- **PostgreSQL 17:** banco de dados relacional.
- **Docker Compose:** execução do banco de dados em ambiente local.
- **pgx:** conexão entre a aplicação Go e o PostgreSQL.
- **golang-migrate:** gerenciamento de migrations do banco de dados.
- **godotenv:** carregamento das variáveis de ambiente do arquivo `.env`.

## Funcionalidades implementadas

- Recebimento de URLs pelo endpoint `/analisar`.
- Extração do domínio da URL.
- Verificação de resolução DNS.
- Verificação de redirecionamentos HTTP.
- Consulta da data de criação do domínio via RDAP.
- Cálculo da idade do domínio.
- Conexão com PostgreSQL usando um pool de conexões.
- Estrutura inicial de banco de dados com as tabelas `brands`, `official_domains` e `schema_migrations`.

**Importante:** a verificação de domínio e idade não garante, isoladamente, que uma URL seja segura ou maliciosa. A aplicação está em desenvolvimento.

## Pré-requisitos

- Go instalado.
- Docker Desktop com Docker Compose.
- Git.
- PowerShell no Windows.

## Configuração do ambiente local

### 1. Clonar o repositório

```powershell
git clone https://github.com/pablo67toledo-dotcom/URL_CHECKER.git
cd URL_CHECKER
```

### 2. Configurar as variáveis de ambiente

Crie seu arquivo `.env` a partir do exemplo:

```powershell
Copy-Item .env.example .env
```

Abra o arquivo `.env` e configure uma senha própria para o PostgreSQL. Não publique esse arquivo no GitHub.

### 3. Iniciar o PostgreSQL

Com o Docker Desktop em execução:

```powershell
docker compose up -d
```

Verifique o estado dos contêineres:

```powershell
docker compose ps
```

### 4. Executar os testes

```powershell
go test ./...
```

### 5. Iniciar a API

```powershell
go run .
```

Quando a conexão com o banco for estabelecida, a aplicação deverá informar que o servidor está disponível em:

`http://localhost:8080`

O endpoint de análise existente é:

`http://localhost:8080/analisar?url=https://example.com`

Você pode abrir essa URL no navegador para testar a rota.

## Banco de dados e migrations

As migrations versionadas definem a estrutura inicial do banco de dados, incluindo as tabelas `brands` e `official_domains`.

As migrations devem ser aplicadas ao banco local antes de utilizar as tabelas. Consulte os arquivos do diretório `migrations/` e as instruções do projeto para executar esse procedimento.

## Segurança e configuração

- Não versione o arquivo `.env`.
- Não publique senhas, tokens ou outras credenciais.
- Use `.env.example` somente para documentar as variáveis necessárias.
- Considere os resultados da análise como indicadores, não como uma garantia de segurança.

## Status do projeto

O projeto está em desenvolvimento. A API, a conexão com PostgreSQL e a estrutura inicial do banco estão sendo construídas de forma incremental.

## Licença

A licença do projeto ainda não foi definida.
