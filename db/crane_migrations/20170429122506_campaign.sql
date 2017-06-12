
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table campaign
(
  id int not null auto_increment
    primary key,
  user_id int not null,
  name varchar(60) not null,
  max_bid int not null,
  frequency int not null,
  created_at timestamp default CURRENT_TIMESTAMP not null,
  updated_at timestamp default CURRENT_TIMESTAMP not null,
  constraint campaign_users_id_fk
    foreign key (user_id) references users (id)
)
;

create index campaign_users_id_fk
  on campaign (user_id)
;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE campaign;
