package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	"webportfolio/internal/httptransport"
	"webportfolio/internal/service"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	dev := os.Getenv("ENV") == "dev"

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	dbConnector := pgdriver.NewConnector(
		pgdriver.WithAddr("postgres:5432"),
		pgdriver.WithDatabase("webportfolio-db"),
		pgdriver.WithUser(dbUser),
		pgdriver.WithPassword(dbPassword),
		pgdriver.WithInsecure(true),
	)

	if dev {
		dbConnector.Config().Addr = "127.0.0.1:5432"
	}

	secret := os.Getenv("SECRET")

	db := bun.NewDB(sql.OpenDB(dbConnector), pgdialect.New())

	service := service.New(db, secret)
	err := service.CreateRelations(context.Background())
	if err != nil {
		panic(err)
	}

	if !dev {
		// wait for postgres to start
		time.Sleep(5 * time.Second)
	}

	panic(httptransport.Handler(service, secret).Start(":8000"))
}
