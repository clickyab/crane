
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE exchange_report
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    target_date DATE NOT NULL,
    supplier_impression_in INT DEFAULT 0 NOT NULL,
    supplier_impression_out INT DEFAULT 0 NOT NULL,
    demand_impression_in INT DEFAULT 0 NOT NULL,
    demand_impression_out INT DEFAULT 0 NOT NULL,
    earn INT DEFAULT 0 NOT NULL,
    spent INT DEFAULT 0 NOT NULL,
    income INT DEFAULT 0 NOT NULL
);
CREATE UNIQUE INDEX exchange_date_uindex ON clickyab.exchange (target_date);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


