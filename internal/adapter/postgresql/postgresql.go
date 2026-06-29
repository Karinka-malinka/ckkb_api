package postgresql

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

type PostgreSQLStorage struct {
	db *sql.DB
}

const DSN = "host=rc1b-cagcn3odpc9bj7ch.mdb.yandexcloud.net port=6432 user=gbu password=1qaz2wsX@ dbname=gbu_config"

func NewDB() (*sql.DB, error) {

	db, err := sql.Open("pgx", DSN)
	if err != nil {
		logrus.Errorf("error open driver pgx: %#v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	logrus.Info("connected to postgres")

	return db, nil
}

func NewPostgreSQLStorage(db *sql.DB) *PostgreSQLStorage {
	return &PostgreSQLStorage{db: db}
}
