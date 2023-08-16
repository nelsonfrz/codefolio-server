package app

import (
	"codefolio/server/database"
	"context"
	"database/sql"
	"os"
)

type App struct {
	db      *sql.DB
	queries *database.Queries
	ctx     context.Context
}

func NewApp(databaseUrl string) (*App, error) {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))

	if err != nil {
		return nil, err
	}

	queries := database.New(db)
	ctx := context.Background()

	return &App{db: db, queries: queries, ctx: ctx}, nil
}
