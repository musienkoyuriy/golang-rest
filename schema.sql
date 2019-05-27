CREATE DATABASE IF NOT EXISTS flourishdb;

USE flourishdb;

CREATE TABLE IF NOT EXISTS strains
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    race VARCHAR(50) NOT NULL
);

TRUNCATE TABLE strains;

CREATE TABLE IF NOT EXISTS flavors
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    UNIQUE (name)
);

TRUNCATE TABLE flavors;

CREATE TABLE IF NOT EXISTS effects
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    UNIQUE (name)
);

TRUNCATE TABLE effects;

CREATE TABLE IF NOT EXISTS strain_flavors
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    strainId INT NOT NULL,
    flavorId INT NOT NULL
    -- FOREIGN KEY (strainId)  REFERENCES strains (id),
    -- FOREIGN KEY (flavorId)  REFERENCES flavors (id)
);

TRUNCATE TABLE strain_flavors;

CREATE TABLE IF NOT EXISTS strain_effects
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    strainId INT NOT NULL,
    effectId INT NOT NULL
    -- FOREIGN KEY (strainId) REFERENCES strains (id),
    -- FOREIGN KEY (effectId) REFERENCES effects (id)
);

TRUNCATE TABLE strain_effects;
