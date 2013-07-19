DROP TABLE IF EXISTS `gamers`;
CREATE TABLE `gamers` (
  `gamerId` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(32) CHARACTER SET ascii NOT NULL,
  PRIMARY KEY (`gamerId`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
  `sessionId` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `gamerId` int(11) unsigned NOT NULL,
  `token` char(32) CHARACTER SET ascii NOT NULL,
  `ip` char(50) NOT NULL,
  `start` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `latest` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `end` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`sessionId`),
  KEY `gamerId` (`gamerId`),
  CONSTRAINT `sessions_ibfk_1` FOREIGN KEY (`gamerId`) REFERENCES `gamers` (`gamerId`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2000 DEFAULT CHARSET=utf8;

