package main

import (
	"github.com/ChristinaFomenko/gophermart/pkg/logger"
)

func main() {
	log, err := logger.InitLogger()
	if err != nil {
		log.Fatal("error init logger")
	}

	defer log.Sync()
	zp := log.Sugar()

}
