
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE demands ADD COLUMN user_id INT(11);
ALTER TABLE demands ADD CONSTRAINT demand_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE suppliers ADD COLUMN user_id INT(11);
ALTER TABLE suppliers ADD CONSTRAINT supplier_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE users ADD COLUMN user_type ENUM('demand','supplier','admin') NOT NULL;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE demands DROP COLUMN user_id;
ALTER TABLE suppliers DROP COLUMN user_id;
ALTER TABLE users DROP COLUMN user_type;

