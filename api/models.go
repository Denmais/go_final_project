package api

import "database/sql"

type Task struct {
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Repeat  string `json:"repeat,omitempty"`
	Comment string `json:"comment"`
}
type Task2 struct {
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
type SliceSheldure struct {
	Tasks []Task2 `json:"tasks"`
}
type IdSer struct {
	ID int `json:"id"`
}
type SchedulerStore struct {
	db *sql.DB
}
