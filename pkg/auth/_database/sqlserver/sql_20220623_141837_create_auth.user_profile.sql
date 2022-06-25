-- +migrate Up
CREATE TABLE auth.user_profile(
    [id] BIGINT IDENTITY(1,1) PRIMARY KEY,
    [user_id] [VARCHAR] (100) NOT NULL,
    [name] [INT] (100) NOT NULL,
    [path] [VARCHAR] (100) NOT NULL,
    created_at [datetime] NOT NULL DEFAULT (getdate()),
    updated_at [datetime] NOT NULL DEFAULT (getdate())
);

-- +migrate Down
DROP TABLE auth.user_profile;
