
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

 CREATE TABLE sup_dem_src
 (
   id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
   demand VARCHAR(25),
   supplier VARCHAR(25),
   source VARCHAR(50),
   time_id INT,
   request_out_count INT COMMENT 'total http request comes from exchange',
   imp_in_count INT COMMENT 'total imp comes into exchange',
   imp_out_count INT COMMENT 'total imp comes from exchange',
   win_count INT COMMENT 'total win count',
   win_bid INT COMMENT 'total win bid of winner request',
   deliver_count INT COMMENT 'total show count of winner request',
   deliver_bid INT COMMENT 'total show of winner requests',
   profit INT COMMENT 'profit',
   CONSTRAINT sup_dem_src_time_table_id_fk FOREIGN KEY (time_id) REFERENCES time_table (id)
 );
 CREATE UNIQUE INDEX sup_dem_src_index_name ON sup_dem_src (time_id,demand,supplier,source);



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE sup_dem_src;