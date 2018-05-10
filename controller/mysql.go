package controller

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var GlobalMySQLDB *sql.DB

func init() {
	db, err := NewMySQLDB("root:@tcp(127.0.0.1:3306)/Goes")
	if err != nil {
		panic(err)
	}
	GlobalMySQLDB = db
}

func NewMySQLDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
