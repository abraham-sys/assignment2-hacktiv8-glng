package main

import (
	"assignment2/db/connection"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	db := connection.ConnectDB()

	err := down(db)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database successfully reverted")
	}
}

func down(db *sql.DB) error {
	defer db.Close()
	var dropTableItems string = `
		DROP TABLE items;
		DROP TABLE orders;
	`

	_, err := db.Exec(dropTableItems)

	if err != nil {
		return err
	}

	return err
}
