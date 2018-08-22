package api_test

import (
	"testing"
	"github.com/messwith/coding_challenge/api"
	"os"
	"log"
	"github.com/messwith/coding_challenge/utils"
	"github.com/messwith/coding_challenge/repository"
	"github.com/kelseyhightower/envconfig"
	"database/sql"
	"github.com/messwith/coding_challenge/models"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"net/http"
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"errors"
	"github.com/messwith/coding_challenge/service"
)

type config struct {
	DbConnString string `default:"postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"`
	Port int `default:"8080"`
}

var (
	cfg config
	db *sql.DB
	logger *log.Logger
)

func TestMain(m *testing.M) {

	logger = log.New(os.Stdout, "", log.LstdFlags)

	if err := envconfig.Process("", &cfg); err != nil {
		logger.Fatal(err)
	}

	var err error
	db, err = utils.DbConnect(cfg.DbConnString)
	if err != nil {
		log.Fatal(err)
	}

	accountRep := repository.NewAccountRepository(db)
	paymentRep := repository.NewPaymentRepository(db)

	accountService := service.NewSqlAccountService(accountRep)
	paymentService := service.NewSqlPaymentService(paymentRep)
	sqlTransactioner := service.NewSqlTransactioner(accountRep, paymentRep)

	server := api.NewServer(sqlTransactioner, accountService, paymentService, logger)
	go server.Run(cfg.Port)
	m.Run()
}

func createRandomAccount(id string) error {
	account := models.Account{
		ID:       id,
		Owner:    "test",
		Balance:  decimal.NewFromFloat(rand.Float64()*10),
		Currency: "PHP",
	}
	_, err := db.Exec(`INSERT INTO accounts (id, owner, balance, currency) VALUES (
					$1, $2, $3, $4	
					)`, account.ID, account.Owner, account.Balance, account.Currency)
	return err
}

func createAccount(id string, balance decimal.Decimal) error {
	account := models.Account{
		ID:       id,
		Owner:    "test",
		Balance:  balance,
		Currency: "PHP",
	}
	_, err := db.Exec(`INSERT INTO accounts (id, owner, balance, currency) VALUES (
					$1, $2, $3, $4	
					)`, account.ID, account.Owner, account.Balance, account.Currency)
	return err
}

func createRandomAccounts(count int) error {
	for i := 0; i < count; i++ {
		if err := createRandomAccount(strconv.Itoa(i+1)); err != nil {
			return err
		}
	}
	return nil
}

func clearDB() {
	logger.Println("startind clearing db")
	_, err := db.Exec(`DELETE FROM payments`)
	if err != nil {
		logger.Fatal(err)
	}
	_, err = db.Exec(`DELETE FROM accounts`)
	if err != nil {
		logger.Fatal(err)
	}
}

func request(method string, path string, payload interface{}, answer interface{}) error {
	var body []byte
	if payload != nil {
		js, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = js
	}

	req, _ := http.NewRequest(method, fmt.Sprintf("http://localhost:%d%s", cfg.Port, path), bytes.NewBuffer(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if answer != nil {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &answer); err != nil {
			return err
		}
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

