package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

const create string = `
  CREATE TABLE IF NOT EXISTS scheduler (
  id INTEGER NOT NULL PRIMARY KEY,
  date TEXT,
  comment TEXT,
  repeat TEXT,
  title TEXT NOT NULL
  );
  CREATE INDEX task_date ON scheduler (date);`

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
