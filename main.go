package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/valentyn88/funds-loader/model"
	"github.com/valentyn88/funds-loader/service"
	"github.com/valentyn88/funds-loader/storage"
)

func main() {
	logErr := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)

	f, err := os.Open("testdata/input.txt")
	if err != nil {
		logErr.Printf("couldn't open a file %s\n", err)
	}
	defer f.Close()

	storage := storage.NewStorage(sync.RWMutex{}, map[string]storage.CustomerFunds{})
	fundsLoader := service.NewFundsLoader(service.WeeklyAmountLimit,
		service.DailyAmountLimit,
		service.AttemptsLimit,
		storage)

	fRes, err := os.Create("testdata/result.txt")
	if err != nil {
		logErr.Printf("couldn't create a file with results %s\n", err)
	}
	defer fRes.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		customerFunds := model.LoadedFunds{}
		if err := json.Unmarshal(scanner.Bytes(), &customerFunds); err != nil {
			logErr.Printf("couldn't Unmarshal customer funds %s\n", err)
			continue
		}

		transaction := fundsLoader.Load(customerFunds)
		bResp, err := json.Marshal(transaction)
		if err = json.Unmarshal(scanner.Bytes(), &customerFunds); err != nil {
			logErr.Printf("couldn't Marshal response %s\n", err)
			continue
		}

		if _, err := fRes.WriteString(string(bResp) + "\n"); err != nil {
			logErr.Printf("couldn't write to the result file %s\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		logErr.Printf("an error occured during reading a file %s\n", err)
	}
}
