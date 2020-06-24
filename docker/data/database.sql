/*
    This file is used by the docker-compose build command to build the mysql db
    
    Tables:
    * category : stores categories (id, name, desc, creation, update)
    * image : stores images (id, name, desc, type, creation, update, category ID)
    * tag : stores tags (id, name, creation date)
    * image_tag : links images to tags by ids (Many to Many relation)
*/

CREATE TABLE user (
    id BIGINT NOT NULL AUTO_INCREMENT, 
    username VARCHAR(50),
    email VARCHAR (255),
    password VARCHAR(255),
    PRIMARY KEY (id)
);

CREATE TABLE mailing_list (
    id BIGINT NOT NULL AUTO_INCREMENT, 
    name VARCHAR(50),
    description TEXT,
    PRIMARY KEY(id)
);

CREATE TABLE user_mailing_list (
    id_user BIGINT NOT NULL,
    id_mailing_list BIGINT NOT NULL
);

CREATE TABLE campaign (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50),
    description TEXT,
    id_mailing_list BIGINT NOT NULL,
    PRIMARY KEY(id)
);