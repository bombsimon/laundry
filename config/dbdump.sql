-- MySQL dump 10.13  Distrib 5.7.19, for Linux (x86_64)
--
-- Host: localhost    Database: laundry
-- ------------------------------------------------------
-- Server version	5.7.19

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
-- Table structure for table `booker`
--

DROP TABLE IF EXISTS `booker`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `booker` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `identifier` varchar(100) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `pin` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `booker`
--

LOCK TABLES `booker` WRITE;
/*!40000 ALTER TABLE `booker` DISABLE KEYS */;
INSERT INTO `booker` VALUES (1,'1001','Simon Sawert','hej@johanhornsten.se',NULL,'1234'),(3,'',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `booker` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bookings`
--

DROP TABLE IF EXISTS `bookings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bookings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `book_date` date NOT NULL,
  `id_slots` int(11) NOT NULL,
  `id_booker` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `id_slots` (`id_slots`),
  KEY `id_booker` (`id_booker`),
  CONSTRAINT `bookings_ibfk_1` FOREIGN KEY (`id_slots`) REFERENCES `slots` (`id`),
  CONSTRAINT `bookings_ibfk_2` FOREIGN KEY (`id_booker`) REFERENCES `booker` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bookings`
--

LOCK TABLES `bookings` WRITE;
/*!40000 ALTER TABLE `bookings` DISABLE KEYS */;
INSERT INTO `bookings` VALUES (1,'2017-08-23',4,1),(2,'2017-09-12',7,1);
/*!40000 ALTER TABLE `bookings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `machines`
--

DROP TABLE IF EXISTS `machines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `machines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `info` varchar(100) DEFAULT NULL,
  `working` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machines`
--

LOCK TABLES `machines` WRITE;
/*!40000 ALTER TABLE `machines` DISABLE KEYS */;
INSERT INTO `machines` VALUES (1,'Electrolux 1',1),(2,'Electrolux 2',1);
/*!40000 ALTER TABLE `machines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notification_types`
--

DROP TABLE IF EXISTS `notification_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `notification_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(25) DEFAULT NULL,
  `description` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notification_types`
--

LOCK TABLES `notification_types` WRITE;
/*!40000 ALTER TABLE `notification_types` DISABLE KEYS */;
INSERT INTO `notification_types` VALUES (1,'on_release','Someone cancels their slot'),(2,'reminder','Before your slot start');
/*!40000 ALTER TABLE `notification_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `notifications` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_notification_types` int(11) NOT NULL,
  `id_bookings` int(11) NOT NULL,
  `ahead` int(4) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_notification_types` (`id_notification_types`),
  KEY `id_bookings` (`id_bookings`),
  CONSTRAINT `notifications_ibfk_1` FOREIGN KEY (`id_notification_types`) REFERENCES `notification_types` (`id`),
  CONSTRAINT `notifications_ibfk_2` FOREIGN KEY (`id_bookings`) REFERENCES `bookings` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notifications`
--

LOCK TABLES `notifications` WRITE;
/*!40000 ALTER TABLE `notifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `notifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `slots`
--

DROP TABLE IF EXISTS `slots`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `slots` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `week_day` enum('0','1','2','3','4','5','6') NOT NULL,
  `start_time` time NOT NULL,
  `end_time` time NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `slots`
--

LOCK TABLES `slots` WRITE;
/*!40000 ALTER TABLE `slots` DISABLE KEYS */;
INSERT INTO `slots` VALUES (1,'1','07:00:00','10:00:00'),(2,'1','10:00:00','14:00:00'),(3,'1','14:00:00','18:00:00'),(4,'1','18:00:00','22:00:00'),(5,'2','07:00:00','10:00:00'),(6,'3','07:00:00','10:00:00'),(7,'4','07:00:00','10:00:00'),(8,'5','07:00:00','10:00:00'),(9,'2','10:00:00','14:00:00'),(10,'3','10:00:00','14:00:00'),(11,'4','10:00:00','14:00:00'),(12,'5','10:00:00','14:00:00'),(13,'2','14:00:00','18:00:00'),(14,'3','14:00:00','18:00:00'),(15,'4','14:00:00','18:00:00'),(16,'5','14:00:00','18:00:00'),(17,'2','18:00:00','22:00:00'),(18,'3','18:00:00','22:00:00'),(19,'4','18:00:00','22:00:00'),(20,'5','18:00:00','22:00:00');
/*!40000 ALTER TABLE `slots` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `slots_machines`
--

DROP TABLE IF EXISTS `slots_machines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `slots_machines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_machines` int(11) NOT NULL,
  `id_slots` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UC_slots_machines` (`id_machines`,`id_slots`),
  KEY `slots_machines_ibfk_2` (`id_slots`),
  CONSTRAINT `slots_machines_ibfk_1` FOREIGN KEY (`id_machines`) REFERENCES `machines` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `slots_machines_ibfk_2` FOREIGN KEY (`id_slots`) REFERENCES `slots` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `slots_machines`
--

LOCK TABLES `slots_machines` WRITE;
/*!40000 ALTER TABLE `slots_machines` DISABLE KEYS */;
INSERT INTO `slots_machines` VALUES (1,1,1),(3,1,2),(5,1,3),(7,1,4),(8,1,7),(2,2,1),(4,2,2),(6,2,3),(10,2,4),(9,2,7);
/*!40000 ALTER TABLE `slots_machines` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-09-02 18:07:10
