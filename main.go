package main

import (
	"final/api"
	"final/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db.CheckDB()
	var err error

	if db.Data, err = db.NewDB("scheduler.db"); err != nil {
		log.Panic(err)
	}

	defer db.Data.Close()

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.Handle("/api/nextdate", http.HandlerFunc(api.HandleDate))
	http.Handle("/api/task", http.HandlerFunc(rout))
	http.Handle("/api/tasks", http.HandlerFunc(api.GetTasks))
	http.Handle("/api/task/done", http.HandlerFunc(api.TaskDone))

	err = http.ListenAndServe(":7540", nil)

	if err != nil {
		panic(err)
	}

}

func rout(res http.ResponseWriter, req *http.Request) {
	method := fmt.Sprint(req.Method)

	if method == "GET" {
		api.GetTask(res, req)
	}
	if method == "POST" {
		api.PostTask(res, req)
	}
	if method == "PUT" {
		api.UpdateTask(res, req)
	}
	if method == "DELETE" {
		api.TaskDelete(res, req)
	}
}
