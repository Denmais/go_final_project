package db

import (
	"database/sql"
	"final/nextdate"
	"fmt"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

func (db DB) AddTask(task Task) (int, error) {

	err := task.CheckTask()
	if err != nil {
		return -1, err
	}

	res, err := db.database.Exec("INSERT INTO scheduler (title, comment, repeat, date) VALUES (:title, :comment, :repeat, :date)",
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("date", task.Date))

	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()

	return int(id), nil
}

func (db DB) GetTasks() (Tasks, error) {

	tasks := Tasks{}
	sliceSheldure := []Task{}
	rows, err := db.database.Query("SELECT id, date, comment, repeat, title FROM scheduler")

	if err != nil {
		return Tasks{}, err
	}

	defer rows.Close()

	for rows.Next() {
		task := Task{}
		var id int
		err := rows.Scan(&id, &task.Date, &task.Comment, &task.Repeat, &task.Title)
		if err != nil {
			return Tasks{}, err
		}
		task.ID = fmt.Sprint(id)
		sliceSheldure = append(sliceSheldure, task)
	}

	tasks.Tasks = sliceSheldure
	return tasks, nil
}

func (db DB) TaskUpdate(task Task) error {
	err := task.CheckTask()

	if err != nil {
		return err
	}

	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}

	row := db.database.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&id)
	if err != nil {
		return err
	}

	db.database.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("title", task.Title),
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", id))

	return nil

}

func (db DB) TaskDone(id int) error {

	var task Task
	row := db.database.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		_, err := db.database.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
		if err != nil {
			return err
		}
		return nil
	} else {
		str, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}
		db.database.Exec("UPDATE scheduler SET title = :title, repeat = :repeat, date = :date, comment = :comment WHERE id = :id",
			sql.Named("title", task.Title),
			sql.Named("repeat", task.Repeat),
			sql.Named("date", str),
			sql.Named("comment", task.Comment),
			sql.Named("id", id))

		return nil
	}
}
func (db DB) Delete(id int) error {
	var task Task
	row := db.database.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		return err
	}

	_, err = db.database.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetTask(id int) (Task, error) {
	var task Task
	row := db.database.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		return Task{}, err
	}
	return task, nil

}
