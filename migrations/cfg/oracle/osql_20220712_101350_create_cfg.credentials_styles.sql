
-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.credentials_styles(
    id VARCHAR2(50) NOT NULL PRIMARY KEY,
    type NUMBER(10,0)  NOT NULL,
    background VARCHAR2 (255) NOT NULL,
    logo VARCHAR2 (255) NOT NULL,
    front VARCHAR2 (255) NOT NULL,
    back VARCHAR2 (255) NOT NULL,
    credential_id VARCHAR2(50)  NOT NULL,
    created_at DATE DEFAULT sysdate NOT NULL,
    updated_at DATE DEFAULT sysdate NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.credentials_styles;
