-- MySQL dump 10.16  Distrib 10.1.20-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: localhost
-- ------------------------------------------------------
-- Server version	10.1.20-MariaDB-1~xenial

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
  `ad_mainText` varchar(128) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_defineText` varchar(128) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_textColor` varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_target` varchar(30) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ad_attribute` text CHARACTER SET utf8mb4,
  `ad_hash_attribute` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`ad_id`),
  KEY `u_id` (`u_id`),
  KEY `ad_size` (`ad_size`,`ad_status`),
  KEY `ad_hash_attribute` (`ad_hash_attribute`(191))
) ENGINE=InnoDB AUTO_INCREMENT=125041 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=latin1;
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
  `app_supplier` varchar(32) NOT NULL DEFAULT 'clickyab',
  `en_app_name` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `app_package` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `am_id` tinyint(1) DEFAULT '0' COMMENT 'Market Id',
  `app_minbid` int(11) DEFAULT '500',
  `app_floor_cpm` int(11) DEFAULT '250',
  `app_div` float DEFAULT '1.7',
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
  PRIMARY KEY (`app_id`),
  KEY `u_id` (`u_id`),
  KEY `app_token` (`app_token`),
  KEY `apps_app_token_app_status_index` (`app_token`,`app_status`)
) ENGINE=InnoDB AUTO_INCREMENT=88145 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=221090 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=3226413 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=519361 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=129258 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=1371648 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=135343 DEFAULT CHARSET=latin1;
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
  PRIMARY KEY (`cp_id`),
  KEY `cp_lock` (`cp_lock`),
  KEY `u_id_idx` (`u_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=23836 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=123642 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=latin1;
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
  `c_status` tinyint(2) DEFAULT '0',
  `c_ip` varchar(20) DEFAULT NULL,
  `c_referaddress` varchar(255) DEFAULT NULL,
  `c_parenturl` varchar(255) DEFAULT NULL,
  `c_fast` int(12) DEFAULT '0',
  `c_os` tinyint(1) DEFAULT '0',
  `c_time` int(11) DEFAULT '0',
  `c_date` int(11) DEFAULT '0',
  `ad_size` int(11) NOT NULL DEFAULT '0',
  `reserved_hash` VARCHAR( 63 ) DEFAULT NULL,
  PRIMARY KEY (`c_id`),
  KEY `c_date` (`c_date`),
  KEY `w_id` (`w_id`,`c_date`),
  KEY `ca_id` (`ca_id`,`c_status`,`c_date`),
  KEY `app_id` (`app_id`,`c_date`),
  KEY `sla_id` (`sla_id`,`c_status`,`c_date`),
  KEY `reserved_hash` (`reserved_hash`)
) ENGINE=InnoDB AUTO_INCREMENT=114677544 DEFAULT CHARSET=latin1;
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
  `cop_id` int(11) DEFAULT '0',
  `imp_id` int(11) DEFAULT '0',
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
  UNIQUE KEY `unique_imp` (`cp_id`,`ad_id`,`imp_id`),
  KEY `c_date` (`c_date`),
  KEY `app_id` (`app_id`),
  KEY `slot_id` (`slot_id`),
  KEY `w_id` (`w_id`),
  KEY `sla_id_single` (`sla_id`)
) ENGINE=InnoDB AUTO_INCREMENT=816337 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=827891 DEFAULT CHARSET=latin1;
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `ctr_stat_pub_id_pub_type_date_uindex` (`pub_id`,`pub_type`,`date`)
) ENGINE=InnoDB AUTO_INCREMENT=122356 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=152208326 DEFAULT CHARSET=utf8;
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
-- Table structure for table `impressions-cells`
--

DROP TABLE IF EXISTS `impressions-cells`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `impressions-cells` (
  `imp_id` int(11) NOT NULL AUTO_INCREMENT,
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
  KEY `app_id` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `impressions20171211`
--

DROP TABLE IF EXISTS `impressions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `impressions` (
  `imp_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `w_id` int(11) DEFAULT '0',
  `app_id` int(11) DEFAULT '0',
  `wp_id` int(11) DEFAULT '0',
  `ad_id` int(11) DEFAULT '0',
  `ad_size` int(11) DEFAULT '0',
  `cop_id` int(11) DEFAULT '0',
  `ca_id` int(11) DEFAULT '0',
  `slot_id` int(11) DEFAULT '0',
  `sla_id` bigint(20) DEFAULT '0',
  `cell_id` int(11) DEFAULT '0',
  `hood_id` int(11) DEFAULT '0',
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
  `reserved_hash` VARCHAR( 63 ) DEFAULT NULL,
  PRIMARY KEY (`imp_id`),
  KEY `imp_date` (`imp_date`),
  KEY `w_id` (`w_id`),
  KEY `ca_id` (`ca_id`),
  KEY `slot_id` (`slot_id`),
  KEY `sla_id` (`sla_id`),
  KEY `app_id` (`app_id`),
  KEY `wb_adsize` (`ad_size`,`w_id`),
  KEY `app_adsize` (`ad_size`,`app_id`),
  KEY `reserved_hash` (`reserved_hash`)
) ENGINE=InnoDB AUTO_INCREMENT=1066346484 DEFAULT CHARSET=latin1;
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
  PRIMARY KEY (`pt_id`),
  KEY `u_id` (`u_id`),
  KEY `pt_authority` (`pt_authority`)
) ENGINE=InnoDB AUTO_INCREMENT=11327 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=3896748 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=5279100674 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `slots_old`
--

DROP TABLE IF EXISTS `slots_old`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `slots_old` (
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
  KEY `app_slots__index` (`slot_pubilc_id`,`app_id`),
  KEY `web_slots__index` (`w_id`,`slot_pubilc_id`),
  KEY `slots_slot_pubilc_id_app_id_index` (`slot_pubilc_id`,`app_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3754823 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statistics_ads`
--

DROP TABLE IF EXISTS `statistics_ads`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statistics_ads` (
  `sa_id` int(11) NOT NULL AUTO_INCREMENT,
  `ad_id` int(11) DEFAULT '0',
  `cp_id` int(11) DEFAULT '0',
  `ca_id` int(11) DEFAULT '0',
  `sa_clicks` int(11) DEFAULT '0',
  `sa_imps` int(11) DEFAULT '0',
  `sa_ctr` float DEFAULT '0',
  `sa_conv` int(11) DEFAULT '0',
  `sa_conv_rate` float DEFAULT '0',
  `sa_cpa` int(11) DEFAULT '0',
  `sa_spend` int(11) DEFAULT '0',
  `sa_day` int(11) DEFAULT '0',
  `sa_done` tinyint(1) DEFAULT '0',
  `sa_trueview` int(11) DEFAULT NULL,
  PRIMARY KEY (`sa_id`),
  KEY `cp_id` (`ca_id`),
  KEY `ca_id` (`ca_id`,`sa_day`),
  KEY `sa_day` (`sa_day`)
) ENGINE=InnoDB AUTO_INCREMENT=1163029 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statistics_alexa`
--

DROP TABLE IF EXISTS `statistics_alexa`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statistics_alexa` (
  `sal_id` int(11) NOT NULL AUTO_INCREMENT,
  `sal_ip` varchar(100) DEFAULT NULL,
  `sal_ref` varchar(255) DEFAULT NULL,
  `sal_user_agent` varchar(255) DEFAULT NULL,
  `sal_show` tinyint(1) DEFAULT '0',
  `sal_time` int(11) DEFAULT '0',
  `sal_day` int(11) DEFAULT '0',
  PRIMARY KEY (`sal_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statistics_apps`
--

DROP TABLE IF EXISTS `statistics_apps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statistics_apps` (
  `sa_id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) DEFAULT '0',
  `sa_clicks` int(11) DEFAULT '0',
  `sa_allclicks` int(11) DEFAULT '0',
  `sa_imps` int(11) DEFAULT '0',
  `sa_ctr` float DEFAULT '0',
  `sa_conv` int(11) DEFAULT '0',
  `sa_conv_rate` float DEFAULT '0',
  `sa_cpa` int(11) DEFAULT '0',
  `sa_spend` int(11) DEFAULT '0',
  `sa_day` int(11) DEFAULT '0',
  `sa_done` tinyint(4) DEFAULT '0',
  PRIMARY KEY (`sa_id`),
  UNIQUE KEY `app_id` (`app_id`,`sa_day`),
  KEY `w_id` (`app_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3519320 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statistics_billing`
--

DROP TABLE IF EXISTS `statistics_billing`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statistics_billing` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `network_debt` bigint(20) NOT NULL DEFAULT '0',
  `wallet` bigint(20) NOT NULL DEFAULT '0',
  `rejected_withdrawal` bigint(20) NOT NULL DEFAULT '0',
  `pending_withdrawal` bigint(20) NOT NULL DEFAULT '0',
  `ready_withdrawal` bigint(20) NOT NULL DEFAULT '0',
  `publishers_income` bigint(20) NOT NULL DEFAULT '0',
  `cpc_income` bigint(20) NOT NULL DEFAULT '0',
  `prepayment_income` bigint(20) NOT NULL DEFAULT '0',
  `publishers_deposit` bigint(20) NOT NULL DEFAULT '0',
  `vertex_number` int(11) NOT NULL DEFAULT '0',
  `vertex_withdrawal` text NOT NULL,
  `action_date` int(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `statistics_billing_action_date_uindex` (`action_date`)
) ENGINE=InnoDB AUTO_INCREMENT=138 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statistics_campaigns`
--

DROP TABLE IF EXISTS `statistics_campaigns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statistics_campaigns` (
  `sc_id` int(11) NOT NULL AUTO_INCREMENT,
  `cp_id` int(11) DEFAULT '0',
  `sc_clicks` int(11) DEFAULT '0',
  `sc_imps` int(11) DEFAULT '0',
  `sc_ctr` float DEFAULT '0',
  `sc_conv` int(11) DEFAULT '0',
  `sc_conv_rate` float DEFAULT '0',
  `sc_cpa` int(11) DEFAULT '0',
  `sc_spend` int(11) DEFAULT '0',
  `sc_day` int(11) DEFAULT '0',
  `sc_done` tinyint(4) DEFAULT '0',
  `sc_trueview` int(11) DEFAULT NULL,
  PRIMARY KEY (`sc_id`),
  KEY `cp_id` (`cp_id`)
) ENGINE=InnoDB AUTO_INCREMENT=210566 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statistics_segments`
--

DROP TABLE IF EXISTS `statistics_segments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statistics_segments` (
  `ss_id` int(11) NOT NULL AUTO_INCREMENT,
  `seg_id` int(11) DEFAULT '0',
  `cs_id` int(11) DEFAULT '0',
  `ss_pageviews` int(11) DEFAULT '0',
  `ss_newvisitors` int(11) DEFAULT '0',
  `ss_conversions` int(11) DEFAULT '0',
  `ss_revenue` int(11) DEFAULT '0',
  `ss_day` int(11) DEFAULT '0',
  `ss_done` tinyint(4) DEFAULT '0',
  PRIMARY KEY (`ss_id`),
  KEY `seg_id` (`seg_id`)
) ENGINE=InnoDB AUTO_INCREMENT=389 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statistics_slot`
--

DROP TABLE IF EXISTS `statistics_slot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statistics_slot` (
  `ssl_id` int(11) NOT NULL AUTO_INCREMENT,
  `ad_id` int(11) DEFAULT '0',
  `slot_id` int(11) DEFAULT '0',
  `sla_id` int(11) DEFAULT '0',
  `ssl_clicks` int(11) DEFAULT '0',
  `ssl_imps` int(11) DEFAULT '0',
  `ssl_ctr` float DEFAULT '0',
  `ssl_conv` int(11) DEFAULT '0',
  `ssl_conv_rate` float DEFAULT '0',
  `ssl_cpa` int(11) DEFAULT '0',
  `ssl_spend` int(11) DEFAULT '0',
  `ssl_day` int(11) DEFAULT '0',
  `ssl_done` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`ssl_id`),
  KEY `slot_id` (`slot_id`),
  KEY `ssl_day` (`ssl_day`),
  KEY `sla_id` (`sla_id`,`ssl_day`),
  KEY `slot_id_day` (`slot_id`,`ssl_day`)
) ENGINE=InnoDB AUTO_INCREMENT=116459483 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statistics_websites`
--

DROP TABLE IF EXISTS `statistics_websites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statistics_websites` (
  `sw_id` int(11) NOT NULL AUTO_INCREMENT,
  `w_id` int(11) DEFAULT '0',
  `sw_clicks` int(11) DEFAULT '0',
  `sw_allclicks` int(11) DEFAULT '0',
  `sw_imps` int(11) DEFAULT '0',
  `sw_ctr` float DEFAULT '0',
  `sw_conv` int(11) DEFAULT '0',
  `sw_conv_rate` float DEFAULT '0',
  `sw_cpa` int(11) DEFAULT '0',
  `sw_spend` int(11) DEFAULT '0',
  `sw_day` int(11) DEFAULT '0',
  `sw_done` tinyint(4) DEFAULT '0',
  `sw_trueview` int(11) DEFAULT NULL,
  PRIMARY KEY (`sw_id`),
  UNIQUE KEY `w_id_2` (`w_id`,`sw_day`),
  KEY `w_id` (`w_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2079092 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `suppliers`
--

DROP TABLE IF EXISTS `suppliers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `suppliers` (
  `name` varchar(32) NOT NULL,
  `token` varchar(100) NOT NULL,
  `user_id` int(11) DEFAULT NULL,
  `default_floor` int(11) NOT NULL,
  `default_soft_floor` int(11) NOT NULL,
  `default_min_bid` int(11) NOT NULL,
  `bid_type` varchar(16) NOT NULL DEFAULT 'cpc',
  `default_ctr` float NOT NULL DEFAULT '0.1',
  `tiny_mark` int(1) NOT NULL DEFAULT '1',
  `show_domain` varchar(60) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `rate` float NOT NULL DEFAULT '1',
  `tiny_logo` varchar(200) NOT NULL,
  `tiny_url` varchar(200) NOT NULL,

  PRIMARY KEY (`name`),
  UNIQUE KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `test_table`
--

DROP TABLE IF EXISTS `test_table`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `test_table` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `test_table_name_uindex` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tickets`
--

DROP TABLE IF EXISTS `tickets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tickets` (
  `ti_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `ti_slug` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_id` int(11) NOT NULL,
  `ti_title` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ti_content` longtext NOT NULL,
  `ti_response` longtext NOT NULL,
  `ti_status` tinyint(4) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`ti_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trueview`
--

DROP TABLE IF EXISTS `trueview`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trueview` (
  `tv_click_id` int(11) NOT NULL,
  PRIMARY KEY (`tv_click_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `unsubscribe`
--

DROP TABLE IF EXISTS `unsubscribe`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `unsubscribe` (
  `email` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `u_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_email` varchar(250) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_password` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_access_key` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_profile_type` tinyint(1) DEFAULT '0',
  `u_role` tinyint(1) DEFAULT '0' COMMENT 'general = 0 , 1 = super admin , 2 => publisher, 3 => accounting , 4 => sells , 5 => support',
  `u_firstname` varchar(64) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_lastname` varchar(64) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_fullname` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_avatar` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_mobile` varchar(20) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_phone` varchar(20) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_postcode` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_address` text CHARACTER SET utf8mb4,
  `u_website` varchar(50) DEFAULT NULL,
  `u_company_name` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_melli_code` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_account_number` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_card_number` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_sheba_number` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_account_holder` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_bank_name` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_province` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_city` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_economic_number` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_register_number` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_mobile_confirm` tinyint(1) DEFAULT '0',
  `u_email_confirm` tinyint(1) DEFAULT '0',
  `u_mobile_confirm_string` varchar(12) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_email_confirm_string` varchar(12) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_balance` int(11) DEFAULT '0',
  `u_today_spend` int(11) DEFAULT '0',
  `u_close` tinyint(1) DEFAULT '0',
  `u_unsubscribe` tinyint(1) DEFAULT '0',
  `u_nlsend` tinyint(1) DEFAULT '0',
  `u_ref` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_ip` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_date` date DEFAULT NULL,
  `u_time` int(11) DEFAULT '0',
  `u_google_id` varchar(25) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_gender` varchar(4) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_longitude` varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_latitude` varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_owner` varchar(128) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_originatinglead` varchar(128) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_customer_code` int(12) DEFAULT NULL,
  `u_gid` varchar(64) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_leadstatus` int(1) NOT NULL DEFAULT '0',
  `u_parent` int(11) DEFAULT '0',
  `u_crm` tinyint(1) DEFAULT '0',
  `remember_token` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `u_added_value` varchar(255) DEFAULT NULL,
  `u_national_card` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `bReadByCRM` tinyint(1) NOT NULL DEFAULT '0',
  `u_api_token` varchar(50) DEFAULT NULL,
  `u_accounting_number` int(11) DEFAULT NULL,
  `publisher_state` enum('official','notofficial') NOT NULL DEFAULT 'notofficial',
  `u_advisor` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  PRIMARY KEY (`u_id`),
  UNIQUE KEY `u_email` (`u_email`(191)) USING BTREE,
  KEY `u_email_2` (`u_email`(191),`u_password`(191)),
  KEY `u_access_key` (`u_access_key`(191)),
  KEY `u_api_token` (`u_api_token`)
) ENGINE=InnoDB AUTO_INCREMENT=51738 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users_log`
--

DROP TABLE IF EXISTS `users_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users_log` (
  `ul_id` int(11) NOT NULL AUTO_INCREMENT,
  `ul_email` varchar(255) CHARACTER SET latin1 NOT NULL,
  `ul_password` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `ul_access_key` varchar(200) CHARACTER SET latin1 DEFAULT NULL,
  `ul_profile_type` tinyint(1) DEFAULT '0',
  `ul_role` tinyint(1) DEFAULT '0' COMMENT 'general = 0 , 1 = super admin , 2 => publisher, 3 => accounting , 4 => sells , 5 => support',
  `ul_fullname` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `ul_avatar` varchar(200) CHARACTER SET latin1 DEFAULT NULL,
  `ul_mobile` varchar(20) CHARACTER SET latin1 DEFAULT '09',
  `ul_phone` varchar(20) CHARACTER SET latin1 DEFAULT NULL,
  `ul_postcode` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `ul_address` text CHARACTER SET latin1,
  `ul_company_name` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `ul_melli_code` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `ul_account_number` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `ul_card_number` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `ul_sheba_number` varchar(255) CHARACTER SET latin1 DEFAULT 'IR',
  `ul_account_holder` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `ul_bank_name` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `ul_province` varchar(200) DEFAULT NULL,
  `ul_city` varchar(200) DEFAULT NULL,
  `ul_economic_number` varchar(200) DEFAULT NULL,
  `ul_register_number` varchar(200) DEFAULT NULL,
  `ul_mobile_confirm` tinyint(1) DEFAULT '0',
  `ul_mobile_confirm_string` varchar(12) CHARACTER SET latin1 DEFAULT '0',
  `ul_email_confirm_string` varchar(12) CHARACTER SET latin1 DEFAULT '0',
  `ul_email_confirm` tinyint(1) DEFAULT '0',
  `ul_balance` int(11) DEFAULT '0',
  `ul_close` tinyint(1) DEFAULT '0',
  `ul_unsubscribe` tinyint(1) DEFAULT '0',
  `ul_nlsend` tinyint(1) DEFAULT '0',
  `ul_ref` varchar(200) CHARACTER SET latin1 DEFAULT NULL,
  `ul_ip` varchar(100) CHARACTER SET latin1 DEFAULT NULL,
  `ul_date` date DEFAULT NULL,
  `ul_time` int(11) DEFAULT '0',
  `ul_google_id` varchar(25) CHARACTER SET latin1 DEFAULT NULL,
  `ul_gender` varchar(4) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`ul_id`),
  UNIQUE KEY `ul_email` (`ul_email`)
) ENGINE=InnoDB AUTO_INCREMENT=19205 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `webpages`
--

DROP TABLE IF EXISTS `webpages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `webpages` (
  `wp_id` int(11) NOT NULL AUTO_INCREMENT,
  `w_id` int(11) DEFAULT '0',
  `wp_url` varchar(255) DEFAULT NULL,
  `wp_md5` varchar(64) DEFAULT NULL,
  `wp_keywords` text,
  PRIMARY KEY (`wp_id`),
  KEY `w_id` (`w_id`),
  KEY `wp_md5` (`wp_md5`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `websites`
--

DROP TABLE IF EXISTS `websites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `websites` (
  `w_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_id` int(11) DEFAULT '0',
  `w_pub_id` bigint(16) DEFAULT '0',
  `w_domain` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `w_supplier` varchar(32) NOT NULL DEFAULT 'clickyab',
  `w_name` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `w_categories` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `w_profile_type` tinyint(4) DEFAULT '0',
  `w_minbid` int(11) NOT NULL DEFAULT '2000',
  `w_floor_cpm` int(5) DEFAULT '500',
  `w_status` tinyint(4) DEFAULT '0' COMMENT '0 => pending,1 => accepted ,2 => rejected, 3 => deleted',
  `w_review` tinyint(1) DEFAULT '0' COMMENT '0 => pending,1 => review,2 => repending',
  `w_alexarank` int(11) DEFAULT '0',
  `w_div` float DEFAULT '2.5',
  `w_mobad` tinyint(1) DEFAULT '0',
  `w_nativead` tinyint(1) DEFAULT '0',
  `w_fatfinger` tinyint(1) DEFAULT '0',
  `w_publish_start` int(11) DEFAULT '0',
  `w_publish_end` int(11) DEFAULT '0',
  `w_publish_cost` int(11) DEFAULT '0',
  `w_prepayment` tinyint(1) DEFAULT '0',
  `w_today_ctr` float DEFAULT '0',
  `w_today_imps` int(11) DEFAULT '0',
  `w_today_clicks` int(11) DEFAULT '0',
  `w_date` int(11) DEFAULT '0',
  `w_notapprovedreason` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`w_id`),
  UNIQUE KEY `w_pub_id` (`w_pub_id`),
  KEY `u_id` (`u_id`),
  KEY `w_domain` (`w_domain`,`w_today_imps`) USING BTREE,
  KEY `w_pub_id_2` (`w_status`,`w_pub_id`) USING BTREE,
  KEY `w_domain_2` (`w_domain`,`w_status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=37792 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `websites_categories`
--

DROP TABLE IF EXISTS `websites_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `websites_categories` (
  `wc_id` int(11) NOT NULL AUTO_INCREMENT,
  `w_id` int(11) DEFAULT '0',
  `cat_id` int(11) DEFAULT '0',
  PRIMARY KEY (`wc_id`),
  KEY `w_id` (`w_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `websites_potential`
--

DROP TABLE IF EXISTS `websites_potential`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `websites_potential` (
  `wpt_id` int(11) NOT NULL AUTO_INCREMENT,
  `wpt_domain` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wpt_fullname` varchar(128) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wpt_status` tinyint(4) DEFAULT NULL,
  `wpt_alexa` int(11) DEFAULT '0',
  `wpt_email` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wpt_phone` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wpt_telephone` varchar(30) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wpt_partner` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wpt_categories` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wpt_clickyab_id` int(11) DEFAULT '0',
  `wpt_description` text CHARACTER SET utf8mb4,
  `wpt_telegram` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wpt_category` varchar(512) CHARACTER SET utf8mb4 DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`wpt_id`),
  KEY `wpt_domain` (`wpt_domain`(191))
) ENGINE=InnoDB AUTO_INCREMENT=6353 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `withdrawal`
--

DROP TABLE IF EXISTS `withdrawal`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `withdrawal` (
  `wd_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_id` int(11) DEFAULT '0',
  `wd_amount` int(11) DEFAULT '0',
  `wd_card_number` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wd_sheba_number` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wd_tracking_code` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wd_open` tinyint(1) DEFAULT '1',
  `wd_date` int(11) DEFAULT '0',
  `wd_date_paid` int(11) DEFAULT '0',
  `wd_reason` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wd_uid_approver` int(10) DEFAULT NULL,
  `wd_date_approvement` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wd_approve` int(1) DEFAULT '0' COMMENT '0 => Not Approved ; 1 => Approved ; ',
  `wd_is_advertiser` int(1) DEFAULT '0' COMMENT '1 = > is Advertiser ; 0 => publisher ;',
  `wd_not_approved_reasons` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wd_review_date` int(11) DEFAULT NULL,
  `wd_review_description` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `wd_can_request_review` int(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`wd_id`),
  KEY `u_id` (`u_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14618 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `withdrawal_history`
--
DROP TABLE IF EXISTS `daily_report`;
CREATE TABLE `daily_report`
(
  id INT PRIMARY KEY AUTO_INCREMENT,
  supplier VARCHAR(63) NOT NULL,
  type ENUM("web", "app") NOT NULL,
  publisher VARCHAR(63) NOT NULL,
  imps INT DEFAULT 0 NOT NULL,
  cpm INT DEFAULT 0 NOT NULL,
  clicks INT DEFAULT 0 NOT NULL,
  cpc INT DEFAULT 0 NOT NULL,
  date INT(11) NOT NULL
);


DROP TABLE IF EXISTS `withdrawal_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `withdrawal_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `withdrawal_id` int(11) NOT NULL,
  `withdrawal_amount` int(11) NOT NULL DEFAULT '0',
  `withdrawal_reason` text,
  `owner_id` int(11) NOT NULL,
  `role_id` int(11) DEFAULT NULL,
  `influencer` int(11) DEFAULT NULL,
  `impersonator` int(11) DEFAULT NULL,
  `withdrawal_status` tinyint(4) NOT NULL,
  `action` varchar(255) NOT NULL,
  `action_date` int(11) NOT NULL,
  `action_time` int(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `withdrawal_history_withdrawal_wd_id_fk` (`withdrawal_id`),
  KEY `withdrawal_history_users_u_id_fk` (`owner_id`),
  KEY `withdrawal_history_users_u_id_fk1` (`influencer`),
  KEY `withdrawal_history_users_u_id_fk2` (`impersonator`),
  KEY `withdrawal_history_roles_id_fk` (`role_id`),
  CONSTRAINT `withdrawal_history_roles_id_fk` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `withdrawal_history_users_u_id_fk` FOREIGN KEY (`owner_id`) REFERENCES `users` (`u_id`),
  CONSTRAINT `withdrawal_history_users_u_id_fk1` FOREIGN KEY (`influencer`) REFERENCES `users` (`u_id`),
  CONSTRAINT `withdrawal_history_users_u_id_fk2` FOREIGN KEY (`impersonator`) REFERENCES `users` (`u_id`),
  CONSTRAINT `withdrawal_history_withdrawal_wd_id_fk` FOREIGN KEY (`withdrawal_id`) REFERENCES `withdrawal` (`wd_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3849 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-12-12 11:39:12
