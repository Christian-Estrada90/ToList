package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DbConnection() (*sql.DB, error) {
	// Cambia 'root' si usas otro usuario, y '12345' es tu contraseña actual
	connectionString := "root:root@tcp(127.0.0.1:3306)/crud_golang?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
