package connection

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type configKey string

const (
	POSTGRES_HOST    configKey = "POSTGRES_HOST"
	POSTGRES_PORT    configKey = "POSTGRES_PORT"
	POSTGRES_USER    configKey = "POSTGRES_USER"
	POSTGRES_PASS    configKey = "POSTGRES_PASS"
	POSTGRES_DBNAME  configKey = "POSTGRES_DBNAME"
	POSTGRES_SSLMODE configKey = "POSTGRES_SSLMODE"
)

func ConnectDB() *sql.DB {
	_ = godotenv.Load(".env")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		loadConfig(POSTGRES_HOST),
		loadConfig(POSTGRES_PORT),
		loadConfig(POSTGRES_USER),
		loadConfig(POSTGRES_PASS),
		loadConfig(POSTGRES_DBNAME),
		loadConfig(POSTGRES_SSLMODE),
	)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		panic(err)
	}

	var errDatabase string

	err = db.Ping()

	if err != nil {
		errDatabase = strings.Split(err.Error(), "pq:")[1]
	}

	if errDatabase == fmt.Sprintf(` database "%s" does not exist`, loadConfig(POSTGRES_DBNAME)) {
		db.Close()

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s",
			loadConfig(POSTGRES_HOST),
			loadConfig(POSTGRES_PORT),
			loadConfig(POSTGRES_USER),
			loadConfig(POSTGRES_PASS),
			loadConfig(POSTGRES_SSLMODE),
		)

		db, err = sql.Open("postgres", dsn)

		if err != nil {
			panic(err)
		}

		dbName := fmt.Sprint(loadConfig(POSTGRES_DBNAME))

		_, err = db.Exec("create database " + dbName)

		if err != nil {
			panic(err)
		}

		db.Close()
	} else {
		db.SetConnMaxIdleTime(1 * time.Minute)
		db.SetConnMaxLifetime(1 * time.Minute)
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(10)

		return db
	}

	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		loadConfig(POSTGRES_HOST),
		loadConfig(POSTGRES_PORT),
		loadConfig(POSTGRES_USER),
		loadConfig(POSTGRES_PASS),
		loadConfig(POSTGRES_DBNAME),
		loadConfig(POSTGRES_SSLMODE),
	)

	db, err = sql.Open("postgres", dsn)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return db
}

func loadConfig(key configKey) string {
	return os.Getenv(string(key))
}
