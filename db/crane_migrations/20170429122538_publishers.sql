
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
  under_floor bit default b'0' null,
  platform tinyint not null,
  active bit default b'1' not null,
  created_at timestamp default CURRENT_TIMESTAMP not null,
  updated_at timestamp default CURRENT_TIMESTAMP not null,
  constraint publishers_users_id_fk
    foreign key (user_id) references users (id)
)
;

create index publishers_users_id_fk
  on publishers (user_id)
;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE publishers;

