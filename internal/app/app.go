package app

import (
	"net/http"
	"task-app/internal/controller"
	"task-app/internal/repository"
	"task-app/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Run() {
	db, err := sqlx.Open("sqlite3", "dev.db")
	if err != nil {
		return
	}
	defer db.Close()

	repo := repository.New(db)
	uc := usecase.New(repo)
	s := controller.NewServer(uc)
	router := chi.NewRouter()

	handler := controller.NewHandler(s, router)
	http.ListenAndServe(":8000", handler)
}
