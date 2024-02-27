package loader

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDb(path string, dbName string) (*sql.DB, error) {
	// Open a connection to the MySQL database
	db, err := sql.Open(dbName, path)
	return db, err
}

func StartLoading(ctx context.Context, db *sql.DB, loadFunc func(*sql.DB), waitSecond time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(waitSecond * time.Second):
		}

		loadFunc(db)
	}
}

