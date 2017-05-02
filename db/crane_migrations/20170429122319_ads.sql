
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table ads
(
  id int not null auto_increment
    primary key,
  type tinyint not null,
  width int not null,
  height int not null,
  created_at timestamp default CURRENT_TIMESTAMP not null,
  updated_at timestamp default CURRENT_TIMESTAMP not null,
  active bit default b'1' not null,
  user_id int not null,
  url text not null,
  constraint ads_users_id_fk
    foreign key (user_id) references users (id)
)
;

create index ads_users_id_fk
  on ads (user_id)
;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE ads
