package postgres

import (
	"os"

	"github.com/jmoiron/sqlx"
)

func ConnectToDatabase() (*sqlx.DB, func()) {
	url := os.Getenv("DATABASE_URI")
	db := sqlx.MustConnect("postgres", url)
	return db, func() {
		_ = db.Close()
	}
}
