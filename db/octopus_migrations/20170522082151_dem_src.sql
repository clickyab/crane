
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE dem_src
(
  id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  demand VARCHAR(25),
  source VARCHAR(50),
  time_id INT,
  request INT COMMENT 'request returns times of requests to a demand',
  win INT COMMENT 'request returns times of requests to demand that has won',
  win_bid INT COMMENT 'request returns total price of requests to a demand that has won',
  `show` INT COMMENT 'request returns times of won ads hat has been shown',
  show_bid INT COMMENT 'request returns total price of won ads hat has been shown',
  CONSTRAINT dem_src_time_table_id_fk FOREIGN KEY (time_id) REFERENCES time_table (id)
);
CREATE UNIQUE INDEX dem_src_id_uindex ON dem_src (id);
CREATE UNIQUE INDEX dem_src_index_gp ON sup_dem_src (demand,time_id)

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE dem_src;