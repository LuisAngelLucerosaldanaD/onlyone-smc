
-- +migrate Up
CREATE TABLE cfg.credentials_styles(
    [id] UNIQUEIDENTIFIER NOT NULL PRIMARY KEY,
    [type] [INT]  NOT NULL,
    [background] [VARCHAR] (255) NOT NULL,
    [logo] [VARCHAR] (255) NOT NULL,
    [front] [VARCHAR] (255) NOT NULL,
    [back] [VARCHAR] (255) NOT NULL,
    [credential_id] [UNIQUEIDENTIFIER]  NOT NULL,
    created_at [datetime] NOT NULL DEFAULT (getdate()),
    updated_at [datetime] NOT NULL DEFAULT (getdate())
);

-- +migrate Down
DROP TABLE cfg.credentials_styles;
