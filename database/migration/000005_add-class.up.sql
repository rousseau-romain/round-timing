CREATE TABLE class (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO class (name) VALUES 
('Féca'),
('Osamodas'),
('Enutrof'),
('Sram'),
('Xélor'),
('Ecaflip'),
('Eniripsa'),
('Iop'),
('Crâ'),
('Sadida'),
('Sacrieur'),
('Pandawa'),
('Global');
