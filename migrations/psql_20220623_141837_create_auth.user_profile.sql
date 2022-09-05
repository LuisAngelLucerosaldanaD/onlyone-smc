
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.user_profile(
    id BIGSERIAL  NOT NULL PRIMARY KEY,
    user_id VARCHAR (100) NOT NULL,
    name VARCHAR (100) NOT NULL,
    path VARCHAR (100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS auth.user_profile;
