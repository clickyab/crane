
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE sup_src
(
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    supplier VARCHAR(25),
    source VARCHAR(50),
    time_id INT,
    imp_bid INT COMMENT 'imp_bid returns total price of winner requests',
    show_bid INT COMMENT 'imp_bid returns total price of winner requests that has been showed',
    request_count INT COMMENT 'impression returns the times of request comes from supplier',
    impression_count INT COMMENT 'impression returns the times of request comes from supplier for each slot, by default its more or equal to request',
    show_count INT COMMENT 'show returns the number of times impressions has been shown',
    CONSTRAINT sup_src_time_table_id_fk FOREIGN KEY (time_id) REFERENCES time_table (id)
);
CREATE UNIQUE INDEX sup_src_index_gp ON sup_src (time_id,supplier,source);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE sup_src;
