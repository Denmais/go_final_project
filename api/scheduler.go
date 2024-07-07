package api

import (
	"database/sql"
	"encoding/json"
	"final/nextdate"
	"fmt"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

func NewStore(db *sql.DB) SchedulerStore {
	return SchedulerStore{db: db}
}
func (s SchedulerStore) Add(p Task) []byte {

	res, err := s.db.Exec("INSERT INTO scheduler (title, comment, repeat, date) VALUES (:title, :comment, :repeat, :date)",
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

func (task *Task) CheckTaskDB() string {
	now := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return "Неправильный формат даты"
	}
	if task.Title == "" {
		return "Не указан заголовок"
	}
	if task.Date < now && task.Repeat == "" {
		task.Date = now
	} else if task.Date < now && task.Repeat != "" {
		task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return "Неправильный формат даты"
		}

	}
	_, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return "ошибка в данных"
	}
	return ""
}
func (task *Task2) CheckTaskDB2() string {
	now := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return "Неправильный формат даты"
	}
	if task.Title == "" {
		return "Не указан заголовок"
	}
	if task.Date < now && task.Repeat == "" {
		task.Date = now
	} else if task.Date < now && task.Repeat != "" {
		task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return "Неправильный формат даты"
		}

	}
	_, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return "ошибка в данных"
	}
	return ""
}
func PostTaskDB(task Task) (int, string) {

	err := task.CheckTaskDB()
	if err != "" {
		return -1, err
	}

	db, errs := sql.Open("sqlite", "scheduler.db")
	if errs != nil {
		return -1, errs.Error()
	}
	defer db.Close()
	store := NewStore(db)
	newres := store.Add(task)
	if string(newres) == "error" {
		return -1, "ошибка в данных"
	}
	ans, _ := strconv.Atoi(string(newres))

	return ans, ""
}

func GetTasksDB() ([]byte, error) {
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

func TaskUpdateDB(task Task2) (int, string) {
	err := task.CheckTaskDB2()

	if err != "" {
		return -1, err
	}
	newres := Update(task)
	if string(newres) == "error" {
		return -1, "ошибка в данных"
	}
	return 0, ""

}

func Update(p Task2) []byte {
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

func TaskDoneDB(id int) string {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err.Error()
	}
	defer db.Close()
	var task Task2
	row := db.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		return err.Error()
	}
	if task.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
		if err != nil {
			return err.Error()
		}
		return ""
	} else {
		str, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err.Error()
		}
		db.Exec("UPDATE scheduler SET title = :title, repeat = :repeat, date = :date, comment = :comment WHERE id = :id",
			sql.Named("title", task.Title),
			sql.Named("repeat", task.Repeat),
			sql.Named("date", str),
			sql.Named("comment", task.Comment),
			sql.Named("id", id))

		return ""
	}
}
func TaskDeleteDB(id int) string {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err.Error()
	}
	defer db.Close()
	var task Task2
	row := db.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		return err.Error()
	}
	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err.Error()
	}
	return ""
}

func GetTaskDB(id int) (Task2, string) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return Task2{}, err.Error()
	}
	defer db.Close()
	var task Task2
	row := db.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		return Task2{}, err.Error()
	}
	return task, ""

}
