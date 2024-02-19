## Prerequisites

1. Install [Go](https://go.dev/doc/install).
2. Install [PostgreSQL](https://www.postgresql.org/download/macosx/)
3. Create a local Postgres database with the name: cointracker, port: 5432 and host: localhost.

## Running the App Locally

This command runs the app on your machine.

```bash
make run
```

You should see "Successfully connected!" which means the app was able to connect to the local cointracker db.

Note: 
The app by default uses fake data to simulate the address transactions list API response from cryptoapis.io since it's expensive to test.
To use the actual cryptoapis.io response, set testing to false in txfetcher.go. 