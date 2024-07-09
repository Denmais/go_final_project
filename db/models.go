package db

import (
	"database/sql"
	"final/nextdate"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Repeat  string `json:"repeat,omitempty"`
	Comment string `json:"comment"`
}
type Void struct {
}
type ErrorstrSer struct {
	Error string `json:"error"`
}
type Tasks struct {
	Tasks []Task `json:"tasks"`
}
type IdSer struct {
	ID int `json:"id"`
}
type SchedulerStore struct {
	db *sql.DB
}

func (task *Task) CheckTask() error {
	now := time.Now().Format(nextdate.Date)
	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse(nextdate.Date, task.Date)
	if err != nil {
		return fmt.Errorf("Неправильный формат даты")
	}
	if task.Title == "" {
		return fmt.Errorf("Не указан заголовок")
	}
	if task.Date < now && task.Repeat == "" {
		task.Date = now
	} else if task.Date < now && task.Repeat != "" {
		task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("Неправильный формат даты")
		}

	}
	_, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("ошибка в данных")
	}
	return nil
}
