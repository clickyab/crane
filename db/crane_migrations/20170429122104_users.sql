
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users
(
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    email VARCHAR(25),
    domain VARCHAR(25) DEFAULT 'clickyab.com' NOT NULL,
    password VARCHAR(60) NOT NULL,
    active ENUM('yes', 'no') NOT NULL
);
CREATE UNIQUE INDEX users_email_uindex ON users (email);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE users;

