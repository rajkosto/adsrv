DROP TABLE IF EXISTS `crexs`;
CREATE TABLE `crexs`  (
  `crexId` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `minSize` smallint(5) UNSIGNED NOT NULL DEFAULT 0,
  `rotDuration` smallint(5) UNSIGNED NOT NULL DEFAULT 0,
  `minAngleDeg` tinyint(3) UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`crexId`)
) AUTO_INCREMENT = 10001;

DROP TABLE IF EXISTS `medias`;
CREATE TABLE `medias`  (
  `mediaId` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `mimeType` smallint(5) UNSIGNED NOT NULL,
  `fileSize` int(11) UNSIGNED NOT NULL,
  `fileMd5` varchar(32) NOT NULL,
  `filePath` varchar(255) NOT NULL,
  `crexId` int(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`mediaId`),
  INDEX `crexId`(`crexId`),
  CONSTRAINT `medias_ibfk_1` FOREIGN KEY (`crexId`) REFERENCES `crexs` (`crexId`) ON DELETE CASCADE ON UPDATE CASCADE
) AUTO_INCREMENT = 100001 CHARACTER SET = ascii COLLATE = ascii_general_ci;

DROP TABLE IF EXISTS `zones`;
CREATE TABLE `zones`  (
  `zoneId` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`zoneId`)
) AUTO_INCREMENT = 101 CHARACTER SET = ascii COLLATE = ascii_general_ci;

DROP TABLE IF EXISTS `targets`;
CREATE TABLE `targets`  (
  `targetId` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `zoneId` int(11) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`targetId`),
  INDEX `zoneId`(`zoneId`),
  CONSTRAINT `targets_ibfk_1` FOREIGN KEY (`zoneId`) REFERENCES `zones` (`zoneId`) ON DELETE CASCADE ON UPDATE CASCADE
) AUTO_INCREMENT = 1001 CHARACTER SET = ascii COLLATE = ascii_general_ci;

DROP TABLE IF EXISTS `target_medias`;
CREATE TABLE `target_medias`  (
  `targetId` int(11) UNSIGNED NOT NULL,
  `mediaId` int(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`targetId`, `mediaId`),
  INDEX `mediaId`(`mediaId`),
  CONSTRAINT `target_medias_ibfk_1` FOREIGN KEY (`targetId`) REFERENCES `targets` (`targetId`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `target_medias_ibfk_2` FOREIGN KEY (`mediaId`) REFERENCES `medias` (`mediaId`) ON DELETE CASCADE ON UPDATE CASCADE
);

DROP TABLE IF EXISTS `gamers`;
CREATE TABLE `gamers`  (
  `gamerId` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid` char(32) NOT NULL,
  PRIMARY KEY (`gamerId`),
  UNIQUE INDEX `uuid`(`uuid`)
) AUTO_INCREMENT = 101 CHARACTER SET = ascii COLLATE = ascii_general_ci;

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions`  (
  `sessionId` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `gamerId` int(11) UNSIGNED NOT NULL,
  `token` char(32) NOT NULL,
  `ip` char(50) NOT NULL,
  `start` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `latest` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `end` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`sessionId`),
  INDEX `gamerId`(`gamerId`),
  CONSTRAINT `sessions_ibfk_1` FOREIGN KEY (`gamerId`) REFERENCES `gamers` (`gamerId`) ON DELETE CASCADE ON UPDATE CASCADE
) AUTO_INCREMENT = 10001 CHARACTER SET = ascii COLLATE = ascii_general_ci;

DROP TABLE IF EXISTS `zone_visits`;
CREATE TABLE `zone_visits`  (
  `visitId` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `sessionId` int(11) UNSIGNED NOT NULL,
  `zoneId` int(11) UNSIGNED NOT NULL,
  `timestampMs` bigint(20) UNSIGNED NOT NULL,
  PRIMARY KEY (`visitId`),
  INDEX `sessionId`(`sessionId`),
  INDEX `zoneId`(`zoneId`),
  CONSTRAINT `zone_visits_ibfk_1` FOREIGN KEY (`sessionId`) REFERENCES `sessions` (`sessionId`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `zone_visits_ibfk_2` FOREIGN KEY (`zoneId`) REFERENCES `zones` (`zoneId`) ON DELETE CASCADE ON UPDATE CASCADE
) AUTO_INCREMENT = 1000001;
