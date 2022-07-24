
-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.lottery_table(
    id uuid NOT NULL PRIMARY KEY,
    block_id BIGINT  NOT NULL,
    registration_start_date TIMESTAMP  NOT NULL,
    registration_end_date TIMESTAMP  NOT NULL,
    lottery_start_date TIMESTAMP  NOT NULL,
    lottery_end_date TIMESTAMP  NOT NULL,
    process_end_date TIMESTAMP  NOT NULL,
    process_status INTEGER  NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.lottery_table;
