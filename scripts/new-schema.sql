CREATE TABLE impressions
(
  id VARCHAR(63) NOT NULL,
  campaign_id INT NOT NULL,
  creative_id INT NOT NULL,
  cpc INT NOT NULL,
  cpm INT NOT NULL,
  delay INT NOT NULL,
  share INT NOT NULL,
  diff INT DEFAULT NULL ,
  publisher_name VARCHAR(63) NOT NULL,
  publisher_type ENUM("app", "web") NOT NULL,
  supplier VARCHAR(63) NOT NULL,
  seat_id VARCHAR(63) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  CONSTRAINT impressions_shadow_created_at_id_pk PRIMARY KEY (created_at, id)
);

CREATE TABLE impression_details
(
  id VARCHAR(63) NOT NULL,
  ref VARCHAR(127),
  page VARCHAR(127),
  ip VARCHAR(31) NOT NULL,
  ua VARCHAR(255) NOT NULL,
  user_id VARCHAR(63) NOT NULL,
  lat DOUBLE,
  lon DOUBLE,
  created_at TIMESTAMP NOT NULL,
  CONSTRAINT impression_detail_created_at_imp_id_pk PRIMARY KEY (created_at, id)
);

CREATE TABLE clicks
(
  id VARCHAR(63) NOT NULL,
  campaign_id INT NOT NULL,
  creative_id INT NOT NULL,
  cpc INT NOT NULL,
  share INT NOT NULL,
  delay INT NOT NULL,
  publisher_name VARCHAR(63) NOT NULL,
  publisher_type ENUM("app", "web") NOT NULL,
  supplier VARCHAR(63) NOT NULL,
  seat_id VARCHAR(63) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  CONSTRAINT click_shadow_created_at_id_pk PRIMARY KEY (created_at, id)
);

CREATE TABLE click_details
(
  id VARCHAR(63) NOT NULL,
  ref VARCHAR(127),
  page VARCHAR(127),
  ip VARCHAR(31) NOT NULL,
  ua VARCHAR(255) NOT NULL,
  user_id VARCHAR(63) NOT NULL,
  lat DOUBLE,
  lon DOUBLE,
  created_at TIMESTAMP NOT NULL,
  CONSTRAINT click_detail_created_at_imp_id_pk PRIMARY KEY (created_at, id)
);
