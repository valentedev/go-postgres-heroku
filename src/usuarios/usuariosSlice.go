package usuarios

import (
	"reflect"
	"strings"
)

//Usuarios é uma struct que servirá como base para o futuro banco de dados
type Usuarios struct {
	ID        int    `json:"id"`
	Nome      string `json:"nome"`
	Sobrenome string `json:"sobrenome"`
	Email     string `json:"email"`
	Senha     string `json:"senha"`
	Admin     bool   `json:"admin"`
	Ativo     bool   `json:"ativo"`
}

// StructFieldsToString ...
func (u Usuarios) StructFieldsToString() string {
	var fields string
	us := reflect.ValueOf(u)
	typeOfU := us.Type()
	for i := 0; i < us.NumField(); i++ {
		fields = fields + ", " + strings.ToLower(typeOfU.Field(i).Name)
	}
	fields = fields[2:]
	return fields
}

//UsuariosSlice vai reunir um ou mais Usuários que serão renderizados no Template index.html
var UsuariosSlice = []Usuarios{
	{
		ID:        1,
		Nome:      "Rodrigo",
		Sobrenome: "Valente",
		Email:     "valentedev.rodrigo@gmail.com",
		Senha:     "senha",
		Admin:     true,
		Ativo:     true,
	},
	{
		ID:        2,
		Nome:      "Steve",
		Sobrenome: "Jobs",
		Email:     "sj@email.com",
		Senha:     "",
		Admin:     false,
		Ativo:     false,
	},
	{
		ID:        3,
		Nome:      "Bill",
		Sobrenome: "Gates",
		Email:     "bg@email.com",
		Senha:     "",
		Admin:     false,
		Ativo:     true,
	},
	{
		ID:        4,
		Nome:      "Rodrigo",
		Sobrenome: "Hotmail Valente",
		Email:     "rodrigovalente@hotmail.com",
		Senha:     "",
		Admin:     true,
		Ativo:     true,
	},
}

// //Usuarios é uma struct que servirá como base para o futuro banco de dados
// type Usuarios struct {
// 	ID           int
// 	Nome         string
// 	Sobrenome    string
// 	Email        string
// 	Perfil       string
// 	Mandato      string
// 	Foto         string
// 	Naturalidade string
// }

