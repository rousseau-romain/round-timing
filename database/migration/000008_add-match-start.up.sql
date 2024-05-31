ALTER TABLE `match` ADD COLUMN round INT NOT NULL DEFAULT 0;

CREATE TABLE match_player_spell (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    player_id INT NOT NULL,
    spell_id INT NOT NULL,
    round_before_recovery INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES `match`(id),
    FOREIGN KEY (player_id) REFERENCES `player`(id),
    FOREIGN KEY (spell_id) REFERENCES `spell`(id)
);
