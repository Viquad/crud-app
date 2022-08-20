package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Viquad/crud-app/internal/repository/psql"
	"github.com/Viquad/crud-app/internal/service"
	"github.com/Viquad/crud-app/internal/transport/rest"
	"github.com/Viquad/crud-app/pkg/config"
	"github.com/Viquad/crud-app/pkg/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// @title       CRUD app API
// @version     1.0
// @description API Server for CRUD application

// @host     localhost:8080
// @BasePath /

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.InfoLevel)
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	cfg, err := config.New("configs", "config")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context": "app.Run()",
			"problem": "can't initialize config",
		}).Fatal(err.Error())
	}

	db, err := database.NewPostgresConnection(cfg.DB)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context": "app.Run()",
			"problem": "can't connect to DB",
		}).Fatal(err.Error())
	}

	defer db.Close()

	repo := psql.NewRepositories(db)
	services := service.NewServices(repo)
	handler := rest.NewHandler(services)

	router := handler.InitRouter()

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logrus.WithFields(logrus.Fields{
			"context": "app.Run()",
			"problem": "server shutdowned",
		}).Error(err.Error())
	}
}
