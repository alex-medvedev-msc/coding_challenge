package main

import (
	"github.com/messwith/coding_challenge/api"
	"github.com/messwith/coding_challenge/utils"
	"log"
	"github.com/messwith/coding_challenge/repository"
	"os"
)

type config struct {
	DbConnString string
	Port int
}

func parseConfig() config {
	return config{}
}

func main() {
	cfg := parseConfig()

	// we are using stdout because deploying will be done via docker
	logger := log.New(os.Stdout, "", log.LstdFlags)

	db, err := utils.DbConnect(cfg.DbConnString)
	if err != nil {
		log.Fatal(err)
	}

	accountRep := repository.NewAccountRepository(db)
	paymentRep := repository.NewPaymentRepository(db)

	server := api.NewServer(accountRep, paymentRep, logger)
	log.Fatal(server.Run(cfg.Port))
}