// //UsuariosSlice vai reunir um ou mais Usuários que serão renderizados no Template index.html
// var UsuariosSlice = []Usuarios{
// 	{
// 		ID:           1,
// 		Nome:         "Rodrigo",
// 		Sobrenome:    "Valente",
// 		Email:        "valentedev.rodrigo@gmail.com",
// 		Perfil:       "admin",
// 		Mandato:      "",
// 		Foto:         "",
// 		Naturalidade: "São Paulo-SP",
// 	},
// 	{
// 		ID:           2,
// 		Nome:         "Campos",
// 		Sobrenome:    "Sales",
// 		Email:        "cs@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1898-1902",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b1/Campos_Sales.jpg/220px-Campos_Sales.jpg",
// 		Naturalidade: "Campinas-SP",
// 	},
// 	{
// 		ID:           3,
// 		Nome:         "Rodrigues",
// 		Sobrenome:    "Alves",
// 		Email:        "ra@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1902-1906",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d3/Rodrigues_Alves_3.jpg/220px-Rodrigues_Alves_3.jpg",
// 		Naturalidade: "Guaratinguetá-SP",
// 	}, {
// 		ID:           4,
// 		Nome:         "Afonso",
// 		Sobrenome:    "Pena",
// 		Email:        "ap@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1906-1909",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/0/0a/Afonso_Pena.jpg/200px-Afonso_Pena.jpg",
// 		Naturalidade: "Santa Bárbara-MG",
// 	},
// 	{
// 		ID:           5,
// 		Nome:         "Nilo",
// 		Sobrenome:    "Peçanha",
// 		Email:        "np@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1906-1909",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c0/Nilo_Pe%C3%A7anha_02.jpg/230px-Nilo_Pe%C3%A7anha_02.jpg",
// 		Naturalidade: "Campos dos Goytacazes-RJ",
// 	},
// 	{
// 		ID:           6,
// 		Nome:         "Hermes",
// 		Sobrenome:    "da Fonseca",
// 		Email:        "hf@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1910-1914",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b3/Hermes_da_Fonseca_%281910%29.jpg/220px-Hermes_da_Fonseca_%281910%29.jpg",
// 		Naturalidade: "São Gabriel-RS",
// 	},
// 	{
// 		ID:           7,
// 		Nome:         "Venceslau",
// 		Sobrenome:    "Brás",
// 		Email:        "vb@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1914-1918",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cb/Venceslau_Br%C3%A1s.jpg/220px-Venceslau_Br%C3%A1s.jpg",
// 		Naturalidade: "Brazópolis-MG",
// 	},
// 	{
// 		ID:           8,
// 		Nome:         "Delfim",
// 		Sobrenome:    "Moreira",
// 		Email:        "dm@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1918-1919",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a7/Delfim_Moreira_%281918%29.jpg/220px-Delfim_Moreira_%281918%29.jpg",
// 		Naturalidade: "Cristina-MG",
// 	},
// 	{
// 		ID:           9,
// 		Nome:         "Epitácio",
// 		Sobrenome:    "Pessoa",
// 		Email:        "ep@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1919-1922",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/1/19/Epitacio_Pessoa_%281919%29.jpg/330px-Epitacio_Pessoa_%281919%29.jpg",
// 		Naturalidade: "Umbuzeiro-PB",
// 	},
// 	{
// 		ID:           10,
// 		Nome:         "Artur",
// 		Sobrenome:    "Bernardes",
// 		Email:        "ab@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1922-1926",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/1/15/Artur_Bernardes_%281922%29.jpg/330px-Artur_Bernardes_%281922%29.jpg",
// 		Naturalidade: "Viçosa-MG",
// 	},
// 	{
// 		ID:           11,
// 		Nome:         "Washington",
// 		Sobrenome:    "Luis",
// 		Email:        "wl@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1926-1930",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fe/Washington_Lu%C3%ADs_%28foto%29.jpg/330px-Washington_Lu%C3%ADs_%28foto%29.jpg",
// 		Naturalidade: "Macaé-RJ",
// 	},
// 	// {
// 	// 	ID:           12,
// 	// 	Nome:         "Julio",
// 	// 	Sobrenome:    "Prestes",
// 	// 	Email:        "jp@gmail.com",
// 	// 	Perfil:       "no-admin",
// 	// },
// 	{
// 		ID:           13,
// 		Nome:         "Getulio",
// 		Sobrenome:    "Vargas",
// 		Email:        "gv@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1930-1945, 1951-1954",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/5/50/Getulio_Vargas_%281930%29.jpg/368px-Getulio_Vargas_%281930%29.jpg",
// 		Naturalidade: "São Borja-RS",
// 	},
// 	{
// 		ID:           14,
// 		Nome:         "Jose",
// 		Sobrenome:    "Linhares",
// 		Email:        "jl@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1945-1946",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Jos%C3%A9_Linhares_TSE.jpg/300px-Jos%C3%A9_Linhares_TSE.jpg",
// 		Naturalidade: "Guaramiranga-CE",
// 	},
// 	{
// 		ID:           15,
// 		Nome:         "Eurico",
// 		Sobrenome:    "Gaspar Dutra",
// 		Email:        "egd@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1946-1951",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/5/55/GASPARDUTRA.jpg/330px-GASPARDUTRA.jpg",
// 		Naturalidade: "Cuiabá-MT",
// 	},
// 	{
// 		ID:           16,
// 		Nome:         "Café",
// 		Sobrenome:    "Filho",
// 		Email:        "cf@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1954-1955",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d3/Caf%C3%A9_Filho.jpg/330px-Caf%C3%A9_Filho.jpg",
// 		Naturalidade: "Natal-RN",
// 	},
// 	{
// 		ID:           17,
// 		Nome:         "Carlos",
// 		Sobrenome:    "Luz",
// 		Email:        "cl@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1955",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1d/Carlos_Luz_Oficial.jpg/330px-Carlos_Luz_Oficial.jpg",
// 		Naturalidade: "Três Corações-MG",
// 	},
// 	{
// 		ID:           18,
// 		Nome:         "Nereu",
// 		Sobrenome:    "Ramos",
// 		Email:        "nr@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1955-1956",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d2/Presidente_Nereu_Ramos.jpg/330px-Presidente_Nereu_Ramos.jpg",
// 		Naturalidade: "Lages-SC",
// 	},
// 	{
// 		ID:           19,
// 		Nome:         "Juscelino",
// 		Sobrenome:    "Kubitscheck",
// 		Email:        "jk@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1956-1961",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Juscelino.jpg/330px-Juscelino.jpg",
// 		Naturalidade: "Diamantina-MG",
// 	},
// 	{
// 		ID:           20,
// 		Nome:         "Jânio",
// 		Sobrenome:    "Quadros",
// 		Email:        "jq@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1961",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Janio_Quadros.png/330px-Janio_Quadros.png",
// 		Naturalidade: "Campo Grande-RJ",
// 	},
// 	{
// 		ID:           21,
// 		Nome:         "Ranieri",
// 		Sobrenome:    "Mazzilli",
// 		Email:        "rm@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1961 e 1964",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Pascoal_Ranieri_Mazzilli%2C_presidente_da_Rep%C3%BAblica..tif/lossy-page1-330px-Pascoal_Ranieri_Mazzilli%2C_presidente_da_Rep%C3%BAblica..tif.jpg",
// 		Naturalidade: "Caconde-MG",
// 	},
// 	{
// 		ID:           22,
// 		Nome:         "João",
// 		Sobrenome:    "Goulart",
// 		Email:        "jg@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1961-1964",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/d/dc/Jango.jpg/330px-Jango.jpg",
// 		Naturalidade: "São Borja-RS",
// 	},
// 	{
// 		ID:           23,
// 		Nome:         "Humberto",
// 		Sobrenome:    "Castelo Branco",
// 		Email:        "hcb@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1964-1967",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/d/df/Castelobranco.jpg/330px-Castelobranco.jpg",
// 		Naturalidade: "Fortaleza-CE",
// 	},
// 	{
// 		ID:           24,
// 		Nome:         "Artur",
// 		Sobrenome:    "da Costa e Silva",
// 		Email:        "acs@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1967-1969",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/6/66/Costa_e_Silva.jpg/330px-Costa_e_Silva.jpg",
// 		Naturalidade: "Taquari-RS",
// 	},
// 	// {
// 	// 	ID:        25,
// 	// 	Nome:      "Pedro",
// 	// 	Sobrenome: "Aleixo",
// 	// 	Email:     "pa@gmail.com",
// 	// 	Perfil:    "no-admin",
// 	// },
// 	{
// 		ID:           26,
// 		Nome:         "Emilio",
// 		Sobrenome:    "Garrastazu Médici",
// 		Email:        "egm@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1969-1974",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/9/92/Em%C3%ADlio_Garrastazu_M%C3%A9dici%2C_presidente_da_Rep%C3%BAblica._%28cropped%29.tif/lossy-page1-330px-Em%C3%ADlio_Garrastazu_M%C3%A9dici%2C_presidente_da_Rep%C3%BAblica._%28cropped%29.tif.jpg",
// 		Naturalidade: "Bagé-RS",
// 	},
// 	{
// 		ID:           27,
// 		Nome:         "Ernesto",
// 		Sobrenome:    "Geisel",
// 		Email:        "eg@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1974-1979",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/c/cc/Geisel-colour-president.jpg",
// 		Naturalidade: "Bento Gonçalves-RS",
// 	},
// 	{
// 		ID:           28,
// 		Nome:         "João",
// 		Sobrenome:    "Figueiredo",
// 		Email:        "jf@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1979-1985",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/a/af/Figueiredo_%28colour%29.jpg/330px-Figueiredo_%28colour%29.jpg",
// 		Naturalidade: "Rio de Janeiro-RJ",
// 	},
// 	{
// 		ID:           29,
// 		Nome:         "José",
// 		Sobrenome:    "Sarney",
// 		Email:        "js@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1985-1990",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fb/Foto_Oficial_Sarney_EBC.jpg/220px-Foto_Oficial_Sarney_EBC.jpg",
// 		Naturalidade: "Pinheiro-MA",
// 	},
// 	{
// 		ID:           30,
// 		Nome:         "Fernando",
// 		Sobrenome:    "Collor",
// 		Email:        "fc@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1990-1992",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/6/63/Fernando_Collor_de_Mello.jpg",
// 		Naturalidade: "Rio de Janeiro-RJ",
// 	},
// 	{
// 		ID:           31,
// 		Nome:         "Itamar",
// 		Sobrenome:    "Franco",
// 		Email:        "if@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1992-1995",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Itamar_Augusto_Cautiero_Franco.gif/330px-Itamar_Augusto_Cautiero_Franco.gif",
// 		Naturalidade: "Mar territorial brasileiro",
// 	},
// 	{
// 		ID:           32,
// 		Nome:         "Fernando",
// 		Sobrenome:    "Henrique Cardoso",
// 		Email:        "fhc@gmail.com",
// 		Perfil:       "no-admin",
// 		Mandato:      "1995-2003",
// 		Foto:         "https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Fernando_Henrique_Cardoso_%281999%29.jpg/330px-Fernando_Henrique_Cardoso_%281999%29.jpg",
// 		Naturalidade: "Rio de Janeiro-RJ",
// 	},
// }
