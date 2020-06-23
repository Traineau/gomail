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
    password VARCHAR(50),
    PRIMARY KEY (id)
)