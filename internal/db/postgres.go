package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func ConnectPostgres() {

	conn, err := pgx.Connect(
		context.Background(),
		"postgres://scheduler:scheduler@localhost:5432/scheduler_db",
	)

	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}

	Conn = conn

	log.Println("Connect to PostgreSQL")
}
