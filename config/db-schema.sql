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
    name    VARCHAR(100),
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

    FOREIGN KEY (id_machines) REFERENCES machines(id),
    FOREIGN KEY (id_slots)    REFERENCES slots(id),
    CONSTRAINT UC_slots_machines UNIQUE (id_machines, id_slots)
);

CREATE TABLE `bookings` (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    book_date   DATE NOT NULL, 
    id_slots    INT NOT NULL,
    id_booker   INT NOT NULL,

    FOREIGN KEY (id_slots)  REFERENCES slots(id),
    FOREIGN KEY (id_booker) REFERENCES booker(id)
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

    FOREIGN KEY (id_notification_types) REFERENCES notification_types(id),
    FOREIGN KEY (id_bookings)           REFERENCES bookings(id)
);
