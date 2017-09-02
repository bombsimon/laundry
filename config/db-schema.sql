DROP DATABASE IF EXISTS `laundry`;

CREATE DATABASE `laundry`;
USE `laundry;

CREATE TABLE `booker` (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    identifier  VARCHAR(100) NOT NULL, -- i.e. apartment no
    name        VARCHAR(100),
    email       VARCHAR(100),
    phone       VARCHAR(20),
    pin         VARCHAR(100)
);

CREATE TABLE `machines` (
    id      INT PRIMARY KEY AUTO_INCREMENT,
    info    VARCHAR(100),
    working TINYINT(1) DEFAULT 1
);

CREATE TABLE `slots` (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    week_day    ENUM('0', '1', '2', '3', '4', '5', '6') NOT NULL,
    start_time  TIME NOT NULL,
    end_time    TIME NOT NULL
);

CREATE TABLE `slots_machines` (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    id_machines INT NOT NULL,
    id_slots    INT NOT NULL,

    FOREIGN KEY (id_machines) REFERENCES machines(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (id_slots)    REFERENCES slots(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT UC_slots_machines UNIQUE (id_machines, id_slots)
);

CREATE TABLE `bookings` (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    book_date   DATE NOT NULL, 
    id_slots    INT NOT NULL,
    id_booker   INT NOT NULL,

    FOREIGN KEY (id_slots)  REFERENCES slots(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (id_booker) REFERENCES booker(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE `notification_types` (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(25) NOT NULL,
    description VARCHAR(255) NOT NULL,
);

CREATE TABLE `notifications` (
    id                      INT PRIMARY KEY AUTO_INCREMENT,
    id_notification_types   INT NOT NULL,
    id_bookings             INT NOT NULL,
    ahead                   INT(4),

    FOREIGN KEY (id_notification_types) REFERENCES notification_types(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (id_bookings)           REFERENCES bookings(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- Test data
/*
INSERT INTO `booker` VALUES (1,'1001','Some User','some.email@domain.com',NULL,'1234'),(2,'1002','Another User',NULL,NULL,NULL);

INSERT INTO `machines` VALUES (1,'Electrolux 1',1),(2,'Electrolux 2',1),(3,'Broken machine',0);

INSERT INTO `slots` VALUES (1,'1','07:00:00','10:00:00'),(2,'1','10:00:00','14:00:00'),(3,'1','14:00:00','18:00:00'),(4,'1','18:00:00','22:00:00'),(5,'2','07:00:00','10:00:00'),(6,'3','07:00:00','10:00:00'),(7,'4','07:00:00','10:00:00'),(8,'5','07:00:00','10:00:00'),(9,'2','10:00:00','14:00:00'),(10,'3','10:00:00','14:00:00'),(11,'4','10:00:00','14:00:00'),(12,'5','10:00:00','14:00:00'),(13,'2','14:00:00','18:00:00'),(14,'3','14:00:00','18:00:00'),(15,'4','14:00:00','18:00:00'),(16,'5','14:00:00','18:00:00'),(17,'2','18:00:00','22:00:00'),(18,'3','18:00:00','22:00:00'),(19,'4','18:00:00','22:00:00'),(20,'5','18:00:00','22:00:00');

INSERT INTO `slots_machines` VALUES (1,1,1),(3,1,2),(5,1,3),(7,1,4),(8,1,7),(2,2,1),(4,2,2),(6,2,3),(10,2,4),(9,2,7);

INSERT INTO `bookings` VALUES (1,'2017-08-23',4,1),(2,'2017-09-12',7,1),(3,'2017-09-12',8,1);

INSERT INTO `notification_types` VALUES (1,'on_release','Someone cancels their slot'),(2,'reminder','Before your slot start');
*/
