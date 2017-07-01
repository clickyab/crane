
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE campaign
(
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    name VARCHAR(20) NOT NULL,
    max_bid INT NOT NULL,
    frequency INT DEFAULT 2 NOT NULL,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null,
    CONSTRAINT campaign_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX campain_user_id_uindex ON ads (user_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE campaign;
