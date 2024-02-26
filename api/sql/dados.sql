insert into usuarios (nome, nick, email, senha) 
values 
 ("natias", "matias01", "matias@gmail.com","$2a$10$CDsA75bXzIrY25avsMCWNOx8e/iLEhxQirCLhknn9hwxHr4JT0NMu"),
 ("lucas 2", "lucas01", "lucas@gmail.com", "$2a$10$CDsA75bXzIrY25avsMCWNOx8e/iLEhxQirCLhknn9hwxHr4JT0NMu"),
 ("pedro 3", "pedro01", "pedro@gmail.com", "$2a$10$CDsA75bXzIrY25avsMCWNOx8e/iLEhxQirCLhknn9hwxHr4JT0NMu");

insert into seguidores(usuario_id, seguidor_id) 
values 
(35, 36),
(37, 35),
(35, 37);

insert into publicacoes(titulo, conteudo,autor_id)
values
("Publicacao do usuario 24", "Publicacao do usuario 24", 24),
("Publicacao do usuario 25", "Publicacao do usuario 25", 25),
("Publicacao do usuario 26", "Publicacao do usuario 26", 26);