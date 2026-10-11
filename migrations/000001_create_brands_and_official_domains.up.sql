CREATE TABLE brands (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT brands_name_not_empty
        CHECK (length(trim(name)) > 0),

    CONSTRAINT brands_slug_not_empty
        CHECK (length(trim(slug)) > 0)
);

CREATE TABLE official_domains (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    brand_id BIGINT NOT NULL
        REFERENCES brands(id) ON DELETE RESTRICT,
    domain TEXT NOT NULL,
    source_url TEXT NOT NULL,
    verified_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT official_domains_domain_not_empty
        CHECK (length(trim(domain)) > 0),

    CONSTRAINT official_domains_source_url_not_empty
        CHECK (length(trim(source_url)) > 0)
);

CREATE UNIQUE INDEX official_domains_domain_lower_unique
    ON official_domains (LOWER(domain));

CREATE INDEX official_domains_brand_id_idx
    ON official_domains (brand_id);