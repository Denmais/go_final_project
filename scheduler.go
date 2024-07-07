package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

type SchedulerStore struct {
	db *sql.DB
}

func NewSchedulerStore(db *sql.DB) SchedulerStore {
	return SchedulerStore{db: db}
}

func Add(p Task, w http.ResponseWriter) []byte {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return []byte("error")
	}
	defer db.Close()
	res, err := db.Exec("INSERT INTO scheduler (title, comment, repeat, date) VALUES (:title, :comment, :repeat, :date)",
		sql.Named("title", p.Title),
		sql.Named("comment", p.Comment),
		sql.Named("repeat", p.Repeat),
		sql.Named("date", p.Date))
	if err != nil {
		return []byte("error")
	}
	id, _ := res.LastInsertId()

	return []byte(strconv.FormatInt(int64(id), 10))
}

func SheldureTaskDB(task Task, w http.ResponseWriter) (int, string) {
	now := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return -1, "Неправильный формат даты"
	}
	if task.Title == "" {
		return -1, "Не указан заголовок"
	}
	if task.Date < now && task.Repeat == "" {
		task.Date = now
	} else if task.Date < now && task.Repeat != "" {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return -1, "Неправильный формат даты"
		}

	}
	_, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return -1, "ошибка в данных"
	}
	newres := Add(task, w)
	if string(newres) == "error" {
		return -1, "ошибка в данных"
	}
	ans, _ := strconv.Atoi(string(newres))

	return ans, ""
}

func SheldureGet() ([]byte, error) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return []byte{}, err
	}
	defer db.Close()
	tasks := SliceSheldure{}
	sliceSheldure := []Task2{}
	rows, err := db.Query("SELECT * FROM scheduler")
	if err != nil {
		return []byte{}, err
	}
	defer rows.Close()
	for rows.Next() {
		task := Task2{}
		var id int
		err := rows.Scan(&id, &task.Title, &task.Comment, &task.Repeat, &task.Date)
		if err != nil {
			return []byte{}, err
		}
		task.ID = fmt.Sprint(id)
		sliceSheldure = append(sliceSheldure, task)
	}
	tasks.Tasks = sliceSheldure
	resp, err := json.Marshal(tasks)
	if err != nil {
		return []byte{}, err
	}
	return resp, nil
}

func SheldureUpdate(task Task2, w http.ResponseWriter) (int, string) {
	now := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return -1, "Неправильный формат даты"
	}
	if task.Title == "" {
		return -1, "Не указан заголовок"
	}
	if task.Date < now && task.Repeat == "" {
		task.Date = now
	} else if task.Date < now && task.Repeat != "" {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return -1, "Неправильный формат даты"
		}

	}
	_, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return -1, "ошибка в данных"
	}
	newres := Update(task, w)
	if string(newres) == "error" {
		return -1, "ошибка в данных"
	}
	return 0, ""

}

func Update(p Task2, w http.ResponseWriter) []byte {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return []byte("error")
	}
	defer db.Close()
	id, err := strconv.Atoi(p.ID)
	if err != nil {
		return []byte("error")
	}
	row := db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&id)
	if err != nil {
		return []byte("error")
	}
	db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("title", p.Title),
		sql.Named("date", p.Date),
		sql.Named("comment", p.Comment),
		sql.Named("repeat", p.Repeat),
		sql.Named("id", id))

	return []byte{}
}
