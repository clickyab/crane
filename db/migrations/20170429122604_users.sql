
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table users
(
  id int not null auto_increment
    primary key,
  email varchar(60) not null,
  domain tinyint not null,
  active bit default b'1' not null,
  created_at timestamp default CURRENT_TIMESTAMP not null,
  updated_at timestamp default CURRENT_TIMESTAMP not null,
  constraint users_email_uindex
    unique (email)
)
;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE users

