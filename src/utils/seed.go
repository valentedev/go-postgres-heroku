package utils

import (
	"database/sql"
	"fmt"

	"github.com/valentedev/go-postgres-heroku/src/usuarios"
)

func seed(db *sql.DB) {
	//apaga tabela USUARIOS anterior, caso ela exista, e cria uma nova tabela com os campos abaixo
	query1 := `
	DROP TABLE IF EXISTS vercod;
	DROP TABLE IF EXISTS usuarios;
	CREATE TABLE usuarios (
		id SERIAL PRIMARY KEY,
		criado_em TIMESTAMPTZ DEFAULT Now() NOT NULL,
		nome VARCHAR(50) NOT NULL,
		sobrenome VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		senha VARCHAR(100),
		admin boolean DEFAULT false NOT NULL,
		ativo boolean DEFAULT false NOT NULL
	);

	CREATE TABLE vercod (
		id SERIAL PRIMARY KEY,
		criado_em TIMESTAMPTZ DEFAULT Now() NOT NULL,
		usuario BIGINT REFERENCES usuarios (id) ON DELETE CASCADE NOT NULL,
		codigo VARCHAR(16) NOT NULL UNIQUE
	);
	`
	_, err := db.Exec(query1)
	if err != nil {
		fmt.Println(err)
	}

	us := usuarios.UsuariosSlice

	for x := range us {
		usuario := us[x]
		query2 := `
		INSERT INTO usuarios(nome, sobrenome, email, senha, admin, ativo)
		VALUES ($1,$2,$3,$4,$5,$6)`
		_, err = db.Exec(query2, usuario.Nome, usuario.Sobrenome, usuario.Email, usuario.Senha, usuario.Admin, usuario.Ativo)
		if err != nil {
			panic(err)
		}
	}
}
