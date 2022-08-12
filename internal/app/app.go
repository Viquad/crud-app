package app

import (
	"log"
	"net/http"
	"time"

	"github.com/Viquad/crud-app/internal/repository/psql"
	"github.com/Viquad/crud-app/internal/service"
	"github.com/Viquad/crud-app/internal/transport/rest"
	"github.com/Viquad/crud-app/pkg/database"
)

func Run() {
	db, err := database.NewPostgresConnection(
		database.ConnectionInfo{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			DBName:   "postgres",
			SSLMode:  "disable",
			Password: "qwerty123",
		},
	)

	if err != nil {
		log.Fatalf("Failed connect to DB %s", err.Error())
	}

	defer db.Close()

	repo := psql.NewBank(db)
	service := service.NewBank(repo)
	handler := rest.NewHandler(service)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: rest.WithLog(handler.InitRouter()),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
