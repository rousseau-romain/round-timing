CREATE TABLE spell (
    id int(11) AUTO_INCREMENT PRIMARY KEY,
    id_class int(11) NOT NULL REFERENCES class(id),
    name VARCHAR(255) NOT NULL,
    delay int(11) NOT NULL DEFAULT 1,
    is_global BOOLEAN NOT NULL DEFAULT FALSE,
    is_team BOOLEAN NOT NULL DEFAULT FALSE,
    is_self BOOLEAN NOT NULL DEFAULT FALSE,
    is_ending_caster BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO spell (id_class, name, delay) VALUES 
(1, "Armure Incandescente", 5),
(1, "Renvoi de Sort", 6),
(1, "Trêve", 8),
(1, "Armure Terrestre", 5),
(1, "Bouclier Féca", 6),
(1, "Glyphe Enflammé", 2),
(1, "Téléportation", 6),
(1, "Glyphe de Silence", 6),
(1, "Armure Venteuse", 5),
(1, "Glyphe d'immobilisation", 5),
(1, "Science du bâton", 6),
(1, "Glyphe Agressif", 2),
(1, "Armure Aqueuse", 5),
(1, "Immunité", 6),
(1, "Invocation de Dopeul Féca", 10);


INSERT INTO spell (id_class, name, delay) VALUES 
(2, "Cri de l'Ours", 2),
(2, "Bénédiction Animale ", 6),
(2, "Piqûre Motivante", 3),
(2, "Crapaud", 3),
(2, "Crocs du Mulou", 4),
(2, "Invocation de Dragonnet Rouge", 6),
(2, "Résistance Naturelle", 3),
(2, "Invocation de Bouftou", 2),
(2, "Invocation de Prespic", 4),
(2, "Invocation de Sanglier", 3),
(2, "Invocation de Bwork Mage", 4),
(2, "Invocation de Craqueleur", 4),
(2, "Invocation de Dopeul Osamodas", 8),
(2, "Laisse Spirituelle", 7);

INSERT INTO spell (id_class, name, delay) VALUES 
(3, "Sac Animé", 7),
(3, "Chance", 6),
(3, "Boîte de Pandore", 3),
(3, "Clé Réductrice", 5),
(3, "Cupidité", 6),
(3, "Maladresse de Masse", 10),
(3, "Accélération", 6),
(3, "Corruption", 5),
(3, "Pelle Animée", 5),
(3, "Coffre Animé", 63),
(3, "Invocation de Dopeul Enutrof", 10),
(3, "Retraite anticipée", 7);

INSERT INTO spell (id_class, name, delay) VALUES 
(4, "Invisibilité", 6),
(4, "Double", 6),
(4, "Repérage", 2),
(4, "Invisibilité d'Autrui", 5),
(4, "Piège Empoisonné", 2),
(4, "Concentration de Chakra", 6),
(4, "Piège d'Immobilisation", 5),
(4, "Piège de Silence", 4),
(4, "Pulsion de Chakra", 6),
(4, "Invocation de Dopeul Sram", 10),
(4, "Poisse", 5);

INSERT INTO spell (id_class, name, delay) VALUES 
(5, "Contre", 6),
(5, "Téléportation", 15),
(5, "Flou", 5),
(5, "Dévouement", 3),
(5, "Protection Aveuglante", 6),
(5, "Momification", 15),
(5, "Cadran de Xélor", 5),
(5, "Invocation de Dopeul Xélor", 10),
(5, "Raulebaque", 7);

INSERT INTO spell (id_class, name, delay) VALUES
(6, "Chance d'Ecaflip", 5),
(6, "Perception", 1),
(6, "Contrecoup", 3),
(6, "Trèfle", 6),
(6, "Roue de la Fortune", 4),
(6, "Griffe Invocatrice", 6),
(6, "Odorat", 7),
(6, "Réflexes", 6),
(6, "Invocation de Dopeul Ecaflip", 10);

INSERT INTO spell (id_class, name, delay) VALUES
(7, "Mot Stimulant", 5),
(7, "Mot de Prévention", 6),
(7, "Mot d'Epine", 3),
(7, "Mot d'Amitié", 6),
(7, "Mot d'Immobilisation", 7),
(7, "Mot de Silence", 3),
(7, "Mot d'Altruisme", 6),
(7, "Mot de Reconstitution", 7),
(7, "Invocation de Dopeul Eniripsa", 10),
(7, "Mot Lotof", 5);

INSERT INTO spell (id_class, name, delay) VALUES
(8, "Compulsion", 6),
(8, "Guide de Bravoure ", 3),
(8, "Souffle", 2),
(8, "Vitalité", 5),
(8, "Puissance", 5),
(8, "Colère de Iop", 4),
(8, "Invocation de Dopeul Iop", 10),
(8, "Brokle", 5);



INSERT INTO spell (id_class, name, delay) VALUES
(9, "Tir Eloigné", 5),
(9, "Flèche d'Expiation", 3),
(9, "Oeil de Taupe", 4),
(9, "Tir Critique", 5),
(9, "Flèche Punitive", 2),
(9, "Tir Puissant", 6),
(9, "Maîtrise de l'Arc", 5),
(9, "Invocation de Dopeul Crâ", 10),
(9, "Flèche de dispersion", 2);

INSERT INTO spell (id_class, name, delay) VALUES
(10, "Poison Paralysant", 2),
(10, "Ronce Apaisante", 4),
(10, "Puissance Sylvestre", 10),
(10, "Tremblement", 5),
(10, "La Sacrifiée", 2),
(10, "Connaissance des Poupées", 6),
(10, "Arbre", 3),
(10, "Vent Empoisonné", 7),
(10, "La Gonflable", 2),
(10, "Ronce Insolente", 3),
(10, "La Surpuissante", 5),
(10, "Invocation de Dopeul Sadida", 10),
(10, "Arbre de vie", 6);

INSERT INTO spell (id_class, name, delay) VALUES
(11, "Châtiment Forcé", 5),
(11, "Dérobade", 5),
(11, "Châtiment Agile", 5),
(11, "Châtiment Osé", 5),
(11, "Châtiment Spirituel", 5),
(11, "Sacrifice", 6),
(11, "Châtiment Vitalesque", 4),
(11, "Coopération", 3),
(11, "Transposition", 3),
(11, "Punition", 2),
(11, "Epée volante", 4),
(11, "Invocation de Dopeul Sacrieur", 10),
(11, "Douleur partagée", 5);

INSERT INTO spell (id_class, name, delay) VALUES
(12, "Stabilisation", 6),
(12, "Colère de Zatoïshwan", 6),
(12, "Pandanlku", 6),
(12, "Lien Spiritueux", 8),
(12, "Invocation de Dopeul Pandawa", 10),
(12, "Ivresse", 5);


INSERT INTO spell (id_class, name, delay) VALUES
(13, "Invocation de Chaferfu", 6),
(13, "Invocation d'Arakne", 6),
(13, "Boomerang perfide", 2),
(13, "Cawotte", 6),
(13, "Marteau de Moon", 3),
(13, "Maîtrise des Bâtons ", 6),
(13, "Maîtrise des Epées", 6),
(13, "Maîtrise des Arcs", 6),
(13, "Maîtrise des Marteaux", 6),
(13, "Maîtrise des Baguettes", 6),
(13, "Maîtrise des Dagues", 6),
(13, "Maîtrise des Pelles", 6),
(13, "Maîtrise des Haches", 6),
(13, "Libération", 2);
