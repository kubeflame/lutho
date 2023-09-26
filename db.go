package main

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

func InitDB() (*bolt.DB, error) {
	db, err := bolt.Open(homeDir+"data.db", 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db, nil
}
