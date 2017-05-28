
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE sup_src
(
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    supplier VARCHAR(25),
    source VARCHAR(50),
    time_id INT,
    request_in_count INT COMMENT 'total http request comes to exchange',
    imp_in_count INT COMMENT 'total imp comes into exchange',
    imp_out_count INT COMMENT 'total imp comes from exchange',
    deliver_count INT COMMENT 'total show count of winner request',
    deliver_bid INT COMMENT 'total show of winner requests',
    profit INT COMMENT 'profit',
    CONSTRAINT sup_src_time_table_id_fk FOREIGN KEY (time_id) REFERENCES time_table (id)
);
CREATE UNIQUE INDEX sup_src_index_gp ON sup_src (time_id,supplier,source);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE sup_src;
