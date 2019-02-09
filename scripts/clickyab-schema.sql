-- MySQL dump 10.15  Distrib 10.0.36-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: clickyab
-- ------------------------------------------------------
-- Server version	10.0.36-MariaDB-1~xenial

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `clickyab`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `clickyab` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `clickyab`;

--
-- Table structure for table `ads`
--

DROP TABLE IF EXISTS `ads`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ads` (
  `ad_id` int(11) NOT NULL AUTO_INCREMENT,
  `ad_size` tinyint(1) DEFAULT '0',
  `u_id` int(11) DEFAULT '0',
  `ad_name` text CHARACTER SET utf8mb4,
  `ad_url` text CHARACTER SET utf8mb4,
  `ad_code` text CHARACTER SET utf8mb4,
  `ad_title` text CHARACTER SET utf8mb4,
  `ad_body` text CHARACTER SET utf8mb4,
  `ad_img` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_status` tinyint(1) DEFAULT '0',
  `ad_reject_reason` varchar(50) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_ctr` float DEFAULT '0.1',
  `ad_conv` mediumint(6) DEFAULT '0',
  `ad_time` int(11) DEFAULT '0',
  `ad_type` tinyint(1) DEFAULT '0',
  `ad_visiable` tinyint(1) DEFAULT '1',
  `ad_mainText` varchar(128) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_defineText` varchar(128) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_textColor` varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_target` varchar(30) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_attribute` text CHARACTER SET utf8mb4,
  `ad_hash_attribute` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `ad_mime` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`ad_id`),
  KEY `u_id` (`u_id`),
  KEY `ad_size` (`ad_size`,`ad_status`),
  KEY `ad_hash_attribute` (`ad_hash_attribute`(191))
) ENGINE=InnoDB AUTO_INCREMENT=144631 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ads_frequency`
--

DROP TABLE IF EXISTS `ads_frequency`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ads_frequency` (
  `af_id` int(11) NOT NULL AUTO_INCREMENT,
  `ad_id` int(11) DEFAULT '0',
  `cp_id` int(11) DEFAULT '0',
  `cop_id` int(11) DEFAULT '0',
  `af_count_total` mediumint(6) DEFAULT '0',
  `af_page_event` int(11) DEFAULT '0',
  `af_count_today` mediumint(6) DEFAULT '0',
  `af_date` int(8) DEFAULT '0',
  PRIMARY KEY (`af_id`),
  KEY `ad_id` (`ad_id`,`cop_id`),
  KEY `cop_id` (`cop_id`),
  KEY `cop_id_2` (`cop_id`,`af_count_today`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `api_users`
--

DROP TABLE IF EXISTS `api_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_users` (
  `api_users_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_id` int(11) DEFAULT NULL,
  `api_users_password` varchar(128) DEFAULT NULL COMMENT 'MD5',
  `api_users_token` varchar(256) DEFAULT NULL,
  `api_users_token_expire` datetime DEFAULT NULL,
  `api_users_access_table` text NOT NULL,
  PRIMARY KEY (`api_users_id`),
  KEY `u_id` (`u_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `app_categories`
--

