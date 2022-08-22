package main

import (
	"context"
	"github.com/ChristinaFomenko/gophermart/config"
	app "github.com/ChristinaFomenko/gophermart/internal/app"
	brokerServ "github.com/ChristinaFomenko/gophermart/internal/broker/service"
	handler "github.com/ChristinaFomenko/gophermart/internal/controller/http/handlers"
	repository "github.com/ChristinaFomenko/gophermart/internal/repository"
	psql "github.com/ChristinaFomenko/gophermart/internal/repository/psql"
	"github.com/ChristinaFomenko/gophermart/internal/service"
	"github.com/ChristinaFomenko/gophermart/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log, err := logger.InitLogger()
	if err != nil {
		log.Fatal("error init logger")
	}

	defer log.Sync()
	zp := log.Sugar()

	conf, err := config.NewConfig()
	if err != nil {
		zp.Fatalf("failed to retrieve env variables, %v", err)
	}

	db, err := psql.NewPsql(conf.DatabaseURI)
	if err != nil {
		zp.Fatalf("DB connection error %v", err)
	}
	err = db.Init()
	if err != nil {
		zp.Fatalf("failed to create db table %v", err)
	}

	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()

	repos := repository.NewRepository(db.DB, log)
	services := service.NewService(nil, repos, log)
	handlers := handler.NewHandler(services, log)

	brokerRepos := repository.NewBrokerRepository(db.DB, log)
	broker := brokerServ.NewBroker(brokerRepos, conf.AccrualSystemAddress, log)
	broker.Start(ctx)

	server := app.NewServer(conf, handlers.InitRoutes())

	//graceful shutdown
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-termChan
		zp.Infof("server shutdown success")
		cansel()
		if err = server.Stop(ctx); err != nil {
			zp.Fatalf("server shutdown error %v", err)
		}
	}()

	if err = server.Run(); err != nil && err != http.ErrServerClosed {
		zp.Fatalf("server run error %v", err)
	}

}
