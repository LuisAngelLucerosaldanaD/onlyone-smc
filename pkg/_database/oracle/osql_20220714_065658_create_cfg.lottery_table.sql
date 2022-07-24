
-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.lottery_table(
    id VARCHAR2(50) NOT NULL PRIMARY KEY,
    block_id NUMBER(20,0)  NOT NULL,
    registration_start_date TIMESTAMP  NOT NULL,
    registration_end_date TIMESTAMP  NOT NULL,
    lottery_start_date TIMESTAMP  NOT NULL,
    lottery_end_date TIMESTAMP  NOT NULL,
    process_end_date TIMESTAMP  NOT NULL,
    process_status NUMBER(10,0)  NOT NULL,
    created_at DATE DEFAULT sysdate NOT NULL,
    updated_at DATE DEFAULT sysdate NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.lottery_table;
