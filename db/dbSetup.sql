CREATE USER IF NOT EXISTS 'rodrigovalente' WITH PASSWORD 'gph_pass';

CREATE DATABASE IF NOT EXISTS gph_db
OWNER rodrigovalente;

CREATE TABLE public.usuarios (
    id SERIAL PRIMARY KEY,
    criado_em TIMESTAMP DEFAULT Now() NOT NULL,
    nome VARCHAR(50),
    sobrenome VARCHAR(50),
    email VARCHAR(100) UNIQUE,
    perfil VARCHAR(50),
    mandato VARCHAR(50),
    foto VARCHAR(150),
    naturalidade VARCHAR(100)
);