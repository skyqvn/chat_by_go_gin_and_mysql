-- MySQL dump 10.13  Distrib 8.0.29, for Win64 (x86_64)
--
-- Host: localhost    Database: chat
-- ------------------------------------------------------
-- Server version	8.0.29

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `chatgroup`
--

DROP TABLE IF EXISTS `chatgroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chatgroup` (
  `name` char(25) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT NULL,
  `password` char(20) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT NULL,
  `introduce` char(200) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT NULL,
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chatgroup`
--

LOCK TABLES `chatgroup` WRITE;
/*!40000 ALTER TABLE `chatgroup` DISABLE KEYS */;
INSERT INTO `chatgroup` VALUES ('sky','rtrt','',22),('sky','1111','',26),('sky','----','',27),('sky','123-','',28),('sky','1','11',29),('sky','12345678','',30),('sky','__________','',31);
/*!40000 ALTER TABLE `chatgroup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `member`
--

DROP TABLE IF EXISTS `member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `member` (
  `owner` bigint unsigned DEFAULT NULL,
  `chatgroup` bigint unsigned DEFAULT NULL,
  KEY `chatgroup` (`chatgroup`) USING BTREE,
  KEY `owner` (`owner`) USING BTREE,
  CONSTRAINT `chatgroup` FOREIGN KEY (`chatgroup`) REFERENCES `chatgroup` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `owner` FOREIGN KEY (`owner`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `member`
--

LOCK TABLES `member` WRITE;
/*!40000 ALTER TABLE `member` DISABLE KEYS */;
INSERT INTO `member` VALUES (23,22),(22,26),(22,27),(22,28),(22,29),(22,30),(22,31);
/*!40000 ALTER TABLE `member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `report`
--

DROP TABLE IF EXISTS `report`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `report` (
  `chatgroup` bigint unsigned NOT NULL,
  `owner` bigint unsigned NOT NULL,
  `value` text CHARACTER SET utf8mb3 COLLATE utf8_general_ci,
  `send_time` datetime DEFAULT CURRENT_TIMESTAMP,
  KEY `re_chatgroup` (`chatgroup`) USING BTREE,
  KEY `userid` (`owner`) USING BTREE,
  CONSTRAINT `re_chatgroup` FOREIGN KEY (`chatgroup`) REFERENCES `chatgroup` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `userid` FOREIGN KEY (`owner`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `report`
--

LOCK TABLES `report` WRITE;
/*!40000 ALTER TABLE `report` DISABLE KEYS */;
INSERT INTO `report` VALUES (22,23,'  trtr','2022-07-12 16:29:52'),(22,23,'ttttttyyyyytyytt','2022-07-12 16:29:59'),(22,23,'ffffffff','2022-07-12 16:30:05');
/*!40000 ALTER TABLE `report` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` char(25) CHARACTER SET utf8mb3 COLLATE utf8_general_ci NOT NULL DEFAULT '匿名用户',
  `password` char(20) CHARACTER SET utf8mb3 COLLATE utf8_general_ci NOT NULL DEFAULT '000000',
  `introduce` char(200) CHARACTER SET utf8mb3 COLLATE utf8_general_ci DEFAULT NULL,
  `login_code` bigint unsigned DEFAULT NULL,
  `last_login_time` datetime DEFAULT (now()),
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_name_on_user` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (21,'qwe','222','',1321432352,'2022-07-08 09:39:15'),(22,'sky','2009917','',10207340356411239215,'2022-08-03 03:16:52'),(23,'tttt','tttt','',1321432352,'2022-07-08 09:39:15');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-08-03 11:37:16
