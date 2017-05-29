
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE supplier_report
(
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    supplier VARCHAR(25) NOT NULL,
    target_date DATE NOT NULL,
    impression_out INT DEFAULT 0 NOT NULL,
    impression_in INT DEFAULT 0 NOT NULL,
    delivered_impression INT DEFAULT 0 NOT NULL,
    earn INT DEFAULT 0 NOT NULL
);
CREATE UNIQUE INDEX supplier_report_unique_gp ON supplier_report (supplier, target_date);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE supplier_report;
