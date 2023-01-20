
-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.shared_credential(
    id BIGSERIAL  NOT NULL PRIMARY KEY,
    data VARCHAR (5000) NOT NULL,
    user_id VARCHAR (50) NOT NULL,
    password VARCHAR (255) NOT NULL,
    expired_at TIMESTAMP  NOT NULL,
    max_number_queries INTEGER  NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.shared_credential;
