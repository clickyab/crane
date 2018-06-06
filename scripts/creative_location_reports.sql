
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
USE `clickyab`;

CREATE TABLE IF NOT EXISTS `seats` (
`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
`slot_id` INT(10) UNSIGNED NOT NULL,
`supplier_name` VARCHAR(200) NOT NULL,
`publisher_id` INT(10) UNSIGNED NOT NULL,
`publisher_domain` VARCHAR(200) NOT NULL,
`creative_size` INT(10) UNSIGNED NOT NULL,
`kind` ENUM('web', 'app') NOT NULL,
`active_days` INT(10) UNSIGNED NOT NULL DEFAULT 1,
`avg_daily_imp` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`avg_daily_clicks` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`today_imp` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`today_clicks` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`today_ctr` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (`id`),
CONSTRAINT slot_per_sup_unique UNIQUE (slot_id, supplier_name, publisher_domain, creative_size))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `publisher_pages` (
`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
`publisher_id` INT(10) UNSIGNED NOT NULL,
`publisher_domain` VARCHAR(200) NOT NULL,
`kind` ENUM('web', 'app') NOT NULL,
`url` TEXT NOT NULL,
`url_key` VARCHAR(255) NOT NULL,
`active_days` INT(10) UNSIGNED NOT NULL DEFAULT 1,
`avg_daily_imp` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`avg_daily_clicks` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`today_imp` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`today_clicks` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`today_ctr` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (`id`),
CONSTRAINT pub_page_unique UNIQUE (url_key, publisher_id))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `creatives_locations_reports` (
`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
`publisher_id` INT(10) UNSIGNED NOT NULL,
`publisher_domain` VARCHAR(200) NOT NULL,
`seat_id` INT(10) UNSIGNED NOT NULL,
`publisher_page_id` INT(10) UNSIGNED NOT NULL,
`url_key` VARCHAR(255) NOT NULL,
`creative_id` INT(10) UNSIGNED NOT NULL,
`creative_size` INT(10) UNSIGNED NOT NULL,
`creative_type` INT(10) UNSIGNED NULL,
`active_days` INT(10) UNSIGNED NOT NULL DEFAULT 1,
`total_imp` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`total_clicks` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`total_ctr` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`yesterday_imp` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`yesterday_clicks` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`yesterday_ctr` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`today_imp` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`today_clicks` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`today_ctr` INT(10) UNSIGNED NOT NULL DEFAULT 0,
`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

ALTER TABLE `creatives_locations_reports` ADD INDEX `publisher_domain` (`publisher_domain`) USING BTREE;
ALTER TABLE `creatives_locations_reports` ADD INDEX `seat_id` (`seat_id`) USING BTREE;
ALTER TABLE `creatives_locations_reports` ADD INDEX `publisher_page_id` (`publisher_page_id`) USING BTREE;
ALTER TABLE `creatives_locations_reports` ADD INDEX `url_key` (`url_key`) USING BTREE;
ALTER TABLE `creatives_locations_reports` ADD INDEX `creative_id` (`creative_id`) USING BTREE;
ALTER TABLE `creatives_locations_reports` ADD INDEX `creative_size` (`creative_size`) USING BTREE;