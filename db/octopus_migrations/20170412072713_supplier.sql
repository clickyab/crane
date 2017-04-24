
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE suppliers
(
  id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  name VARCHAR(40) NOT NULL,
  type VARCHAR(10),
  `key` VARCHAR(50),
  floor_cpm INT NOT NULL,
  soft_floor_cpm INT NOT NULL,
  under_floor INT NOT NULL,
  excluded_demands TEXT,
  share INT NOT NULL,
  active INT DEFAULT 1 NOT NULL
);
CREATE UNIQUE INDEX suppliers_name_uindex ON suppliers (name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE suppliers CASCADE ;

