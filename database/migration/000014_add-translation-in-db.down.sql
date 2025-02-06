ALTER TABLE class_translation DROP FOREIGN KEY fk_class_translation_language;
ALTER TABLE class_translation DROP FOREIGN KEY fk_class_translation_class;
DROP TABLE IF EXISTS class_translation;

ALTER TABLE spell_translation DROP FOREIGN KEY fk_spell_translation_language;
ALTER TABLE spell_translation DROP FOREIGN KEY fk_spell_translation_spell;
DROP TABLE IF EXISTS spell_translation;