
CREATE TABLE IF NOT EXISTS usuarios(
    id serial PRIMARY KEY, 
    nome varchar(100) NOT NULL, 
    nick varchar(100) NOT NULL UNIQUE, 
    email varchar(100) NOT NULL UNIQUE, 
    senha varchar(100) NOT NULL,
    criadoem timestamp DEFAULT current_timestamp NOT NULL,
    updateem timestamp DEFAULT current_timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS seguidores (
    usuario_id int NOT NULL,
    FOREIGN KEY(usuario_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

    seguidor_id int NOT NULL,
    FOREIGN KEY(seguidor_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,
     
    PRIMARY KEY(usuario_id, seguidor_id)
 );

CREATE TABLE IF NOT EXISTS publicacoes (
    id SERIAL PRIMARY KEY,
    titulo varchar(100) NOT NULL,
    conteudo varchar(250) NOT NULL,
    autor_id INT NOT NULL,
    curtidas INT DEFAULT 0,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    FOREIGN KEY (autor_id)
        REFERENCES usuarios (id)
        ON DELETE CASCADE
);
     