CREATE TABLE language (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    locale VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO language (locale) VALUES ('en'), ('fr'), ('it'), ('es'), ('pt');

ALTER TABLE user ADD COLUMN id_language INT;

UPDATE user SET id_language=1;

ALTER TABLE user MODIFY COLUMN id_language INT NOT NULL;

ALTER TABLE user ADD CONSTRAINT fk_user_language FOREIGN KEY (id_language) REFERENCES language(id);