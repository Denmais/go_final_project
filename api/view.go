package api

import (
	"bytes"
	"encoding/json"
	"final/nextdate"
	"net/http"
	"strconv"
	"time"
)

func HandleDate(res http.ResponseWriter, req *http.Request) {
	date := req.URL.Query().Get("date")
	now := req.URL.Query().Get("now")
	nowTime, err := time.Parse("20060102", now)
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
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res, errs := PostTaskDB(task)
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
	res, err := GetTasksDB()
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
	res, errs := GetTaskDB(idint)
	if errs != "" {
		resp, _ := json.Marshal(ErrorstrSer{errs})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	ans, err := json.Marshal(res)
	if err != nil {
		resp, _ := json.Marshal(ErrorstrSer{errs})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ans)
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

	_, errs := TaskUpdateDB(task)
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
	errs := TaskDeleteDB(idint)
	if errs != "" {
		res, _ := json.Marshal(ErrorstrSer{errs})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
	}

	res, _ := json.Marshal(Void{})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return

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
	errs := TaskDoneDB(idint)
	if errs != "" {
		res, _ := json.Marshal(ErrorstrSer{errs})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
	}

	res, _ := json.Marshal(Void{})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return

}
