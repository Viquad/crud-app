package app

import (
	"log"
	"net/http"
	"time"

	"github.com/Viquad/crud-app/internal/repository/psql"
	"github.com/Viquad/crud-app/internal/service"
	"github.com/Viquad/crud-app/internal/transport/rest"
	"github.com/Viquad/crud-app/pkg/config"
	"github.com/Viquad/crud-app/pkg/database"
)

func Run() {
	cfg, err := config.New("configs", "config")
	if err != nil {
		log.Fatalf("can't initialize config: %s", err.Error())
	}

	db, err := database.NewPostgresConnection(cfg.DB)
	if err != nil {
		log.Fatalf("can't connect to DB %s", err.Error())
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
