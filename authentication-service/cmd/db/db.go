package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var counts int64

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	//dsn := "postgresql://gocommmicro:password@localhost:5432/gocommmicro?sslmode=disable"

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds...")

		time.Sleep(time.Second * 2)
		continue
	}
}
func ConnectToTestDB() *sql.DB {
	testDSN := os.Getenv("TESTDBDSN")
	if testDSN == "" {
		testDSN = "postgresql://gocommmicro:password@localhost:5432/testdb?sslmode=disable"
	}
	//dsn := "postgresql://gocommmicro:password@localhost:5432/gocommmicro?sslmode=disable"
	//log.Println("TEST DSN:", testDSN)
	for {
		connection, err := openDB(testDSN)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds...")

		time.Sleep(time.Second * 2)
		continue
	}
}
