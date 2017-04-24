
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE demands
(
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(40) NOT NULL,
  type VARCHAR(40) NOT NULL,
  get_url VARCHAR(300) NOT NULL,
  win_url VARCHAR(300) NOT NULL,
  white_countrie TEXT,
  excluded_suppliers TEXT,
  minute_limit INT,
  hour_limit INT,
  day_limit INT,
  week_limit INT,
  month_limit INT,
  idle_connection INT,
  timeout INT,
  call_rate INT NOT NULL,
  active INT DEFAULT 1 NOT NULL,
  handicap INT DEFAULT 100 NOT NULL,
  share INT DEFAULT 100 NOT NULL
);
CREATE UNIQUE INDEX demands_name_uindex ON demands (name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE demands CASCADE ;
