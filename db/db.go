package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

const create string = `
  CREATE TABLE IF NOT EXISTS scheduler (
  id INTEGER NOT NULL PRIMARY KEY,
  date CHAR(8),
  comment TEXT,
  repeat VARCHAR(128),
  title TEXT NOT NULL
  );
  CREATE INDEX task_date ON scheduler (date);`

type DB struct {
	db *sql.DB
}

func NewDB(filePath string) (DB, error) {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return DB{}, err
	}

	return DB{db: db}, nil
}

var Data DB

func (db DB) Close() error {
	return db.db.Close()
}

func CheckDB() {

	dbFile := filepath.Join("scheduler.db")
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		fmt.Println(err)
		install = true
	}

	if install {
		_, err = os.Create("scheduler.db")
		if err != nil {
			fmt.Println(err)
			return
		}

		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = db.Exec(create)
		if err != nil {
			fmt.Println(err)
		}
		db.Close()
	}

}