DROP TABLE IF EXISTS `app_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_categories` (
  `cat_id` int(11) NOT NULL AUTO_INCREMENT,
  `cat_title` varchar(90) CHARACTER SET utf8 DEFAULT NULL,
  `cat_title_persian` varchar(90) CHARACTER SET utf8 DEFAULT NULL,
  `cat_count_w` int(11) NOT NULL DEFAULT '0',
  `cat_count_a` int(11) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`cat_id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps`
--

DROP TABLE IF EXISTS `apps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps` (
  `app_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_id` int(11) DEFAULT '0',
  `app_token` varchar(100) DEFAULT NULL,
  `app_name` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `app_supplier` varchar(32) CHARACTER SET utf8mb4 NOT NULL DEFAULT 'clickyab',
  `en_app_name` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `app_package` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `am_id` tinyint(1) DEFAULT '0' COMMENT 'Market Id',
  `app_minbid` int(11) DEFAULT '700',
  `app_floor_cpm` int(11) DEFAULT '700',
  `app_div` float DEFAULT '3',
  `app_status` tinyint(1) DEFAULT '0',
  `app_review` tinyint(1) DEFAULT '0' COMMENT '0 => pending,1 => review,2 => repending',
  `app_today_ctr` int(11) DEFAULT '0',
  `app_today_imps` int(11) DEFAULT '0',
  `app_today_clicks` int(11) DEFAULT '0',
  `app_date` int(11) DEFAULT '0' COMMENT 'Date(''Ymd'')',
  `app_cat` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `app_notapprovedreason` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `app_fatfinger` tinyint(4) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `app_prepayment` tinyint(4) DEFAULT '0',
  `app_publish_cost` tinyint(4) DEFAULT '0',
  `app_min_cpc` text,
  PRIMARY KEY (`app_id`),
  UNIQUE KEY `app_token` (`app_token`),
  KEY `u_id` (`u_id`),
  KEY `apps_app_token_app_status_index` (`app_token`,`app_status`)
) ENGINE=InnoDB AUTO_INCREMENT=105537 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps_android_ver`
--

DROP TABLE IF EXISTS `apps_android_ver`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps_android_ver` (
  `aav_id` int(11) NOT NULL AUTO_INCREMENT,
  `aav_version` int(11) DEFAULT '0',
  PRIMARY KEY (`aav_id`),
  UNIQUE KEY `aav_android_version` (`aav_version`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps_brand_models`
--

DROP TABLE IF EXISTS `apps_brand_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps_brand_models` (
  `abm_id` int(11) NOT NULL AUTO_INCREMENT,
  `ab_id` int(11) DEFAULT '0',
  `abm_model` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`abm_id`),
  UNIQUE KEY `ab_id_2` (`ab_id`,`abm_model`),
  KEY `ab_id` (`ab_id`),
  KEY `abm_model` (`abm_model`)
) ENGINE=InnoDB AUTO_INCREMENT=22957 DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps_brands`
--

DROP TABLE IF EXISTS `apps_brands`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps_brands` (
  `ab_id` int(11) NOT NULL AUTO_INCREMENT,
  `ab_brand` varchar(255) DEFAULT NULL,
  `ab_show` tinyint(4) DEFAULT '1',
  `ab_count` int(12) DEFAULT '0',
  PRIMARY KEY (`ab_id`),
  UNIQUE KEY `ab_brand` (`ab_brand`)
) ENGINE=InnoDB AUTO_INCREMENT=221340 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps_carriers`
--

DROP TABLE IF EXISTS `apps_carriers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps_carriers` (
  `ac_id` int(11) NOT NULL AUTO_INCREMENT,
  `ac_carrier` varchar(255) DEFAULT NULL,
  `ac_show` tinyint(4) DEFAULT '1',
  `ac_count` int(12) DEFAULT '0',
  PRIMARY KEY (`ac_id`),
  UNIQUE KEY `ac_carrier` (`ac_carrier`)
) ENGINE=InnoDB AUTO_INCREMENT=3231336 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps_install`
--

DROP TABLE IF EXISTS `apps_install`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps_install` (
  `api_id` int(11) NOT NULL,
  `u_id` int(11) DEFAULT '0',
  `api_token` varchar(200) DEFAULT NULL,
  `api_name` varchar(200) DEFAULT NULL,
  `api_package` varchar(200) DEFAULT NULL,
  `api_status` varchar(200) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps_langs`
--

DROP TABLE IF EXISTS `apps_langs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps_langs` (
  `al_id` int(11) NOT NULL AUTO_INCREMENT,
  `al_lang` varchar(255) DEFAULT NULL,
  `al_show` tinyint(4) DEFAULT '1',
  `al_count` int(12) DEFAULT NULL,
  PRIMARY KEY (`al_id`),
  UNIQUE KEY `al_lang` (`al_lang`)
) ENGINE=InnoDB AUTO_INCREMENT=183805 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps_market`
--

DROP TABLE IF EXISTS `apps_market`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps_market` (
  `am_id` int(11) NOT NULL AUTO_INCREMENT,
  `am_market` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `am_market_persian` varchar(50) CHARACTER SET utf8mb4 DEFAULT NULL,
  `am_market_os` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  PRIMARY KEY (`am_id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps_networks`
--

DROP TABLE IF EXISTS `apps_networks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps_networks` (
  `an_id` int(11) NOT NULL AUTO_INCREMENT,
  `an_network` varchar(255) DEFAULT NULL,
  `an_show` tinyint(4) DEFAULT '1',
  `an_count` int(12) DEFAULT '0',
  PRIMARY KEY (`an_id`),
  UNIQUE KEY `an_network` (`an_network`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `apps_potential`
--

DROP TABLE IF EXISTS `apps_potential`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apps_potential` (
  `send` tinyint(1) DEFAULT '0',
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `unsub` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1056 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `audit_log`
--

DROP TABLE IF EXISTS `audit_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `audit_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `impersonator` int(11) DEFAULT NULL,
  `for_who` int(11) DEFAULT NULL,
  `action` char(30) DEFAULT NULL,
  `target_id` int(10) unsigned DEFAULT NULL,
  `target_type` varchar(255) NOT NULL,
  `description` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `audit_log_target_id_target_type_index` (`target_id`,`target_type`),
  KEY `audit_log_user_id_foreign` (`user_id`),
  KEY `audit_log_impersonator_foreign` (`impersonator`),
  KEY `audit_log_role_id_foreign` (`role_id`),
  KEY `audit_log_for_who_foreign` (`for_who`),
  KEY `audit_log_action_index` (`action`),
  CONSTRAINT `audit_log_for_who_foreign` FOREIGN KEY (`for_who`) REFERENCES `users` (`u_id`),
  CONSTRAINT `audit_log_impersonator_foreign` FOREIGN KEY (`impersonator`) REFERENCES `users` (`u_id`),
  CONSTRAINT `audit_log_role_id_foreign` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `audit_log_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`u_id`)
) ENGINE=InnoDB AUTO_INCREMENT=643936 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `audit_log_details`
--

DROP TABLE IF EXISTS `audit_log_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `audit_log_details` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `audit_id` int(10) unsigned NOT NULL,
  `data` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `audit_log_details_audit_id_foreign` (`audit_id`),
  CONSTRAINT `audit_log_details_audit_id_foreign` FOREIGN KEY (`audit_id`) REFERENCES `audit_log` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=150829 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `billing`
--

DROP TABLE IF EXISTS `billing`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `billing` (
  `bi_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_id` int(11) DEFAULT '0',
  `income_id` int(11) DEFAULT '0',
  `bi_is_crm` tinyint(1) DEFAULT '0',
  `bi_title` varchar(255) DEFAULT NULL,
  `bi_amount` int(11) DEFAULT '0',
  `bi_type` int(11) DEFAULT '0',
  `bi_balance` bigint(20) DEFAULT '0',
  `bi_time` int(11) DEFAULT '0',
  `bi_date` int(11) DEFAULT '0',
  `bi_reason` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `factor_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`bi_id`),
  UNIQUE KEY `u_id_2` (`u_id`,`income_id`,`bi_amount`,`bi_time`) USING BTREE,
  KEY `u_id` (`u_id`),
  KEY `billing_billing_factor_id_fk` (`factor_id`),
  CONSTRAINT `billing_billing_factor_id_fk` FOREIGN KEY (`factor_id`) REFERENCES `billing_factor` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1824701 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `billing_factor`
--

DROP TABLE IF EXISTS `billing_factor`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `billing_factor` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `prefix` varchar(15) NOT NULL,
  `amount` int(11) NOT NULL DEFAULT '0',
  `paid_amount` int(11) NOT NULL DEFAULT '0',
  `vat` int(11) NOT NULL DEFAULT '0',
  `creator` int(11) NOT NULL,
  `for_who` int(11) NOT NULL,
  `discount` int(11) NOT NULL DEFAULT '0',
  `tax` int(11) NOT NULL DEFAULT '0',
  `description` varchar(255) DEFAULT NULL,
  `date` int(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `billing_factor_usersc_u_id_fk` (`creator`),
  KEY `billing_factor_users_u_id_fk` (`for_who`),
  CONSTRAINT `billing_factor_users_u_id_fk` FOREIGN KEY (`for_who`) REFERENCES `users` (`u_id`),
  CONSTRAINT `billing_factor_usersc_u_id_fk` FOREIGN KEY (`creator`) REFERENCES `users` (`u_id`)
) ENGINE=InnoDB AUTO_INCREMENT=724 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `billing_slot`
--

DROP TABLE IF EXISTS `billing_slot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `billing_slot` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `amount` int(11) NOT NULL,
  `date` varchar(30) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  UNIQUE KEY `billing_slot_date_uindex` (`date`)
) ENGINE=InnoDB AUTO_INCREMENT=347148 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns`
--

DROP TABLE IF EXISTS `campaigns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns` (
  `cp_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_type` tinyint(2) DEFAULT '0',
  `cp_billing_type` varchar(4) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_id` int(11) DEFAULT '0',
  `cp_name` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_network` tinyint(1) NOT NULL DEFAULT '0',
  `cp_placement` varchar(2550) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_wfilter` varchar(2550) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_retargeting` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_frequency` int(3) DEFAULT '2',
  `cp_segment_id` int(11) DEFAULT '0',
  `cp_app_brand` varchar(200) DEFAULT NULL,
  `cp_net_provider` varchar(200) DEFAULT NULL,
  `cp_app_lang` varchar(200) DEFAULT NULL,
  `cp_app_market` int(11) DEFAULT NULL,
  `cp_web_mobile` tinyint(1) DEFAULT '0',
  `cp_web` tinyint(1) DEFAULT '0',
  `cp_application` tinyint(1) DEFAULT '0',
  `cp_video` tinyint(1) DEFAULT '0',
  `cp_apps_carriers` varchar(200) DEFAULT NULL,
  `cp_longmap` varchar(200) DEFAULT NULL,
  `cp_latmap` varchar(200) DEFAULT NULL,
  `cp_radius` int(11) DEFAULT '0',
  `cp_opt_ctr` tinyint(1) DEFAULT '0',
  `cp_opt_conv` tinyint(1) DEFAULT '0',
  `cp_opt_br` tinyint(1) DEFAULT '0',
  `cp_gender` tinyint(1) DEFAULT '0',
  `cp_alexa` tinyint(1) DEFAULT '0',
  `cp_fatfinger` tinyint(1) DEFAULT '1',
  `cp_under` tinyint(1) DEFAULT '0',
  `cp_geos` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_region` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_country` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_hoods` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_isp_blacklist` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_cat` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_like_app` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_app` varchar(2550) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_app_filter` varchar(2550) DEFAULT NULL,
  `cp_keywords` text CHARACTER SET utf8mb4,
  `cp_platforms` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_platform_version` varchar(200) DEFAULT NULL,
  `cp_maxbid` int(11) DEFAULT '0',
  `cp_weekly_budget` int(11) DEFAULT '0',
  `cp_daily_budget` int(11) DEFAULT '0',
  `cp_total_budget` int(11) DEFAULT '0',
  `cp_weekly_spend` int(11) DEFAULT '0',
  `cp_total_spend` int(11) DEFAULT '0',
  `cp_today_spend` int(11) DEFAULT '0',
  `cp_clicks` int(11) DEFAULT '0',
  `cp_ctr` float DEFAULT '0',
  `cp_imps` int(11) DEFAULT '0',
  `cp_cpm` int(11) DEFAULT '0',
  `cp_cpa` int(11) DEFAULT '0',
  `cp_cpc` int(11) DEFAULT '0',
  `cp_conv` int(11) DEFAULT '0',
  `cp_conv_rate` float DEFAULT '0',
  `cp_revenue` int(11) DEFAULT '0',
  `cp_roi` int(4) DEFAULT '0',
  `cp_start` int(11) DEFAULT '0',
  `cp_end` int(11) DEFAULT '0',
  `cp_status` int(11) DEFAULT '1',
  `cp_lastupdate` int(11) DEFAULT '0',
  `cp_hour_start` tinyint(4) DEFAULT '0',
  `cp_hour_end` tinyint(4) DEFAULT '24',
  `cp_time_duration` varchar(255) DEFAULT NULL,
  `is_crm` tinyint(4) DEFAULT '0',
  `cp_lock` int(11) NOT NULL DEFAULT '0' COMMENT 'determine if the campaign was created through crm',
  `created_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `cp_isp` varchar(255) DEFAULT NULL,
  `cp_app_brand_name` text,
  `cp_app_carrier_name` text,
  `cp_net_provider_name` text,
  PRIMARY KEY (`cp_id`),
  KEY `cp_lock` (`cp_lock`),
  KEY `u_id_idx` (`u_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=27108 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_ads`
--

DROP TABLE IF EXISTS `campaigns_ads`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_ads` (
  `ca_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `ad_id` int(11) DEFAULT '0',
  `ca_status` tinyint(4) DEFAULT '1',
  `ca_imps` int(11) DEFAULT '0',
  `ca_cpm` int(11) DEFAULT '0',
  `ca_cpc` int(11) DEFAULT '0',
  `ca_clicks` int(11) DEFAULT '0',
  `ca_ctr` float DEFAULT '0.1',
  `ca_conv` tinyint(4) DEFAULT '0',
  `ca_conv_rate` float DEFAULT '0',
  `ca_cpa` int(11) DEFAULT '0',
  `ca_spend` int(11) DEFAULT '0',
  `ca_lastupdate` int(11) DEFAULT '0',
  PRIMARY KEY (`ca_id`),
  KEY `cp_id` (`cp_id`),
  KEY `ad_id` (`ad_id`),
  KEY `cp_id_2` (`cp_id`,`ca_status`)
) ENGINE=InnoDB AUTO_INCREMENT=143201 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_interests`
--

DROP TABLE IF EXISTS `campaigns_interests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_interests` (
  `cpin_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `in_id` int(11) DEFAULT '0',
  PRIMARY KEY (`cpin_id`),
  KEY `cp_id` (`cp_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_keywords`
--

DROP TABLE IF EXISTS `campaigns_keywords`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_keywords` (
  `cpk_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `k_id` int(11) DEFAULT '0',
  PRIMARY KEY (`cpk_id`),
  KEY `cp_id` (`cp_id`)
) ENGINE=InnoDB AUTO_INCREMENT=697160 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_locations`
--

DROP TABLE IF EXISTS `campaigns_locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_locations` (
  `cpl_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `location_id` int(11) DEFAULT '0',
  PRIMARY KEY (`cpl_id`),
  KEY `cp_id` (`cp_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5074 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_new`
--

DROP TABLE IF EXISTS `campaigns_new`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_new` (
  `cp_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_type` tinyint(2) DEFAULT '0',
  `cp_billing_type` varchar(4) DEFAULT 'cpm',
  `u_id` int(11) DEFAULT '0',
  `cp_name` longtext CHARACTER SET utf8,
  `cp_network` tinyint(1) DEFAULT '0',
  `cp_placement` varchar(2550) DEFAULT NULL,
  `cp_wfilter` varchar(2550) DEFAULT '#',
  `cp_retargeting` varchar(255) DEFAULT NULL,
  `cp_frequency` int(3) DEFAULT '2',
  `cp_segment_id` int(11) DEFAULT '0',
  `cp_opt_ctr` tinyint(1) DEFAULT '0',
  `cp_opt_conv` tinyint(1) DEFAULT '0',
  `cp_opt_br` tinyint(1) DEFAULT '0',
  `cp_gender` tinyint(1) DEFAULT '0',
  `cp_alexa` tinyint(1) DEFAULT '0',
  `cp_fatfinger` tinyint(1) DEFAULT '1',
  `cp_under` tinyint(1) DEFAULT '0',
  `cp_geos` varchar(200) DEFAULT NULL,
  `cp_region` varchar(200) DEFAULT NULL,
  `cp_hoods` varchar(200) DEFAULT NULL,
  `cp_isp_blacklist` varchar(200) DEFAULT '#',
  `cp_cat` varchar(200) DEFAULT NULL,
  `cp_like_app` varchar(200) DEFAULT NULL,
  `cp_app` varchar(2550) DEFAULT NULL,
  `cp_app_filter` varchar(2550) DEFAULT NULL,
  `cp_keywords` text,
  `cp_platforms` varchar(100) DEFAULT NULL,
  `cp_platform_version` varchar(200) DEFAULT NULL COMMENT 'aav_id from   `apps_android_ver`  table',
  `cp_maxbid` int(11) DEFAULT '0',
  `cp_weekly_budget` int(11) DEFAULT '0',
  `cp_daily_budget` int(11) DEFAULT '0',
  `cp_total_budget` int(11) DEFAULT '0',
  `cp_weekly_spend` int(11) DEFAULT '0',
  `cp_total_spend` int(11) DEFAULT '0',
  `cp_today_spend` int(11) DEFAULT '0',
  `cp_clicks` int(11) DEFAULT '0',
  `cp_ctr` float DEFAULT '0',
  `cp_imps` int(11) DEFAULT '0',
  `cp_cpm` int(11) DEFAULT '0',
  `cp_cpa` int(11) DEFAULT '0',
  `cp_cpc` int(11) DEFAULT '0',
  `cp_conversions` int(11) DEFAULT '0',
  `cp_revenue` int(11) DEFAULT '0',
  `cp_roi` int(4) DEFAULT '0',
  `cp_start` int(11) DEFAULT '0',
  `cp_end` int(11) DEFAULT '0',
  `cp_status` int(11) DEFAULT '1',
  `cp_lastupdate` int(11) DEFAULT '0',
  `cp_hour_start` tinyint(4) NOT NULL DEFAULT '0',
  `cp_hour_end` tinyint(4) NOT NULL DEFAULT '24',
  `cp_app_brand` varchar(200) DEFAULT NULL COMMENT 'ab_id from apps_brands table',
  `cp_net_provider` varchar(200) DEFAULT NULL COMMENT 'an_id from apps_networks',
  `cp_app_lang` varchar(200) DEFAULT NULL COMMENT 'al_id from `apps_langs` ',
  `cp_app_market` int(11) DEFAULT NULL,
  `cp_web_mobile` tinyint(1) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `cp_apps_carriers` varchar(200) DEFAULT NULL,
  `cp_web` tinyint(1) NOT NULL,
  `cp_application` tinyint(1) NOT NULL,
  `longmap` varchar(200) DEFAULT NULL,
  `latmap` varchar(200) DEFAULT NULL,
  `radius` int(11) DEFAULT NULL,
  PRIMARY KEY (`cp_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_placement`
--

DROP TABLE IF EXISTS `campaigns_placement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_placement` (
  `cpp_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  `cpp_status` int(11) DEFAULT '0',
  PRIMARY KEY (`cpp_id`),
  KEY `cp_id` (`cp_id`),
  KEY `w_id` (`w_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_placement_copy`
--

DROP TABLE IF EXISTS `campaigns_placement_copy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_placement_copy` (
  `cp_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_type` tinyint(2) DEFAULT '0',
  `cp_billing_type` varchar(4) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_id` int(11) DEFAULT '0',
  `cp_name` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_network` tinyint(1) NOT NULL DEFAULT '0',
  `cp_placement` varchar(2550) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_wfilter` varchar(2550) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_retargeting` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_frequency` int(3) DEFAULT '2',
  `cp_segment_id` int(11) DEFAULT '0',
  `cp_app_brand` varchar(200) DEFAULT NULL,
  `cp_net_provider` varchar(200) DEFAULT NULL,
  `cp_app_lang` varchar(200) DEFAULT NULL,
  `cp_app_market` int(11) DEFAULT NULL,
  `cp_web_mobile` tinyint(1) DEFAULT '0',
  `cp_web` tinyint(1) DEFAULT '0',
  `cp_application` tinyint(1) DEFAULT '0',
  `cp_video` tinyint(1) DEFAULT '0',
  `cp_apps_carriers` varchar(200) DEFAULT NULL,
  `cp_longmap` varchar(200) DEFAULT NULL,
  `cp_latmap` varchar(200) DEFAULT NULL,
  `cp_radius` int(11) DEFAULT '0',
  `cp_opt_ctr` tinyint(1) DEFAULT '0',
  `cp_opt_conv` tinyint(1) DEFAULT '0',
  `cp_opt_br` tinyint(1) DEFAULT '0',
  `cp_gender` tinyint(1) DEFAULT '0',
  `cp_alexa` tinyint(1) DEFAULT '0',
  `cp_fatfinger` tinyint(1) DEFAULT '1',
  `cp_under` tinyint(1) DEFAULT '0',
  `cp_geos` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_region` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_country` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_hoods` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_isp_blacklist` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_cat` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_like_app` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_app` varchar(2550) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_app_filter` varchar(2550) DEFAULT NULL,
  `cp_keywords` text CHARACTER SET utf8mb4,
  `cp_platforms` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cp_platform_version` varchar(200) DEFAULT NULL,
  `cp_maxbid` int(11) DEFAULT '0',
  `cp_weekly_budget` int(11) DEFAULT '0',
  `cp_daily_budget` int(11) DEFAULT '0',
  `cp_total_budget` int(11) DEFAULT '0',
  `cp_weekly_spend` int(11) DEFAULT '0',
  `cp_total_spend` int(11) DEFAULT '0',
  `cp_today_spend` int(11) DEFAULT '0',
  `cp_clicks` int(11) DEFAULT '0',
  `cp_ctr` float DEFAULT '0',
  `cp_imps` int(11) DEFAULT '0',
  `cp_cpm` int(11) DEFAULT '0',
  `cp_cpa` int(11) DEFAULT '0',
  `cp_cpc` int(11) DEFAULT '0',
  `cp_conv` int(11) DEFAULT '0',
  `cp_conv_rate` float DEFAULT '0',
  `cp_revenue` int(11) DEFAULT '0',
  `cp_roi` int(4) DEFAULT '0',
  `cp_start` int(11) DEFAULT '0',
  `cp_end` int(11) DEFAULT '0',
  `cp_status` int(11) DEFAULT '1',
  `cp_lastupdate` int(11) DEFAULT '0',
  `cp_hour_start` tinyint(4) DEFAULT '0',
  `cp_hour_end` tinyint(4) DEFAULT '24',
  `cp_time_duration` varchar(255) DEFAULT NULL,
  `is_crm` tinyint(4) DEFAULT '0',
  `cp_lock` int(11) NOT NULL DEFAULT '0' COMMENT 'determine if the campaign was created through crm',
  `created_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `cp_isp` varchar(255) DEFAULT NULL,
  `cp_app_brand_name` text,
  `cp_app_carrier_name` text,
  `cp_net_provider_name` text,
  PRIMARY KEY (`cp_id`),
  KEY `cp_lock` (`cp_lock`),
  KEY `u_id_idx` (`u_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=25560 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_platform`
--

DROP TABLE IF EXISTS `campaigns_platform`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_platform` (
  `cpp_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `platform_id` int(11) DEFAULT '0',
  PRIMARY KEY (`cpp_id`),
  KEY `cp_id` (`cp_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2689 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_retargeting`
--

DROP TABLE IF EXISTS `campaigns_retargeting`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_retargeting` (
  `cpr_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  PRIMARY KEY (`cpr_id`),
  KEY `cp_id` (`cp_id`),
  KEY `w_id` (`w_id`)
) ENGINE=InnoDB AUTO_INCREMENT=440 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `campaigns_segments`
--

DROP TABLE IF EXISTS `campaigns_segments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `campaigns_segments` (
  `cs_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `seg_id` int(11) DEFAULT '0',
  `cs_conversions` int(11) DEFAULT '0',
  `cs_revenue` int(11) DEFAULT '0',
  `cs_lastupdate` int(11) DEFAULT '0',
  PRIMARY KEY (`cs_id`),
  KEY `cp_id` (`cp_id`)
) ENGINE=InnoDB AUTO_INCREMENT=657 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `categories` (
  `cat_id` int(11) NOT NULL AUTO_INCREMENT,
  `cat_title` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cat_title_persian` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cat_count_w` int(11) NOT NULL DEFAULT '0',
  `cat_count_a` int(11) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`cat_id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `categories_old`
--

DROP TABLE IF EXISTS `categories_old`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `categories_old` (
  `cat_id` int(11) NOT NULL AUTO_INCREMENT,
  `cat_code` varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cat_title` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cat_parent` varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cat_title_persian` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `cat_active` int(11) NOT NULL,
  `cat_count_w` int(5) NOT NULL DEFAULT '0',
  `cat_count_a` int(5) NOT NULL DEFAULT '0',
  PRIMARY KEY (`cat_id`)
) ENGINE=InnoDB AUTO_INCREMENT=393 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cities`
--

DROP TABLE IF EXISTS `cities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cities` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `province_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cities_name_provinces_id_unidex` (`name`,`province_id`),
  KEY `cities_provinces_id_fk` (`province_id`),
  CONSTRAINT `cities_provinces_id_fk` FOREIGN KEY (`province_id`) REFERENCES `provinces` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=441 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `clicks`
--

DROP TABLE IF EXISTS `clicks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `clicks` (
  `c_id` int(11) NOT NULL AUTO_INCREMENT,
  `c_winnerbid` int(11) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  `app_id` int(11) DEFAULT '0',
  `wp_id` int(11) DEFAULT '0',
  `cp_id` int(11) DEFAULT '0',
  `ca_id` int(11) DEFAULT '0',
  `slot_id` int(11) DEFAULT '0',
  `sla_id` int(11) DEFAULT '0',
  `ad_id` int(11) DEFAULT '0',
  `cop_id` int(11) DEFAULT '0',
  `imp_id` int(11) DEFAULT '0',
  `creatives_location_id` int(10) unsigned DEFAULT NULL,
  `publisher_page_id` int(10) unsigned DEFAULT NULL,
  `seat_id` int(10) unsigned DEFAULT NULL,
  `c_status` tinyint(2) DEFAULT '0',
  `c_ip` varchar(20) DEFAULT NULL,
  `c_referaddress` varchar(255) DEFAULT NULL,
  `c_parenturl` varchar(255) DEFAULT NULL,
  `c_fast` int(12) DEFAULT '0',
  `c_os` tinyint(1) DEFAULT '0',
  `c_time` int(11) DEFAULT '0',
  `c_date` int(11) DEFAULT '0',
  `ad_size` int(11) NOT NULL DEFAULT '0',
  `reserved_hash` varchar(63) DEFAULT NULL,
  `c_supplier` varchar(20) NOT NULL DEFAULT 'clickyab',
  PRIMARY KEY (`c_id`),
  KEY `c_date` (`c_date`),
  KEY `w_id` (`w_id`,`c_date`),
  KEY `ca_id` (`ca_id`,`c_status`,`c_date`),
  KEY `app_id` (`app_id`,`c_date`),
  KEY `sla_id` (`sla_id`,`c_status`,`c_date`),
  KEY `clicks_reserved_hash_index` (`reserved_hash`),
  KEY `seat_id` (`seat_id`),
  KEY `publisher_page_id` (`publisher_page_id`),
  KEY `creatives_location_id` (`creatives_location_id`)
) ENGINE=InnoDB AUTO_INCREMENT=180244908 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `clicks_conv`
--

DROP TABLE IF EXISTS `clicks_conv`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `clicks_conv` (
  `cc_id` int(11) NOT NULL AUTO_INCREMENT,
  `c_id` int(11) DEFAULT '0',
  `c_winnerbid` int(11) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  `app_id` int(11) DEFAULT '0',
  `wp_id` int(11) DEFAULT '0',
  `cp_id` int(11) DEFAULT '0',
  `ca_id` int(11) DEFAULT '0',
  `slot_id` int(11) DEFAULT '0',
  `sla_id` int(11) DEFAULT '0',
  `ad_id` int(11) DEFAULT '0',
  `cop_id` bigint(20) DEFAULT '0',
  `imp_id` int(11) DEFAULT '0',
  `reserved_hash` varchar(50) DEFAULT NULL,
  `c_status` tinyint(2) DEFAULT '0',
  `c_ip` varchar(20) DEFAULT NULL,
  `c_referaddress` varchar(255) DEFAULT NULL,
  `c_parenturl` varchar(255) DEFAULT NULL,
  `c_ua` text,
  `c_fast` int(12) DEFAULT '0',
  `c_os` tinyint(1) DEFAULT '0',
  `c_time` int(11) DEFAULT '0',
  `c_date` int(11) DEFAULT '0',
  `c_action` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`cc_id`),
  UNIQUE KEY `reserved_hash` (`reserved_hash`),
  KEY `c_date` (`c_date`),
  KEY `app_id` (`app_id`),
  KEY `slot_id` (`slot_id`),
  KEY `w_id` (`w_id`),
  KEY `sla_id_single` (`sla_id`)
) ENGINE=InnoDB AUTO_INCREMENT=836745 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `conversions`
--

DROP TABLE IF EXISTS `conversions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `conversions` (
  `conv_id` int(11) NOT NULL AUTO_INCREMENT,
  `cs_id` int(11) DEFAULT '0',
  `seg_id` int(11) DEFAULT '0',
  `seg_convvalue` int(11) NOT NULL DEFAULT '0',
  `cp_id` int(11) DEFAULT '0',
  `ad_id` int(11) DEFAULT '0',
  `ca_id` int(11) DEFAULT '0',
  `imp_id` int(11) DEFAULT '0',
  `c_id` int(11) DEFAULT '0',
  `wp_id` int(11) DEFAULT '0',
  `cop_id` int(11) DEFAULT '0',
  `conv_time` int(11) DEFAULT '0',
  `conv_date` int(11) DEFAULT '0',
  PRIMARY KEY (`conv_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2603 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cookie_profiles`
--

DROP TABLE IF EXISTS `cookie_profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cookie_profiles` (
  `cop_id` int(11) NOT NULL AUTO_INCREMENT,
  `cop_key` varchar(20) DEFAULT NULL,
  `cop_email` varchar(100) DEFAULT NULL,
  `cop_last_ip` varchar(50) DEFAULT NULL,
  `cop_gender` tinyint(1) DEFAULT '0',
  `cop_alexa` tinyint(1) DEFAULT '0',
  `cop_os` tinyint(1) DEFAULT '0',
  `cop_browser` tinyint(2) DEFAULT '0',
  `cop_city` smallint(4) DEFAULT '0',
  `cop_age` tinyint(1) DEFAULT '0',
  `cop_keywords` text,
  `cop_active_date` int(9) DEFAULT '0',
  PRIMARY KEY (`cop_id`),
  UNIQUE KEY `cop_key` (`cop_key`)
) ENGINE=InnoDB AUTO_INCREMENT=11558948 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cookie_webpages`
--

DROP TABLE IF EXISTS `cookie_webpages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cookie_webpages` (
  `cwp_id` int(11) NOT NULL AUTO_INCREMENT,
  `wp_id` int(11) DEFAULT '0',
  `w_id` int(12) DEFAULT '0',
  `cop_id` int(11) DEFAULT '0',
  `cwp_time` int(11) DEFAULT '0',
  `cwp_date` int(8) DEFAULT '0',
  PRIMARY KEY (`cwp_id`),
  KEY `wp_id` (`wp_id`),
  KEY `cop_id` (`cop_id`),
  KEY `cwp_date` (`cwp_date`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cookie_websites`
--

DROP TABLE IF EXISTS `cookie_websites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cookie_websites` (
  `cw_id` int(11) NOT NULL AUTO_INCREMENT,
  `w_id` int(12) DEFAULT '0',
  `cop_id` int(11) DEFAULT '0',
  `cw_time` int(11) DEFAULT '0',
  `cw_date` int(11) DEFAULT '0',
  PRIMARY KEY (`cw_id`),
  KEY `cop_id` (`cop_id`),
  KEY `w_id` (`w_id`,`cop_id`),
  KEY `cw_date` (`cw_date`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `countries`
--

DROP TABLE IF EXISTS `countries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `countries` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `countries_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country`
--

DROP TABLE IF EXISTS `country`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `iso` char(2) NOT NULL,
  `name` varchar(80) NOT NULL,
  `nicename` varchar(80) NOT NULL,
  `iso3` char(3) DEFAULT NULL,
  `numcode` smallint(6) DEFAULT NULL,
  `phonecode` int(5) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `iso` (`iso`)
) ENGINE=InnoDB AUTO_INCREMENT=240 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `coupons`
--

DROP TABLE IF EXISTS `coupons`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `coupons` (
  `cpn_id` int(11) NOT NULL AUTO_INCREMENT,
  `cpn_code` varchar(16) NOT NULL,
  `cpn_value` int(11) NOT NULL,
  `u_id` int(11) DEFAULT '0',
  `cpn_date_used` int(11) NOT NULL DEFAULT '0',
  `cpn_date_expire` int(11) NOT NULL,
  PRIMARY KEY (`cpn_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1202 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cp_zero`
--

DROP TABLE IF EXISTS `cp_zero`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cp_zero` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `q` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `creatives_locations_reports`
--

DROP TABLE IF EXISTS `creatives_locations_reports`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `creatives_locations_reports` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `publisher_id` int(10) unsigned NOT NULL,
  `publisher_domain` varchar(200) NOT NULL,
  `seat_id` int(10) unsigned NOT NULL,
  `publisher_page_id` int(10) unsigned NOT NULL,
  `url_key` varchar(255) NOT NULL,
  `creative_id` int(10) unsigned NOT NULL,
  `creative_size` int(10) unsigned NOT NULL,
  `creative_type` int(10) unsigned DEFAULT NULL,
  `active_days` int(10) unsigned NOT NULL DEFAULT '1',
  `total_imp` int(10) unsigned NOT NULL DEFAULT '0',
  `total_clicks` int(10) unsigned NOT NULL DEFAULT '0',
  `total_ctr` int(10) unsigned NOT NULL DEFAULT '0',
  `yesterday_imp` int(10) unsigned NOT NULL DEFAULT '0',
  `yesterday_clicks` int(10) unsigned NOT NULL DEFAULT '0',
  `yesterday_ctr` int(10) unsigned NOT NULL DEFAULT '0',
  `today_imp` int(10) unsigned NOT NULL DEFAULT '0',
  `today_clicks` int(10) unsigned NOT NULL DEFAULT '0',
  `today_ctr` int(10) unsigned NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `publisher_domain` (`publisher_domain`) USING BTREE,
  KEY `seat_id` (`seat_id`) USING BTREE,
  KEY `publisher_page_id` (`publisher_page_id`) USING BTREE,
  KEY `url_key` (`url_key`) USING BTREE,
  KEY `creative_id` (`creative_id`) USING BTREE,
  KEY `creative_size` (`creative_size`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5647435 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `crm_emails`
--

DROP TABLE IF EXISTS `crm_emails`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `crm_emails` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ctr_stat`
--

DROP TABLE IF EXISTS `ctr_stat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ctr_stat` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pub_id` int(11) NOT NULL,
  `pub_type` enum('app','web') NOT NULL,
  `date` int(11) NOT NULL,
  `imp_1` int(11) NOT NULL DEFAULT '0',
  `imp_2` int(11) NOT NULL DEFAULT '0',
  `imp_3` int(11) NOT NULL DEFAULT '0',
  `imp_4` int(11) NOT NULL DEFAULT '0',
  `imp_5` int(11) NOT NULL DEFAULT '0',
  `imp_6` int(11) NOT NULL DEFAULT '0',
  `imp_7` int(11) NOT NULL DEFAULT '0',
  `imp_8` int(11) NOT NULL DEFAULT '0',
  `imp_9` int(11) NOT NULL DEFAULT '0',
  `imp_10` int(11) NOT NULL DEFAULT '0',
  `imp_11` int(11) NOT NULL DEFAULT '0',
  `imp_12` int(11) NOT NULL DEFAULT '0',
  `imp_13` int(11) NOT NULL DEFAULT '0',
  `imp_14` int(11) NOT NULL DEFAULT '0',
  `imp_15` int(11) NOT NULL DEFAULT '0',
  `imp_16` int(11) NOT NULL DEFAULT '0',
  `imp_17` int(11) NOT NULL DEFAULT '0',
  `imp_18` int(11) NOT NULL DEFAULT '0',
  `imp_19` int(11) NOT NULL DEFAULT '0',
  `imp_20` int(11) NOT NULL DEFAULT '0',
  `imp_21` int(11) NOT NULL DEFAULT '0',
  `click_1` int(11) NOT NULL DEFAULT '0',
  `click_2` int(11) NOT NULL DEFAULT '0',
  `click_3` int(11) NOT NULL DEFAULT '0',
  `click_4` int(11) NOT NULL DEFAULT '0',
  `click_5` int(11) NOT NULL DEFAULT '0',
  `click_6` int(11) NOT NULL DEFAULT '0',
  `click_7` int(11) NOT NULL DEFAULT '0',
  `click_8` int(11) NOT NULL DEFAULT '0',
  `click_9` int(11) NOT NULL DEFAULT '0',
  `click_10` int(11) NOT NULL DEFAULT '0',
  `click_11` int(11) NOT NULL DEFAULT '0',
  `click_12` int(11) NOT NULL DEFAULT '0',
  `click_13` int(11) NOT NULL DEFAULT '0',
  `click_14` int(11) NOT NULL DEFAULT '0',
  `click_15` int(11) NOT NULL DEFAULT '0',
  `click_16` int(11) NOT NULL DEFAULT '0',
  `click_17` int(11) NOT NULL DEFAULT '0',
  `click_18` int(11) NOT NULL DEFAULT '0',
  `click_19` int(11) NOT NULL DEFAULT '0',
  `click_20` int(11) NOT NULL DEFAULT '0',
  `click_21` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ctr_stat_pub_id_pub_type_date_uindex` (`pub_id`,`pub_type`,`date`)
) ENGINE=InnoDB AUTO_INCREMENT=7893945 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `daily_report`
--

DROP TABLE IF EXISTS `daily_report`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `daily_report` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `supplier` varchar(63) NOT NULL,
  `type` enum('web','app') NOT NULL,
  `publisher` varchar(150) NOT NULL,
  `imps` int(11) NOT NULL DEFAULT '0',
  `cpm` int(11) NOT NULL DEFAULT '0',
  `clicks` int(11) NOT NULL DEFAULT '0',
  `cpc` int(11) NOT NULL DEFAULT '0',
  `date` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `daily_report_supplier_publisher_date_type_uindex` (`supplier`,`publisher`,`date`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=32898423 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `docker`
--

DROP TABLE IF EXISTS `docker`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `docker` (
  `docker_id` int(11) NOT NULL AUTO_INCREMENT,
  `docker_ip_client` varchar(255) DEFAULT NULL,
  `docker_ip_server` varchar(255) DEFAULT NULL,
  `docker_time` int(11) DEFAULT '0',
  PRIMARY KEY (`docker_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `finder_admin`
--

DROP TABLE IF EXISTS `finder_admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `finder_admin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `pass` char(32) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `finder_cells`
--

DROP TABLE IF EXISTS `finder_cells`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `finder_cells` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `cellname` varchar(6) NOT NULL,
  `top_left_lat` varchar(50) NOT NULL,
  `top_left_long` varchar(50) NOT NULL,
  `bottom_left_lat` varchar(50) NOT NULL,
  `bottom_left_long` varchar(50) NOT NULL,
  `bottom_right_lat` varchar(50) NOT NULL,
  `bottom_right_long` varchar(50) NOT NULL,
  `top_right_lat` varchar(50) NOT NULL,
  `top_right_long` varchar(50) NOT NULL,
  `neighborhoods_id` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `theindex` (`top_left_lat`,`bottom_left_lat`,`top_left_long`,`bottom_right_long`)
) ENGINE=InnoDB AUTO_INCREMENT=3601 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `finder_city_parts`
--

DROP TABLE IF EXISTS `finder_city_parts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `finder_city_parts` (
  `id` mediumint(7) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(4) NOT NULL,
  `cellgroup` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `finder_logs`
--

DROP TABLE IF EXISTS `finder_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `finder_logs` (
  `id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned DEFAULT '0',
  `cell_id` int(10) unsigned DEFAULT '0',
  `imei` varchar(50) DEFAULT NULL,
  `android_id` varchar(50) DEFAULT NULL,
  `carrier` varchar(35) DEFAULT NULL,
  `mcc` mediumint(6) unsigned DEFAULT '0',
  `mnc` mediumint(6) unsigned DEFAULT '0',
  `lac` int(10) unsigned DEFAULT '0',
  `cid` int(10) unsigned DEFAULT '0',
  `ip` varchar(16) DEFAULT NULL,
  `l_time` int(10) unsigned DEFAULT '0',
  `locations` varchar(101) DEFAULT NULL,
  `time` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `finder_logs_sdk_old`
--

DROP TABLE IF EXISTS `finder_logs_sdk_old`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `finder_logs_sdk_old` (
  `id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `cell_id` int(10) unsigned DEFAULT '0',
  `android_id` varchar(50) DEFAULT NULL,
  `android_version` varchar(20) DEFAULT NULL,
  `parameters` text,
  `carrier` varchar(35) DEFAULT NULL,
  `mcc` mediumint(6) unsigned DEFAULT '0',
  `mnc` mediumint(6) unsigned DEFAULT '0',
  `lac` int(10) unsigned DEFAULT '0',
  `cid` int(10) unsigned DEFAULT '0',
  `ip` varchar(16) DEFAULT NULL,
  `locations` varchar(101) DEFAULT NULL,
  `time` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `mcc` (`mcc`,`mnc`,`lac`,`cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `finder_logs_sdk_true`
--

DROP TABLE IF EXISTS `finder_logs_sdk_true`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `finder_logs_sdk_true` (
  `id` int(20) unsigned NOT NULL AUTO_INCREMENT,
  `cell_id` int(10) DEFAULT '0',
  `recheck` tinyint(1) DEFAULT '0',
  `android_id` varchar(50) DEFAULT NULL,
  `android_version` varchar(20) DEFAULT NULL,
  `parameters` text,
  `carrier` varchar(35) DEFAULT NULL,
  `mcc` mediumint(6) unsigned DEFAULT '0',
  `mnc` mediumint(6) unsigned DEFAULT '0',
  `lac` int(10) unsigned DEFAULT '0',
  `cid` int(10) unsigned DEFAULT '0',
  `ip` varchar(100) DEFAULT NULL,
  `locations` varchar(255) DEFAULT NULL,
  `time` int(10) unsigned DEFAULT NULL,
  `neighborhoods_id` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `mcc` (`mcc`,`mnc`,`lac`,`cid`),
  KEY `cell_id` (`cell_id`,`locations`)
) ENGINE=InnoDB AUTO_INCREMENT=110413 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `finder_users`
--

DROP TABLE IF EXISTS `finder_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `finder_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `android_id` varchar(50) NOT NULL,
  `imei` varchar(50) NOT NULL,
  `carrier` varchar(50) NOT NULL,
  `brand` varchar(50) NOT NULL,
  `model` varchar(50) NOT NULL,
  `time` int(10) unsigned NOT NULL,
  `ip` varchar(16) NOT NULL,
  `status` tinyint(2) unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=400 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `google_users`
--

DROP TABLE IF EXISTS `google_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `google_users` (
  `google_id` decimal(21,0) NOT NULL,
  `google_name` varchar(60) NOT NULL,
  `google_email` varchar(60) NOT NULL,
  `google_link` varchar(60) NOT NULL,
  `google_picture_link` varchar(60) NOT NULL,
  PRIMARY KEY (`google_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gorp_migrations`
--

DROP TABLE IF EXISTS `gorp_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gorp_migrations` (
  `id` varchar(255) NOT NULL,
  `applied_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `hits`
--

DROP TABLE IF EXISTS `hits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `hits` (
  `hit_id` int(11) NOT NULL AUTO_INCREMENT,
  `w_id` int(11) DEFAULT '0',
  `hit_date` int(11) DEFAULT '0',
  PRIMARY KEY (`hit_id`),
  KEY `w_id` (`w_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `impressions`
--

DROP TABLE IF EXISTS `impressions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `impressions` (
  `imp_id` int(11) NOT NULL AUTO_INCREMENT,
  `creatives_location_id` int(10) unsigned DEFAULT NULL,
  `publisher_page_id` int(10) unsigned DEFAULT NULL,
  `seat_id` int(10) unsigned DEFAULT NULL,
  `cp_id` int(11) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  `wp_id` int(11) DEFAULT '0',
  `app_id` int(11) DEFAULT '0',
  `ad_id` int(11) DEFAULT '0',
  `cop_id` int(11) DEFAULT '0',
  `ca_id` int(11) DEFAULT '0',
  `imp_ipaddress` varchar(50) DEFAULT NULL,
  `imp_referaddress` text,
  `imp_parenturl` text,
  `imp_url` text,
  `imp_winnerbid` int(11) DEFAULT '0',
  `imp_status` tinyint(1) DEFAULT '0',
  `imp_cookie` tinyint(1) DEFAULT '1',
  `imp_alexa` tinyint(1) DEFAULT '0',
  `imp_conv` tinyint(1) DEFAULT '0',
  `imp_flash` tinyint(1) DEFAULT '1',
  `imp_time` int(11) DEFAULT '0',
  `imp_date` int(8) DEFAULT '0',
  PRIMARY KEY (`imp_id`),
  KEY `imp_date` (`imp_date`),
  KEY `w_id` (`w_id`,`imp_date`),
  KEY `ca_id` (`ca_id`,`imp_date`),
  KEY `app_id` (`app_id`,`imp_date`),
  KEY `seat_id` (`seat_id`),
  KEY `publisher_page_id` (`publisher_page_id`),
  KEY `creatives_location_id` (`creatives_location_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `impressions-cells`
--

DROP TABLE IF EXISTS `impressions-cells`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `impressions-cells` (
  `imp_id` int(11) NOT NULL AUTO_INCREMENT,
  `creatives_location_id` int(10) unsigned DEFAULT NULL,
  `publisher_page_id` int(10) unsigned DEFAULT NULL,
  `seat_id` int(10) unsigned DEFAULT NULL,
  `cp_id` int(11) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  `app_id` int(11) DEFAULT '0',
  `wp_id` int(11) DEFAULT '0',
  `ad_id` int(11) DEFAULT '0',
  `cop_id` int(11) DEFAULT '0',
  `ca_id` int(11) DEFAULT '0',
  `slot_id` int(11) DEFAULT '0',
  `sla_id` int(11) DEFAULT '0',
  `cell_id` int(11) DEFAULT '0',
  `imp_ipaddress` varchar(50) DEFAULT NULL,
  `imp_referaddress` text,
  `imp_parenturl` text,
  `imp_url` text,
  `imp_winnerbid` int(11) DEFAULT '0',
  `imp_status` tinyint(1) DEFAULT '0',
  `imp_cookie` tinyint(1) DEFAULT '1',
  `imp_alexa` tinyint(1) DEFAULT '0',
  `imp_flash` tinyint(1) DEFAULT '1',
  `imp_time` int(11) DEFAULT '0',
  `imp_date` int(8) DEFAULT '0',
  PRIMARY KEY (`imp_id`),
  KEY `imp_date` (`imp_date`),
  KEY `w_id` (`w_id`),
  KEY `ca_id` (`ca_id`),
  KEY `slot_id` (`slot_id`),
  KEY `sla_id` (`sla_id`),
  KEY `app_id` (`app_id`),
  KEY `seat_id` (`seat_id`),
  KEY `publisher_page_id` (`publisher_page_id`),
  KEY `creatives_location_id` (`creatives_location_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `interests`
--

DROP TABLE IF EXISTS `interests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `interests` (
  `in_id` int(11) NOT NULL AUTO_INCREMENT,
  `in_parent_id` int(11) DEFAULT '0',
  `in_gender` tinyint(1) DEFAULT '0',
  `in_age` tinyint(1) DEFAULT '0',
  `in_name` int(11) DEFAULT NULL,
  PRIMARY KEY (`in_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `invoices`
--

DROP TABLE IF EXISTS `invoices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invoices` (
  `in_id` int(11) NOT NULL AUTO_INCREMENT,
  `in_accept` tinyint(2) NOT NULL DEFAULT '0',
  `in_serial` int(11) NOT NULL DEFAULT '0',
  `in_date` int(11) NOT NULL DEFAULT '0',
  `u_id` int(11) NOT NULL,
  `in_price` int(11) NOT NULL,
  `in_title` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `in_user_register_id` int(11) NOT NULL,
  `in_sale_condition` int(1) DEFAULT '0',
  PRIMARY KEY (`in_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `invoices_details`
--

DROP TABLE IF EXISTS `invoices_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invoices_details` (
  `ind_id` int(11) NOT NULL AUTO_INCREMENT,
  `in_id` int(11) NOT NULL,
  `ind_description` varchar(500) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ind_timing` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ind_price` int(11) NOT NULL,
  `ind_price_off` int(11) NOT NULL,
  PRIMARY KEY (`ind_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ip2location`
--

DROP TABLE IF EXISTS `ip2location`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ip2location` (
  `ip_from` int(10) unsigned DEFAULT NULL,
  `ip_to` int(10) unsigned DEFAULT NULL,
  `country_code` char(2) COLLATE utf8_bin DEFAULT NULL,
  `country_name` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `region_name` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `city_name` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  KEY `CC` (`country_code`,`ip_to`),
  KEY `ip_from` (`ip_from`,`ip_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=COMPACT;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ip2location3`
--

DROP TABLE IF EXISTS `ip2location3`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ip2location3` (
  `ip_from` int(10) unsigned DEFAULT NULL,
  `ip_to` int(10) unsigned DEFAULT NULL,
  `country_code` char(2) DEFAULT NULL,
  `country_name` varchar(64) DEFAULT NULL,
  `region_name` varchar(128) DEFAULT NULL,
  `city_name` varchar(128) DEFAULT NULL,
  `isp` varchar(128) DEFAULT NULL,
  KEY `CC` (`country_code`,`ip_to`),
  KEY `ip_from` (`ip_from`,`ip_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ip2location_ir`
--

DROP TABLE IF EXISTS `ip2location_ir`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ip2location_ir` (
  `ip_from` int(10) unsigned DEFAULT NULL,
  `ip_to` int(10) unsigned DEFAULT NULL,
  `country_code` char(2) DEFAULT NULL,
  `country_name` varchar(64) DEFAULT NULL,
  `region_name` varchar(128) DEFAULT NULL,
  `city_name` varchar(128) DEFAULT NULL,
  KEY `CC` (`country_code`,`ip_to`),
  KEY `ip_from` (`ip_from`,`ip_to`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ip2location_ir_old`
--

DROP TABLE IF EXISTS `ip2location_ir_old`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ip2location_ir_old` (
  `ip_from` int(10) unsigned DEFAULT NULL,
  `ip_to` int(10) unsigned DEFAULT NULL,
  `country_code` char(2) COLLATE utf8_bin DEFAULT NULL,
  `country_name` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `region_name` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `city_name` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  KEY `CC` (`country_code`,`ip_to`),
  KEY `ip_from` (`ip_from`,`ip_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=COMPACT;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ip2location_range`
--

DROP TABLE IF EXISTS `ip2location_range`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ip2location_range` (
  `id` int(11) NOT NULL,
  `range_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ip2location_test`
--

DROP TABLE IF EXISTS `ip2location_test`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ip2location_test` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip_from` int(10) unsigned DEFAULT NULL,
  `ip_to` int(10) unsigned DEFAULT NULL,
  `country_code` char(2) COLLATE utf8_bin DEFAULT NULL,
  `country_name` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `region_name` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `city_name` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `CC` (`country_code`,`ip_to`),
  KEY `ip_from` (`ip_from`,`ip_to`)
) ENGINE=InnoDB AUTO_INCREMENT=20757 DEFAULT CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=COMPACT;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `keywords`
--

DROP TABLE IF EXISTS `keywords`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `keywords` (
  `k_id` int(11) NOT NULL AUTO_INCREMENT,
  `k_string` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `k_string_md5` varchar(64) CHARACTER SET utf8mb4 DEFAULT NULL,
  `k_confirm` tinyint(1) DEFAULT '0',
  `k_count` int(11) DEFAULT '0',
  PRIMARY KEY (`k_id`),
  UNIQUE KEY `k_string_md5` (`k_string_md5`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `keywords_interests`
--

DROP TABLE IF EXISTS `keywords_interests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `keywords_interests` (
  `ki_id` int(11) NOT NULL AUTO_INCREMENT,
  `k_id` int(11) DEFAULT '0',
  `in_id` int(11) DEFAULT '0',
  PRIMARY KEY (`ki_id`),
  KEY `k_id` (`k_id`),
  KEY `in_id` (`in_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `keywords_webpages`
--

DROP TABLE IF EXISTS `keywords_webpages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `keywords_webpages` (
  `kwp_id` int(11) NOT NULL,
  `k_id` int(11) DEFAULT '0',
  `wp_id` int(11) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  PRIMARY KEY (`kwp_id`),
  KEY `wp_id` (`wp_id`),
  KEY `k_id_2` (`k_id`),
  KEY `k_id` (`k_id`,`wp_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `list_asset`
--

DROP TABLE IF EXISTS `list_asset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `list_asset` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `hash` varchar(40) NOT NULL,
  `url` text NOT NULL,
  `img` text NOT NULL,
  `title` varchar(255) NOT NULL,
  `price` int(11) DEFAULT NULL,
  `discount` int(11) DEFAULT NULL,
  `sku` varchar(255) DEFAULT NULL,
  `is_available` tinyint(1) DEFAULT '1',
  `cat` varchar(255) DEFAULT NULL,
  `brand` varchar(255) DEFAULT NULL,
  `mime` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `hash` (`hash`),
  KEY `is_available` (`is_available`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `list_browser`
--

DROP TABLE IF EXISTS `list_browser`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `list_browser` (
  `browser_id` int(11) NOT NULL AUTO_INCREMENT,
  `browser_value` varchar(100) DEFAULT NULL,
  `browser_name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`browser_id`),
  UNIQUE KEY `bor_value` (`browser_value`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `list_city`
--

DROP TABLE IF EXISTS `list_city`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `list_city` (
  `location_id` int(11) NOT NULL AUTO_INCREMENT,
  `location_name` text CHARACTER SET utf8mb4,
  `location_name_persian` text CHARACTER SET utf8mb4,
  `location_master` mediumint(6) DEFAULT '0',
  `location_select` tinyint(1) DEFAULT '0',
  `location_code` int(11) NOT NULL,
  `location_country` varchar(3) CHARACTER SET utf8mb4 DEFAULT NULL,
  `location_region` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`location_id`)
) ENGINE=InnoDB AUTO_INCREMENT=182 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `list_locations`
--

DROP TABLE IF EXISTS `list_locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `list_locations` (
  `location_id` int(11) NOT NULL AUTO_INCREMENT,
  `location_name` text CHARACTER SET utf8mb4,
  `location_name_persian` text CHARACTER SET utf8mb4,
  `location_master` mediumint(6) DEFAULT '0',
  `location_select` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`location_id`),
  KEY `location_name` (`location_name`(50))
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `list_platform`
--

DROP TABLE IF EXISTS `list_platform`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `list_platform` (
  `platform_id` int(11) NOT NULL AUTO_INCREMENT,
  `platform_network` tinyint(1) DEFAULT '0',
  `platform_value` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `platform_name` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  PRIMARY KEY (`platform_id`),
  UNIQUE KEY `osl_value` (`platform_value`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `list_region`
--

DROP TABLE IF EXISTS `list_region`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `list_region` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `region_name` varchar(50) CHARACTER SET utf8mb4 DEFAULT NULL,
  `region_name_persian` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `region_code` int(2) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `list_region2`
--

DROP TABLE IF EXISTS `list_region2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `list_region2` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `region_name` varchar(50) DEFAULT NULL,
  `region_name_persian` varchar(100) DEFAULT NULL,
  `region_code` int(2) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mastertabletest`
--

DROP TABLE IF EXISTS `mastertabletest`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mastertabletest` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mastertabletesto`
--

DROP TABLE IF EXISTS `mastertabletesto`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mastertabletesto` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mastertabletests`
--

DROP TABLE IF EXISTS `mastertabletests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mastertabletests` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mastertabletestt`
--

DROP TABLE IF EXISTS `mastertabletestt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mastertabletestt` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `neighborhoods`
--

DROP TABLE IF EXISTS `neighborhoods`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `neighborhoods` (
  `id` mediumint(7) unsigned NOT NULL AUTO_INCREMENT,
  `title` text CHARACTER SET utf8mb4,
  `cellsgroup` text CHARACTER SET utf8mb4,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `neighborhoods_old`
--

DROP TABLE IF EXISTS `neighborhoods_old`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `neighborhoods_old` (
  `id` mediumint(7) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `cellsgroup` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `password_resets`
--

DROP TABLE IF EXISTS `password_resets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `password_resets` (
  `u_email` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `token` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `password_resets_email_index` (`u_email`),
  KEY `password_resets_token_index` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `payment_transaction`
--

DROP TABLE IF EXISTS `payment_transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_transaction` (
  `pt_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_id` int(11) DEFAULT '0',
  `pt_amount` int(11) DEFAULT '0',
  `pt_type` tinyint(2) DEFAULT '0',
  `pt_gate` varchar(100) DEFAULT NULL,
  `pt_status` tinyint(2) DEFAULT '0',
  `pt_authority` varchar(255) DEFAULT NULL,
  `pt_refid` varchar(255) DEFAULT NULL,
  `pt_time` int(11) DEFAULT '0',
  `pt_date` int(11) DEFAULT '0',
  `pt_flag` tinyint(2) DEFAULT '0',
  `pt_vat` int(11) NOT NULL DEFAULT '0',
  `factor_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`pt_id`),
  KEY `u_id` (`u_id`),
  KEY `pt_authority` (`pt_authority`),
  KEY `payment_transaction_billing_factor_id_fk` (`factor_id`),
  CONSTRAINT `payment_transaction_billing_factor_id_fk` FOREIGN KEY (`factor_id`) REFERENCES `billing_factor` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12008 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `permission_role`
--

DROP TABLE IF EXISTS `permission_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permission_role` (
  `permission_id` int(10) unsigned NOT NULL,
  `role_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`permission_id`,`role_id`),
  KEY `permission_role_role_id_foreign` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `permissions`
--

DROP TABLE IF EXISTS `permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permissions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `label` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `access` varchar(100) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'list',
  `action` varchar(100) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'own',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=211 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `provinces`
--

DROP TABLE IF EXISTS `provinces`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `provinces` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `country_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `provinces_name_country_id_unidex` (`name`,`country_id`),
  KEY `provinces_countries_id_fk` (`country_id`),
  CONSTRAINT `provinces_countries_id_fk` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `publisher_pages`
--

DROP TABLE IF EXISTS `publisher_pages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `publisher_pages` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `publisher_id` int(10) unsigned NOT NULL,
  `publisher_domain` varchar(200) NOT NULL,
  `kind` enum('web','app') NOT NULL,
  `url` text NOT NULL,
  `url_key` varchar(255) NOT NULL,
  `active_days` int(10) unsigned NOT NULL DEFAULT '1',
  `avg_daily_imp` int(10) unsigned NOT NULL DEFAULT '0',
  `avg_daily_clicks` int(10) unsigned NOT NULL DEFAULT '0',
  `today_imp` int(10) unsigned NOT NULL DEFAULT '0',
  `today_clicks` int(10) unsigned NOT NULL DEFAULT '0',
  `today_ctr` int(10) unsigned NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `pub_page_unique` (`url_key`,`publisher_id`)
) ENGINE=InnoDB AUTO_INCREMENT=238606 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `qlog`
--

DROP TABLE IF EXISTS `qlog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `qlog` (
  `q_id` int(11) NOT NULL AUTO_INCREMENT,
  `q_content` text,
  `q_time` int(11) DEFAULT '0',
  PRIMARY KEY (`q_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `retargeting_item`
--

DROP TABLE IF EXISTS `retargeting_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `retargeting_item` (
  `ri_id` int(11) NOT NULL AUTO_INCREMENT,
  `ri_key` varchar(11) DEFAULT NULL,
  `rl_id` int(11) DEFAULT '0',
  `ri_name` text,
  `ri_url` text,
  `ri_img_url` text,
  `ri_price` int(11) DEFAULT '-1',
  PRIMARY KEY (`ri_id`),
  KEY `ri_key` (`ri_key`),
  KEY `rl_id` (`rl_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `retargeting_list`
--

DROP TABLE IF EXISTS `retargeting_list`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `retargeting_list` (
  `rl_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_id` int(11) DEFAULT '0',
  `rl_key` varchar(11) DEFAULT NULL,
  `rl_domain` varchar(255) NOT NULL,
  `rl_name` varchar(255) DEFAULT NULL,
  `rl_date` int(11) DEFAULT '0',
  PRIMARY KEY (`rl_id`),
  KEY `rl_key` (`rl_key`),
  KEY `u_id` (`u_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `role_user`
--

DROP TABLE IF EXISTS `role_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role_user` (
  `role_id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`role_id`,`user_id`),
  KEY `role_user_user_id_foreign` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `label` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `childes` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `seats`
--

DROP TABLE IF EXISTS `seats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `seats` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `slot_id` varchar(255) NOT NULL,
  `supplier_name` varchar(200) NOT NULL,
  `publisher_id` int(10) unsigned NOT NULL,
  `publisher_domain` varchar(200) NOT NULL,
  `creative_size` int(10) unsigned NOT NULL,
  `kind` enum('web','app') NOT NULL,
  `active_days` int(10) unsigned NOT NULL DEFAULT '1',
  `avg_daily_imp` int(10) unsigned NOT NULL DEFAULT '0',
  `avg_daily_clicks` int(10) unsigned NOT NULL DEFAULT '0',
  `today_imp` int(10) unsigned NOT NULL DEFAULT '0',
  `today_clicks` int(10) unsigned NOT NULL DEFAULT '0',
  `today_ctr` int(10) unsigned NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slot_per_sup_unique` (`slot_id`,`supplier_name`,`publisher_domain`,`creative_size`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=1054681 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `segments`
--

DROP TABLE IF EXISTS `segments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `segments` (
  `seg_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_id` int(11) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  `seg_pattern` varchar(255) DEFAULT NULL,
  `seg_type` tinyint(1) DEFAULT '0',
  `seg_name` varchar(255) DEFAULT NULL,
  `seg_isconv` tinyint(1) DEFAULT '0',
  `seg_convvalue` int(11) DEFAULT '0',
  `seg_conversions` int(11) DEFAULT '0',
  `seg_visitors` int(11) DEFAULT '0',
  `seg_pageview` int(11) DEFAULT '0',
  `seg_lastupdate` int(11) DEFAULT '0',
  PRIMARY KEY (`seg_id`),
  KEY `u_id` (`u_id`)
) ENGINE=InnoDB AUTO_INCREMENT=846 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `slot_pin`
--

DROP TABLE IF EXISTS `slot_pin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `slot_pin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `slot_id` int(11) NOT NULL,
  `chance` int(11) NOT NULL,
  `ad_id` int(11) NOT NULL,
  `bid` int(11) NOT NULL,
  `start` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `end` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `direct` tinyint(1) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  UNIQUE KEY `slot_pin_slot_id_uindex` (`slot_id`),
  KEY `slot_pin_ads_ad_id_fk` (`ad_id`),
  CONSTRAINT `slot_pin_ads_ad_id_fk` FOREIGN KEY (`ad_id`) REFERENCES `ads` (`ad_id`),
  CONSTRAINT `slot_pin_slots_slot_id_fk` FOREIGN KEY (`slot_id`) REFERENCES `slots` (`slot_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `slots`
--

DROP TABLE IF EXISTS `slots`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `slots` (
  `slot_id` int(11) NOT NULL AUTO_INCREMENT,
  `slot_pubilc_id` bigint(20) DEFAULT '0',
  `slot_name` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `slot_size` tinyint(2) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  `app_id` int(11) DEFAULT '0',
  `slot_avg_daily_imps` int(11) DEFAULT '0',
  `slot_avg_daily_clicks` int(11) DEFAULT '0',
  `slot_floor_cpm` int(11) DEFAULT '0',
  `slot_total_monthly_cost` int(11) DEFAULT '0',
  `slot_lastupdate` int(11) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`slot_id`),
  UNIQUE KEY `web_slots__index` (`w_id`,`slot_pubilc_id`,`app_id`),
  KEY `app_slots__index` (`slot_pubilc_id`,`app_id`),
  KEY `slots_slot_pubilc_id_app_id_index` (`slot_pubilc_id`,`app_id`)
) ENGINE=InnoDB AUTO_INCREMENT=79231502 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `slots_ads`
--

DROP TABLE IF EXISTS `slots_ads`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `slots_ads` (
  `sla_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `slot_id` int(11) DEFAULT '0',
  `ad_id` int(11) DEFAULT '0',
  `sla_imps` int(11) DEFAULT '0',
  `sla_clicks` int(11) DEFAULT '0',
  `sla_ctr` float DEFAULT '0',
  `sla_conv` int(11) DEFAULT '0',
  `sla_conv_rate` float DEFAULT '0',
  `sla_cpa` int(11) DEFAULT '0',
  `sla_cpm` int(11) DEFAULT '0',
  `sla_spend` int(11) DEFAULT '0',
  `sla_lastupdate` int(11) DEFAULT '0',
  PRIMARY KEY (`sla_id`),
  UNIQUE KEY `uniq_slot_ad` (`slot_id`,`ad_id`),
  KEY `slot_id` (`slot_id`,`ad_id`,`sla_ctr`)
) ENGINE=InnoDB AUTO_INCREMENT=9826778494 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;
