package db

import (
	"database/sql"
	"final/nextdate"
	"fmt"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

func NewStore(db *sql.DB) SchedulerStore {
	return SchedulerStore{db: db}
}
func (db DB) Add(p Task) string {

	res, err := db.db.Exec("INSERT INTO scheduler (title, comment, repeat, date) VALUES (:title, :comment, :repeat, :date)",
		sql.Named("title", p.Title),
		sql.Named("comment", p.Comment),
		sql.Named("repeat", p.Repeat),
		sql.Named("date", p.Date))
	if err != nil {
		return "error"
	}
	id, _ := res.LastInsertId()
	idstr := strconv.FormatInt(id, 10)

	return idstr
}

func (task *Task) CheckTaskDB() string {
	now := time.Now().Format(nextdate.Date)
	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse(nextdate.Date, task.Date)
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

func (db DB) PostTaskDB(task Task) (int, string) {

	err := task.CheckTaskDB()
	if err != "" {
		return -1, err
	}

	newres := db.Add(task)
	if newres == "error" {
		return -1, "ошибка в данных"
	}
	ans, _ := strconv.Atoi(newres)

	return ans, ""
}

func (db DB) GetTasksDB() (SliceSheldure, error) {

	tasks := SliceSheldure{}
	sliceSheldure := []Task{}
	rows, err := db.db.Query("SELECT * FROM scheduler")
	if err != nil {
		return SliceSheldure{}, err
	}
	defer rows.Close()
	for rows.Next() {
		task := Task{}
		var id int
		err := rows.Scan(&id, &task.Title, &task.Comment, &task.Repeat, &task.Date)
		if err != nil {
			return SliceSheldure{}, err
		}
		task.ID = fmt.Sprint(id)
		sliceSheldure = append(sliceSheldure, task)
	}
	tasks.Tasks = sliceSheldure
	return tasks, nil
}

func (db DB) TaskUpdateDB(task Task) (int, string) {
	err := task.CheckTaskDB()

	if err != "" {
		return -1, err
	}
	errnew := db.Update(task)
	if errnew != nil {
		return -1, errnew.Error()
	}
	return 0, ""

}

func (db DB) Update(p Task) error {

	id, err := strconv.Atoi(p.ID)
	if err != nil {
		return err
	}
	row := db.db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&id)
	if err != nil {
		return err
	}
	db.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("title", p.Title),
		sql.Named("date", p.Date),
		sql.Named("comment", p.Comment),
		sql.Named("repeat", p.Repeat),
		sql.Named("id", id))

	return nil
}

func (db DB) TaskDoneDB(id int) string {

	var task Task
	row := db.db.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		return err.Error()
	}
	if task.Repeat == "" {
		_, err := db.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
		if err != nil {
			return err.Error()
		}
		return ""
	} else {
		str, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err.Error()
		}
		db.db.Exec("UPDATE scheduler SET title = :title, repeat = :repeat, date = :date, comment = :comment WHERE id = :id",
			sql.Named("title", task.Title),
			sql.Named("repeat", task.Repeat),
			sql.Named("date", str),
			sql.Named("comment", task.Comment),
			sql.Named("id", id))

		return ""
	}
}
func (db DB) TaskDeleteDB(id int) string {
	var task Task
	row := db.db.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		return err.Error()
	}
	_, err = db.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err.Error()
	}
	return ""
}

func (db DB) GetTaskDB(id int) (Task, string) {
	var task Task
	row := db.db.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		return Task{}, err.Error()
	}
	return task, ""

}
