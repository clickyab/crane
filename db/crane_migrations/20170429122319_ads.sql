
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE ads
(
    id INT NOT NULL AUTO_INCREMENT primary key,
    type ENUM('native', 'video', 'banner', 'dynamic') NOT NULL,
    width INT NOT NULL,
    height INT NOT NULL,
    campain_id INT,
    user_id INT NOT NULL,
    active ENUM('yes', 'no') NOT NULL,
    url VARCHAR(200) NOT NULL,
    attribute TEXT,
    updated_at TIMESTAMP DEFAULT current_timestamp NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp NOT NULL,
    CONSTRAINT ads_user_id_fk FOREIGN KEY (user_id) REFERENCES ads (id),
    CONSTRAINT ads_campaign_id_fk FOREIGN KEY (campain_id) REFERENCES ads (id)
);

CREATE INDEX ads_user_uindex ON ads (user_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE ads;
