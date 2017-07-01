
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table publishers
(
  id int not null auto_increment
    primary key,
  user_id int not null,
  floor_cpm int not null,
  soft_floor_cpm int null,
  name varchar(60) not null,
  bid_type tinyint not null,
  under_floor ENUM('yes', 'no') NOT NULL,
  platform ENUM('app', 'web', 'vast') NOT NULL,
  active ENUM('yes', 'no') NOT NULL,
  created_at timestamp default CURRENT_TIMESTAMP not null,
  updated_at timestamp default CURRENT_TIMESTAMP not null,
  CONSTRAINT publisher_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX publisher_user_uindex ON publishers (user_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE publishers;

