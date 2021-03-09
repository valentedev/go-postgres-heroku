package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// SelectByID ...
func SelectByID(db *sql.DB, fields string, table string, id int) (string, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id=%d", fields, table, id)
	rows, _ := db.Query(query)
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	m := make(map[string]interface{})
	columnPointers := make([]interface{}, len(cols))

	for rows.Next() {
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			mj, _ := json.Marshal(m)
			return string(mj), err
		}
		for j, colName := range cols {
			val := columnPointers[j].(*interface{})
			m[colName] = *val
		}
	}

	mjson, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(mjson), nil
}

// SelectByEmail ...
func SelectByEmail(db *sql.DB, fields string, table string, email string) (string, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE email=%s", fields, table, email)
	rows, _ := db.Query(query)
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	m := make(map[string]interface{})
	columnPointers := make([]interface{}, len(cols))

	for rows.Next() {
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			mj, _ := json.Marshal(m)
			return string(mj), err
		}
		for j, colName := range cols {
			val := columnPointers[j].(*interface{})
			m[colName] = *val
		}
	}

	mjson, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(mjson), nil
}
