
-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.credential_page(
    id uuid NOT NULL PRIMARY KEY,
    url VARCHAR (2000) NOT NULL,
    ttl INTEGER  NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.credential_page;
