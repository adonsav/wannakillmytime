package repository

import (
	"database/sql"
	"github.com/adonsav/fgoapp/internal/config"
)

type postgresDBRepo struct {
	dbrepoAppConfig *config.AppConfig
	DB              *sql.DB
}

func NewPostgresDBRepo(conn *sql.DB, ac *config.AppConfig) DatabaseRepo {
	return &postgresDBRepo{
		dbrepoAppConfig: ac,
		DB:              conn,
	}
}

type testDBRepo struct {
	dbrepoAppConfig *config.AppConfig
	DB              *sql.DB
}

func NewTestingDBRepo(ac *config.AppConfig) DatabaseRepo {
	return &testDBRepo{
		dbrepoAppConfig: ac,
	}
}
