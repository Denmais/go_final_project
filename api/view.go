package api

import (
	"bytes"
	"encoding/json"
	"final/db"
	"final/nextdate"
	"net/http"
	"strconv"
	"time"
)

func HandleDate(res http.ResponseWriter, req *http.Request) {
	date := req.URL.Query().Get("date")
	now := req.URL.Query().Get("now")
	nowTime, err := time.Parse(nextdate.Date, now)

	if err != nil {
		res.Write([]byte(err.Error()))
	}

	repeat := req.URL.Query().Get("repeat")
	lastDate, err := nextdate.NextDate(nowTime, date, repeat)

	if err != nil {
		res.Write([]byte(err.Error()))
	} else {
		res.Write([]byte(lastDate))
	}

}

func PostTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res, err := db.Data.AddTask(task)

	if err != nil {
		resp, _ := json.Marshal(db.ErrorstrSer{err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	serializer := db.IdSer{ID: int(res)}
	newres, err := json.Marshal(serializer)

	if err != nil {
		resp, _ := json.Marshal(db.ErrorstrSer{err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(newres)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	res, err := db.Data.GetTasks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(res)

	if err != nil {
		resp, _ := json.Marshal(db.ErrorstrSer{err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		res, _ := json.Marshal(db.ErrorstrSer{"Не указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	idint, err := strconv.Atoi(id)

	if err != nil {
		res, _ := json.Marshal(db.ErrorstrSer{"Неправильно указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	res, err := db.Data.GetTask(idint)

	if err != nil {
		resp, _ := json.Marshal(db.ErrorstrSer{err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	ans, err := json.Marshal(res)

	if err != nil {
		resp, _ := json.Marshal(db.ErrorstrSer{err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ans)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		res, _ := json.Marshal(db.ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		res, _ := json.Marshal(db.ErrorstrSer{"Ошибка сервера"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	err = db.Data.TaskUpdate(task)
	if err != nil {
		resp, _ := json.Marshal(db.ErrorstrSer{err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	newres, err := json.Marshal(db.Void{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(newres)
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		res, _ := json.Marshal(db.ErrorstrSer{"Не указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	idint, err := strconv.Atoi(id)

	if err != nil {
		res, _ := json.Marshal(db.ErrorstrSer{"Неправильно указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	err = db.Data.Delete(idint)

	if err != nil {
		res, _ := json.Marshal(db.ErrorstrSer{err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
	}

	res, err := json.Marshal(db.Void{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return

}

func TaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		res, _ := json.Marshal(db.ErrorstrSer{"Не указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	idint, err := strconv.Atoi(id)

	if err != nil {
		res, _ := json.Marshal(db.ErrorstrSer{"Неправильно указан идентификатор"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	err = db.Data.TaskDone(idint)
	if err != nil {
		res, _ := json.Marshal(db.ErrorstrSer{err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
	}

	res, _ := json.Marshal(db.Void{})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return

}
