
-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.user_profile(
    id BIGSERIAL  NOT NULL PRIMARY KEY,
    user_id VARCHAR (100) NOT NULL,
    name VARCHAR (100) NOT NULL,
    path VARCHAR (100) NOT NULL,
    created_at DATE DEFAULT sysdate NOT NULL,
    updated_at DATE DEFAULT sysdate NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS auth.user_profile;
