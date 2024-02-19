package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	db_port  = 5432
	user     = ""
	password = ""
	dbname   = "cointracker"
)

const address = "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h"

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	psqlInfo := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable",
		host, db_port, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	createTablesSQL := `
    CREATE TABLE IF NOT EXISTS addresses (
        id SERIAL PRIMARY KEY,
        address TEXT UNIQUE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
	
	CREATE TABLE IF NOT EXISTS syncs (
		id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		address text,
		status text,
		txs_synced integer,
		total_sync_txs integer,
		total_txs integer,
		created_at timestamp with time zone,
		finished_at timestamp with time zone
	);

	CREATE TABLE IF NOT EXISTS sync_errors (
		id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		sync_id text,
		page integer,
		error text
	);

	CREATE TABLE IF NOT EXISTS txs (
		id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		address text,
		tx_id text UNIQUE,
		raw jsonb,
		page integer,
		created_at timestamp with time zone,
		sync_id text
	);
	`
	_, err = db.Exec(createTablesSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	} else {
		log.Println("Table created or already exists")
	}

	txFetcher := GetNewTxFetcher(db, address, pageLimit)
	txFetcher.SyncTxs()
}
