package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/Kaikawa1028/go-template/app/logger"
	"github.com/Kaikawa1028/go-template/app/wire"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	godotenv.Load()

	if err := logger.SetupLogger(); err != nil {
		log.Fatalf("Logger initialization failed: %+v", err)
		return
	}

	di, cleanup, err := wire.InitializeDIContainer()
	if err != nil {
		logger.SimpleFatal(err, nil)
		return
	}
	defer cleanup()

	di.Router.Attach(di.Server.Echo)

	go func() {
		err = di.Server.Start()
		if err != nil {
			logger.SimpleFatal(err, nil)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := di.Server.Echo.Shutdown(ctx); err != nil {
		logger.SimpleFatal(err, nil)
	}
}
