
-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.credentials_styles(
    id uuid NOT NULL PRIMARY KEY,
    type INTEGER  NOT NULL,
    background VARCHAR (255) NOT NULL,
    logo VARCHAR (255) NOT NULL,
    front VARCHAR (255) NOT NULL,
    back VARCHAR (255) NOT NULL,
    credential_id UUID  NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.credentials_styles;
