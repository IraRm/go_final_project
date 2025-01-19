package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go_final_project/database"
	"go_final_project/server"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() *sql.DB {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	if install {
		if database.CreateTableSQL(db) != nil {
			log.Fatalf("Failed to create table: %v", err)
		}

		if database.CreateIndexSQL(db) != nil {
			log.Fatalf("Failed to create index: %v", err)
		}
	}
	return db
}

func main() {
	db := initDB()

	r := chi.NewRouter()
	server := server.NewServer(db)

	fs := http.FileServer(http.Dir("./web/"))
	r.Handle("/*", http.StripPrefix("/", fs))
	r.Get("/api/tasks", server.GetTasks)
	r.Get("/api/task", server.GetTask)
	r.Get("/api/nextdate", server.GetNextDate)
	r.Post("/api/task", server.PostTask)
	r.Post("/api/task/done", server.PostTaskDone)
	r.Put("/api/task", server.PutTask)
	r.Delete("/api/task", server.DeleteTask)

	fmt.Println("Starting server at :7540")
	if err := http.ListenAndServe(":7540", r); err != nil {
		fmt.Printf("Server failed to start: %s\n", err)
	}
}
