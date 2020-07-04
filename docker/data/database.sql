/*
    This file is used by the docker-compose build command to build the mysql db
    
    Tables:
    * category : stores categories (id, name, desc, creation, update)
    * image : stores images (id, name, desc, type, creation, update, category ID)
    * tag : stores tags (id, name, creation date)
    * image_tag : links images to tags by ids (Many to Many relation)
*/

CREATE TABLE api_user (
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

CREATE TABLE recipient (
    id BIGINT NOT NULL AUTO_INCREMENT, 
    email VARCHAR(255),
    first_name VARCHAR(50) DEFAULT "",
    last_name VARCHAR(50) DEFAULT "",
    username VARCHAR(50) DEFAULT "",
    PRIMARY KEY (id)
);

CREATE TABLE recipient_mailing_list (
    id_recipient BIGINT NOT NULL,
    id_mailing_list BIGINT NOT NULL
);

CREATE TABLE campaign (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) DEFAULT "",
    description TEXT,
    template_name VARCHAR(255) DEFAULT "",
    template_path VARCHAR(255) DEFAULT "",
    id_mailing_list BIGINT NOT NULL,
    PRIMARY KEY(id)
);