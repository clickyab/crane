
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

 CREATE TABLE demand_report
 (
   id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
   demand VARCHAR(25) NOT NULL ,
   target_date DATE  NOT NULL ,
   request_out_count INT  DEFAULT 0  NOT NULL  COMMENT 'total http request comes from exchange',
   imp_in_count INT DEFAULT 0  NOT  NULL COMMENT 'total imp comes into exchange',
   imp_out_count INT DEFAULT 0  NOT NULL COMMENT 'total imp comes from exchange',
   win_count INT DEFAULT 0  NOT NULL COMMENT 'total win count',
   win_bid INT DEFAULT 0  NOT NULL  COMMENT 'total win bid of winner request',
   deliver_count INT DEFAULT 0  NOT NULL COMMENT 'total show count of winner request',
   deliver_bid INT DEFAULT 0  NOT NULL COMMENT 'total show of winner requests',
   profit INT DEFAULT 0  NOT NULL COMMENT 'profit'
 );
 CREATE UNIQUE INDEX demand_report_demand_time ON demand_report (demand,target_date DESC);
 CREATE INDEX idx_demand_report_target_date ON demand_report (target_date DESC);



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE demand_report;