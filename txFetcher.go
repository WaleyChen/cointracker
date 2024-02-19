package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const apiKey = "f45402fbaf2f7649ce26a55aa9ab555033c48b7f"

var client = &http.Client{}

const testing = true
const testTotalNumTxs = 20

const limitNumResultsInProd = true
const limitNumResults = 10

// const testTotalNumTxs = 1305156

const pageLimit = 1
const maxRequestsPerSecond = 10

type TxFetcher struct {
	Address     string
	CurrentPage int
	db          *sql.DB

	TotalNumPages int // The total number of pages for the address
	TotalNumTxs   int // The total number of txs for the address

	TotalSyncPages int // The number of pages we'll attempt to sync
	TotalSyncTxs   int // The number of txs we'll attempt to sync
	TotalTxsSynced int // The total number of txs we've already synced
}

func GetNewTxFetcher(db *sql.DB, address string, pageLimit int) TxFetcher {
	totalNumTxs, _ := getTotalNumberOfTxs(address)
	totalNumPages := int(math.Ceil(float64(totalNumTxs) / float64(pageLimit)))

	// Record the address
	now := time.Now().UTC()
	stmt := `INSERT INTO addresses (address, created_at) VALUES ($1, $2) ON CONFLICT (address) DO NOTHING`
	_, err := db.Exec(stmt, address, now)
	if err != nil {
		log.Fatal(err)
	}

	// Get the total number of txs we've already synced
	query := `SELECT COALESCE(SUM(txs_synced), 0) FROM syncs`
	var totalTxsSynced int // The total number of txs we've already synced
	err = db.QueryRow(query).Scan(&totalTxsSynced)
	if err != nil {
		log.Fatal(err)
	}

	totalSyncPages := (totalNumTxs - totalTxsSynced + pageLimit - 1) / pageLimit

	fmt.Println("Total Number of Address's Txs: ", totalNumTxs)
	fmt.Println("Total Number of Addresses Synced: ", totalTxsSynced)
	fmt.Println("Total Number of Address's Pages: ", totalNumPages)
	fmt.Println("Total Number of Pages to Sync: ", totalSyncPages)

	return TxFetcher{
		Address:        address,
		db:             db,
		CurrentPage:    0,
		TotalNumPages:  totalNumPages,
		TotalNumTxs:    totalNumTxs,
		TotalSyncPages: totalSyncPages,
		TotalSyncTxs:   totalNumTxs - totalTxsSynced,
		TotalTxsSynced: totalTxsSynced,
	}
}

func (txFetcher TxFetcher) SyncTxs() {
	fmt.Println("Syncing txs for address: ", txFetcher.Address)

	// Record the sync txs attempt
	now := time.Now().UTC()
	stmt := `INSERT INTO syncs (address, status, txs_synced, total_sync_txs, total_txs, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var syncId string
	err := txFetcher.db.QueryRow(stmt, txFetcher.Address, "STARTED", 0, txFetcher.TotalSyncTxs, txFetcher.TotalNumTxs, now).Scan(&syncId)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second / time.Duration(maxRequestsPerSecond))
	defer ticker.Stop()

	var wg sync.WaitGroup
	for txFetcher.CurrentPage < txFetcher.TotalSyncPages {
		select {
		case t := <-ticker.C:
			wg.Add(1)

			// Send a request
			go func(page int) {
				defer wg.Done()

				err := txFetcher.worker(txFetcher.db, page, int64(t.UnixMilli()), syncId)
				if err != nil {
					// TODO: retry 2 more times with exponential backoff
					// TODO: store error into sync error table
					log.Fatal(err)
				}
			}(txFetcher.CurrentPage)
			txFetcher.CurrentPage++
		}
	}

	wg.Wait() // Wait for all workers to finish

	// Get the number of txs synced
	query := `SELECT COUNT(*) FROM txs WHERE sync_id = $1`
	var count int
	err = txFetcher.db.QueryRow(query, syncId).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	// Update the sync txs status to COMPLETED
	now = time.Now().UTC()
	stmt = `UPDATE syncs SET STATUS = $1, FINISHED_AT = $2, TXS_SYNCED = $3 WHERE id = $4`
	_, err = txFetcher.db.Exec(stmt, "COMPLETED", now, count, syncId)
	if err != nil {
		log.Fatal(err)
	}
}

func (txFetcher TxFetcher) worker(db *sql.DB, page int, timestamp int64, sync_id string) error {
	fmt.Printf("Making request for page %d at timestamp %d\n", page, timestamp)

	var body []byte
	var response APIResponse
	var txBytes []byte
	if !testing {
		url := fmt.Sprintf(
			"https://rest.cryptoapis.io/blockchain-data/bitcoin/mainnet/addresses/%s/transactions?context=yourExampleString&limit=%d&offset=%d",
			address,
			pageLimit,
			page*pageLimit,
		)
		method := "GET"

		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return errors.New("failed to create request")
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", apiKey)

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return errors.New("failed to send request")
		}
		defer res.Body.Close()

		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return errors.New("failed to read response body")
		}
	} else {
		// Simulate a request that takes between 0.5 and 2 seconds
		max := 2000
		min := 500
		reqTimeInMs := rand.Intn(max-min) + min
		time.Sleep(time.Duration(reqTimeInMs) * time.Millisecond)

		body = []byte(fakeResponse)
	}

	// Unmarshal the JSON data into the struct
	err := json.Unmarshal(body, &response)
	if err != nil {
		return errors.New("failed to unmarshal JSON")
	}

	if testing {
		fakeItems, err := generateFakeItems(page, pageLimit, txFetcher.TotalSyncTxs)
		if err != nil {
			return errors.New("failed to generate fake items")
		}
		response.Data.Items = fakeItems
	}

	fmt.Println("Num Txs: ", len(response.Data.Items))

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return errors.New("failed to begin transaction")
	}

	// Prepare the statement for batch inserting
	fmt.Println("Prepare statement for batch inserting")
	stmt, err := tx.Prepare("INSERT INTO txs (address, tx_id, sync_id, raw, page, created_at) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (tx_id) DO NOTHING")
	if err != nil {
		return errors.New("failed to prepare statement for batch inserting")
	}
	defer stmt.Close()

	// Execute the statement for each tx
	now := time.Now().UTC()
	fmt.Println("Execute statement for each tx")
	for _, item := range response.Data.Items {
		fmt.Println("Inserting tx: ", item.TransactionId)
		txBytes, err = json.Marshal(item)
		if err != nil {
			return errors.New("failed to marshal JSON")
		}

		_, err = stmt.Exec(address, item.TransactionId, sync_id, string(txBytes), page, now)
		if err != nil {
			return errors.New("failed to execute statement for batch inserting")
		}
	}

	// Commit the transaction
	fmt.Println("Commit the transaction")
	if err := tx.Commit(); err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}

func getTotalNumberOfTxs(address string) (int, error) {
	if !testing {
		url := "https://rest.cryptoapis.io/blockchain-data/bitcoin/mainnet/addresses/" + address + "/transactions?context=yourExampleString&limit=" + strconv.Itoa(pageLimit) + "&offset=1304404"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			return -1, errors.New("failed to create request")
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", apiKey)

		res, err := client.Do(req)
		if err != nil {
			return -1, errors.New("failed to send request")
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return -1, errors.New("failed to read response body")
		}

		var response APIResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return -1, errors.New("failed to unmarshal JSON")
		}

		if limitNumResultsInProd {
			return limitNumResults, nil
		}
		return response.Data.Total, nil
	} else {
		return testTotalNumTxs, nil
	}
}
