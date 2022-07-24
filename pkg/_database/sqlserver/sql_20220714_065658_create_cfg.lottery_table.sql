
-- +migrate Up
CREATE TABLE cfg.lottery_table(
    [id] UNIQUEIDENTIFIER NOT NULL PRIMARY KEY,
    [block_id] [BIGINT]  NOT NULL,
    [registration_start_date] [TIMESTAMP]  NOT NULL,
    [registration_end_date] [TIMESTAMP]  NOT NULL,
    [lottery_start_date] [TIMESTAMP]  NOT NULL,
    [lottery_end_date] [TIMESTAMP]  NOT NULL,
    [process_end_date] [TIMESTAMP]  NOT NULL,
    [process_status] [INT]  NOT NULL,
    created_at [datetime] NOT NULL DEFAULT (getdate()),
    updated_at [datetime] NOT NULL DEFAULT (getdate())
);

-- +migrate Down
DROP TABLE cfg.lottery_table;
