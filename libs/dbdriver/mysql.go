package dbdriver

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func MysqlConn() (*sql.DB, error) {
	dsn := "root:amitabha@tcp(localhost:3306)/test?charset=utf8"
	db, err := sql.Open("mysql", dsn)

	return db, err
}
