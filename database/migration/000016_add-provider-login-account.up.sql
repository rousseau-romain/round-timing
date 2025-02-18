ALTER TABLE user ADD COLUMN provider_login VARCHAR(255) DEFAULT NULL;

UPDATE user SET provider_login="discord";