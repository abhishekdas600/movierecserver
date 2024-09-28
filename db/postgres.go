package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitializePostgreSQL(db **sql.DB) error {
	connStr := os.Getenv("POSTGRES_STR") 
	var err error

	*db, err = sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	(*db).SetMaxOpenConns(25)
	(*db).SetMaxIdleConns(25)
	(*db).SetConnMaxLifetime(0) 

	if err := (*db).Ping(); err != nil {
		return fmt.Errorf("cannot reach the database: %v", err)
	}

	log.Println("Connected to the database successfully")
	return nil
}
