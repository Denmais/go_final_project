package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

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

func HandleDate(res http.ResponseWriter, req *http.Request) {
	date := req.URL.Query().Get("date")
	now := req.URL.Query().Get("now")
	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		res.Write([]byte(err.Error()))
	}
	repeat := req.URL.Query().Get("repeat")
	lastDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		res.Write([]byte(err.Error()))
	} else {
		res.Write([]byte(lastDate))
	}
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res, errs := SheldureTaskDB(task, w)
	if errs != "" {
		resp, _ := json.Marshal(ErrorstrSer{errs})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	Serializer := IdSer{ID: int(res)}
	newres, _ := json.Marshal(Serializer)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(newres)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	res, err := SheldureGet()
	if err != nil {
		if res != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		res, _ := json.Marshal(ErrorstrSer{"Не указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	idint, err := strconv.Atoi(id)

	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Неправильно указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	defer db.Close()
	var task Task2
	row := db.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", idint))
	err = row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	res, err := json.Marshal(task)
	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task Task2
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	_, errs := SheldureUpdate(task, w)
	if errs != "" {
		resp, _ := json.Marshal(ErrorstrSer{errs})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	newres, _ := json.Marshal(Void{})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(newres)
}

func TaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		res, _ := json.Marshal(ErrorstrSer{"Не указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	idint, err := strconv.Atoi(id)

	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Неправильно указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	defer db.Close()
	var task Task2
	row := db.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", idint))
	err = row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	if task.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", idint))
		if err != nil {
			res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
			return
		}
		res, _ := json.Marshal(Void{})
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	} else {
		str, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			res, _ := json.Marshal(ErrorstrSer{str})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
			return
		}
		db.Exec("UPDATE scheduler SET title = :title, repeat = :repeat, date = :date, comment = :comment WHERE id = :id",
			sql.Named("title", task.Title),
			sql.Named("repeat", task.Repeat),
			sql.Named("date", str),
			sql.Named("comment", task.Comment),
			sql.Named("id", idint))

		res, _ := json.Marshal(Void{})
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		res, _ := json.Marshal(ErrorstrSer{"Не указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	idint, err := strconv.Atoi(id)

	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Неправильно указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	defer db.Close()
	var task Task2
	row := db.QueryRow("SELECT id, comment, repeat, title, date FROM scheduler WHERE id = :id", sql.Named("id", idint))
	err = row.Scan(&task.ID, &task.Comment, &task.Repeat, &task.Title, &task.Date)
	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", idint))
	if err != nil {
		res, _ := json.Marshal(ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	res, _ := json.Marshal(Void{})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}
