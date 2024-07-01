package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Scheduler struct {
	ID      int
	Title   string
	Comment int
	Repeat  string
	Date    string
}

type SchedulerService struct {
	store SchedulerStore
}

func NewSchedulerService(store SchedulerStore) SchedulerService {
	return SchedulerService{store: store}
}

const create string = `
  CREATE TABLE IF NOT EXISTS scheduler (
  id INTEGER NOT NULL PRIMARY KEY,
  date DATETIME NOT NULL,
  comment TEXT,
  repeat TEXT,
  title TEXT NOT NULL
  );`

const Value = 1048576

func main() {

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	store := NewSchedulerStore(db)
	// service := NewSchedulerService(store)

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	fmt.Println(dbFile)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	fmt.Println(install)

	if install == true {
		if _, err := store.db.Exec(create); err != nil {
			fmt.Println(err)
			return
		}
	}

}
