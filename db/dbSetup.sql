-- CREATE TABLE usuarios (
--     id SERIAL PRIMARY KEY,
--     criado_em TIMESTAMP DEFAULT Now() NOT NULL,
--     nome VARCHAR(50),
--     sobrenome VARCHAR(50),
--     email VARCHAR(100) UNIQUE,
--     perfil VARCHAR(50),
--     mandato VARCHAR(50),
--     foto VARCHAR(250),
--     naturalidade VARCHAR(100)
-- );

CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    criado_em TIMESTAMP DEFAULT Now() NOT NULL,
    nome VARCHAR(50) NOT NULL,
    sobrenome VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    senha VARCHAR(100),
    acesso VARCHAR(50) NOT NULL
);

CREATE TABLE vercod (
    id SERIAL PRIMARY KEY,
    criado_em TIMESTAMP DEFAULT Now() NOT NULL,
    usuario BIGINT REFERENCES usuarios (id) ON DELETE CASCADE NOT NULL,
    codigo VARCHAR(16) NOT NULL UNIQUE
);

