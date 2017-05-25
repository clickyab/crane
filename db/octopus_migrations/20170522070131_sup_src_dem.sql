
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

 CREATE TABLE sup_dem_src
 (
   id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
   demand VARCHAR(25),
   supplier VARCHAR(25),
   source VARCHAR(50),
   time_id INT,
   imp_bid INT COMMENT 'imp_bid is the price of biding for all ads',
   show_bid INT COMMENT 'show_bid is same as imp_bid but only for shown ads',
   win_bid INT COMMENT 'win_bid is same as imp_bid but only for shown ads',
   show_count INT COMMENT 'returns the total time of an ad been showed',
   win_count INT COMMENT 'win returns the total times of an ad win, shown + not shown',
   impression_count INT COMMENT 'impression returns the times of we request a demand for each slot',
   request_count INT COMMENT 'impression returns the times of we request a demand, by default its less or equal to impression',
   CONSTRAINT sup_dem_src_time_table_id_fk FOREIGN KEY (time_id) REFERENCES time_table (id)
 );
 CREATE UNIQUE INDEX sup_dem_src_index_name ON sup_dem_src (time_id,demand,supplier,source);



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE sup_dem_src;