package main

import (
	"fmt"
	"net/http"
)

type Scheduler struct {
	ID      int
	Title   string
	Comment string
	Repeat  string
	Date    string
}

type SchedulerService struct {
	store SchedulerStore
}

func NewSchedulerService(store SchedulerStore) SchedulerService {
	return SchedulerService{store: store}
}

const Value = 1048576

func main() {
	CheckDB()
	//r := chi.NewRouter()
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.Handle("/api/nextdate", http.HandlerFunc(HandleDate))
	http.Handle("/api/task", http.HandlerFunc(rout))
	http.Handle("/api/tasks", http.HandlerFunc(GetTasks))
	http.Handle("/api/task/done", http.HandlerFunc(TaskDone))

	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}

}

func rout(res http.ResponseWriter, req *http.Request) {
	method := fmt.Sprint(req.Method)
	if method == "GET" {
		GetTask(res, req)
	}
	if method == "POST" {
		PostTask(res, req)
	}
	if method == "PUT" {
		UpdateTask(res, req)
	}
	if method == "DELETE" {
		TaskDelete(res, req)
	}
}
