package main

import (
	"assignment2/db/connection"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	db := connection.ConnectDB()

	err := initTable(db)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database successfully initialized")
	}
}

func initTable(db *sql.DB) error {
	defer db.Close()
	var createTableOrders string = `
		CREATE TABLE orders (
			order_id serial primary key,
			customer_name varchar(255),
			ordered_at timestamptz,
			quantity integer
		);

		CREATE TABLE items (
			item_id serial,
			item_code varchar(255),
			description text,
			quantity integer,
			order_id integer,
			PRIMARY KEY (item_id),
			FOREIGN KEY (order_id) REFERENCES orders(order_id)
		);
		`

	_, err := db.Exec(createTableOrders)

	if err != nil {
		return err
	}

	return err
}
