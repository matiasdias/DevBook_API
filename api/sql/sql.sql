
CREATE TABLE IF NOT EXISTS  usuarios(
    id serial PRIMARY KEY, 
    nome varchar(100) NOT NULL, 
    nick varchar(100) NOT NULL UNIQUE, 
    email varchar(100) NOT NULL UNIQUE, 
    senha varchar(100) NOT NULL,
    criadoem timestamp DEFAULT current_timestamp NOT NULL
);

/*CREATE TABLE seguidores (
    usuario_id int not null,
    FOREIGN KEY(usuario_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

    seguidor_id int not null,
    FOREIGN KEY(seguidor_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,
     
    primary key(usuario_id, seguidor_id)
 );
*/
CREATE TABLE IF NOT EXISTS publicacoes (
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(100) NOT NULL,
    conteudo TEXT NOT NULL,
    autor_id INT NOT NULL,
    curtidas INT DEFAULT 0,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    FOREIGN KEY (autor_id)
        REFERENCES usuarios (id)
        ON DELETE CASCADE
);
     