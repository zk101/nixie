--
-- Database: `nixie_telemetry`
--

CREATE DATABASE IF NOT EXISTS `nixie_telemetry` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `nixie_telemetry`;

--
-- Table structure for table `telemetry`
--

DROP TABLE IF EXISTS `telemetry`;

CREATE TABLE `telemetry` (
  `telemetry_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `telemetry_user_string` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `telemetry_client_string` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `telemetry_client_version` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `telemetry_date` bigint(20) unsigned NOT NULL DEFAULT 0,
  `telemetry_data` text COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `telemetry_dt_created` datetime NOT NULL DEFAULT current_timestamp(),
  `telemetry_dt_modified` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  `telemetry_status` tinyint(4) NOT NULL DEFAULT 0,
  PRIMARY KEY (`telemetry_id`),
  KEY `telemetry_user_string` (`telemetry_user_string`),
  KEY `telemetry_client_string` (`telemetry_client_string`),
  KEY `telemetry_client_version` (`telemetry_client_version`),
  KEY `telemetry_date` (`telemetry_date`),
  KEY `telemetry_status` (`telemetry_status`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- EOF
