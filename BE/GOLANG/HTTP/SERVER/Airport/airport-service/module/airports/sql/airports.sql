/*关闭外键约束*/
SET FOREIGN_KEY_CHECKS=0;

/*机场表*/
DROP TABLE IF EXISTS `airports`;
CREATE TABLE `airports` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) DEFAULT '' COMMENT 'Public name of the airport',
    `iata_code` VARCHAR(16) DEFAULT '' COMMENT 'Official IATA code',
    `icao_code` VARCHAR(16) DEFAULT '' COMMENT 'Official ICAO code',
    `lat` DOUBLE DEFAULT 0 COMMENT 'lat',
    `lng` DOUBLE DEFAULT 0 COMMENT 'lng',
    `city` VARCHAR(255) DEFAULT '' COMMENT 'Airport metropolitan city name',
    `city_code` VARCHAR(16) DEFAULT '' COMMENT 'Airport metropolitan 3 letter city code',
    `un_locode` VARCHAR(16) DEFAULT '' COMMENT 'United Nations location code',
    `timezone` VARCHAR(16) DEFAULT '' COMMENT 'Airport location timezone',
    `country_code` VARCHAR(16) DEFAULT '' COMMENT 'ISO 2 country code',
    `departures`INT(10) DEFAULT 0 COMMENT 'Total departures from airport per year',
    `created_on` INT(10) UNSIGNED DEFAULT 0 COMMENT 'create at this time',
    `modified_on` INT(10) UNSIGNED DEFAULT 0 COMMENT 'update at this time',
    `deleted_on` INT(10) UNSIGNED DEFAULT 0 COMMENT 'delete at this time',
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='airports manage';