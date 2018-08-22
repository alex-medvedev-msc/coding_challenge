package main

import (
	"github.com/messwith/coding_challenge/api"
	"github.com/messwith/coding_challenge/utils"
	"log"
	"github.com/messwith/coding_challenge/repository"
	"os"
	"github.com/kelseyhightower/envconfig"
	"github.com/messwith/coding_challenge/service"
)

type config struct {
	DbConnString string `default:"postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"`
	Port int `default:"8080"`
}

func parseConfig() (config, error) {
	cfg := config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func main() {
	// we are using stdout because deploying will be done via docker
	logger := log.New(os.Stdout, "", log.LstdFlags)

	cfg, err := parseConfig()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := utils.DbConnect(cfg.DbConnString)
	if err != nil {
		log.Fatal(err)
	}

	accountRep := repository.NewAccountRepository(db)
	paymentRep := repository.NewPaymentRepository(db)
	accountService := service.NewSqlAccountService(accountRep)
	paymentService := service.NewSqlPaymentService(paymentRep)
	sqlTransactioner := service.NewSqlTransactioner(accountRep, paymentRep)

	server := api.NewServer(sqlTransactioner, accountService, paymentService, logger)
	log.Fatal(server.Run(cfg.Port))
}
