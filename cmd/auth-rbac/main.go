package main

import (
	"auth-rbac/internal/config"
	"auth-rbac/internal/postgres"
	"auth-rbac/internal/rbac"
	"auth-rbac/pkg/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	cfg, err := config.LoadConfiguration("./config/config.json")
	if err != nil {
		log.Fatalf("Could not instantiate config %+s", err)
	}

	newLogger, err := logger.NewLogger(cfg.Log)
	if err != nil {
		log.Fatalf("Could not instantiate log %+s", err)
	}

	db := postgres.New(newLogger, cfg)

	defer db.Close()

	user, err := postgres.NewUserStorage(db)
	if err != nil {
		newLogger.Fatalf("Could not instantiate user statements %+s", err)
	}

	session, err := postgres.NewSessionStorage(db)
	if err != nil {
		newLogger.Fatalf("Could not instantiate session statements %+s", err)
	}

	rbac, err := rbac.FromFile("./config/rbac.yml")
	if err != nil {
		newLogger.Fatalf("Could not instantiate permissons %+s", err)
	}

	templates := parseTemplates()

	handler := newHandler(newLogger, cfg, rbac, user, session, templates)

	r := chi.NewRouter()

	handler.routers(r)

	srv := &http.Server{
		Addr:    cfg.Http.Adress,
		Handler: r,
	}

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			newLogger.Fatalf("server stopped %+s", err)
		}
	}()

	newLogger.Debugf("server started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		newLogger.Debugf("system call:%+v", oscall)
		cancel()
	}()

	<-ctx.Done()

	newLogger.Debugf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cancel()
	}()

	err = srv.Shutdown(ctxShutDown)
	if err != nil {
		newLogger.Fatalf("server shutdown failed:%+s", err)
	}

	newLogger.Debugf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}
