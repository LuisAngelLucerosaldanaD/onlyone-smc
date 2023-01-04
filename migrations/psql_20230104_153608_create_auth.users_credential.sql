
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.users_credential(
    id uuid NOT NULL PRIMARY KEY,
    private_key VARCHAR (1000) NOT NULL,
    identity_number VARCHAR (100) NOT NULL,
    mnemonic VARCHAR (500) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.users_credential;
