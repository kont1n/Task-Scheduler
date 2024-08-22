package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"Task-Scheduler/config"
	"Task-Scheduler/database"
	"Task-Scheduler/internal"
)

func main() {

	db, err := database.CheckDB()
	if err != nil {
		fmt.Println("Error checking database connection:", err)
		return
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	store := internal.NewStorage(db)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("web")))
	r.Get("/api/nextdate", internal.NextDateHandle)
	r.Get("/api/tasks", internal.GetTasksHandle(store))
	r.Post("/api/task", internal.CreateTaskHandle(store))
	r.Get("/api/task", internal.ReadTaskHandle(store))
	r.Put("/api/task", internal.UpdateTaskHandle(store))
	r.Delete("/api/task", internal.DeleteTaskHandle(store))
	r.Post("/api/task/done", internal.DoneTaskHandle(store))

	serverPort := config.GetServerPort()
	fmt.Println("Server started on localhost port", serverPort)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", serverPort), r); err != nil {
		fmt.Println("Web server error:", err)
		return
	}
}
