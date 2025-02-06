ALTER TABLE user DROP FOREIGN KEY fk_user_language;
ALTER TABLE user DROP COLUMN id_language;

DROP TABLE IF EXISTS language;