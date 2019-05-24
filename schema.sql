CREATE DATABASE IF NOT EXISTS flourishdb;

USE flourishdb;

CREATE TABLE strains
(
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    raceId INT
);

TRUNCATE TABLE strains;

CREATE TABLE races
(
    id INT PRIMARY KEY,
    name VARCHAR(50)
);

CREATE TABLE flavors
(
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

TRUNCATE TABLE flavors;

CREATE TABLE effects
(
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL
);

TRUNCATE TABLE effects;

CREATE TABLE strain_flavors
(
    id INT PRIMARY KEY,
    strainId INT NOT NULL,
    flavorId INT NOT NULL
);

TRUNCATE TABLE strain_flavors;

CREATE TABLE strain_effects
(
    id INT PRIMARY KEY,
    strainId INT NOT NULL,
    effectId INT NOT NULL
);

TRUNCATE TABLE strain_effects
