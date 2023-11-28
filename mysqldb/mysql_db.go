package mysqldb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{conn: db}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Query(sql string, args ...interface{}) (string, error) {
	rows, err := db.conn.Query(sql, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func (db *DB) Get(key string) ([]byte, error) {
	query := "SELECT phone_number FROM contacts WHERE name = ?"
	result, err := db.Query(query, key)
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}
