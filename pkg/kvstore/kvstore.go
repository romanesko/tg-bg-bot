package kvstore

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/nutsdb/nutsdb"
)

var (
	db     *nutsdb.DB
	once   sync.Once
	err    error
	bucket = "default_bucket"
)

func initDB() {
	once.Do(func() {
		// Get the current working directory
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}

		// Use a subdirectory for the database
		dbPath := filepath.Join(dir, "storage")
		opt := nutsdb.DefaultOptions
		opt.Dir = dbPath

		// Open the database
		db, err = nutsdb.Open(opt)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}

		err = db.Update(func(tx *nutsdb.Tx) error {
			return tx.NewBucket(nutsdb.DataStructureBTree, bucket)
		})

	})
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}
}

// Write stores a key-value pair with an optional TTL (in seconds).
// If ttlSeconds is 0, the key persists indefinitely.
func Write(key string, value []byte, ttlSeconds uint32) error {
	initDB()
	return db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(bucket, []byte(key), value, ttlSeconds)
	})
}

// Read retrieves the value for a key, returning "" if the key doesn't exist or has expired.
func Read(key string) []byte {
	initDB()
	var value []byte
	err := db.View(func(tx *nutsdb.Tx) error {
		entry, err := tx.Get(bucket, []byte(key))
		if err != nil {
			if errors.Is(err, nutsdb.ErrNotFoundKey) {
				return nil // Key not found or expired, return empty string
			}
			return err
		}
		value = entry
		return nil
	})
	if err != nil {
		log.Printf("Error reading key %s: %v", key, err)
		return nil
	}
	return value
}

// Close shuts down the database (optional, for cleanup).
func Close() error {
	initDB()
	return db.Close()
}
