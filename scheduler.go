package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SchedulerStore struct {
	db *sql.DB
}

func NewSchedulerStore(db *sql.DB) SchedulerStore {
	return SchedulerStore{db: db}
}

func (s SchedulerStore) Add(p Scheduler) (int, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (title, comment, repeat, date) VALUES (:client, :status, :address, :created_at)",
		sql.Named("title", p.Title),
		sql.Named("comment", p.Comment),
		sql.Named("repeat", p.Repeat),
		sql.Named("date", p.Date))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	id, err := res.LastInsertId()

	return int(id), err
}
