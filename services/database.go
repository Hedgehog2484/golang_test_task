package services

import (
	"context"
	"os"
	_ "fmt"
	
	"github.com/jackc/pgx/v5"

	"test_golang_task/custom_types"
)


type Database struct {
	Conn *pgx.Conn
}


func NewDatabase(database_url string) Database {
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	return Database{Conn: conn}
}


func (db *Database) CloseConnection() {
	db.Conn.Close(context.Background())
}


func (db *Database) GetAllShelfsNames() ([]string, error) {
	rows, _ := db.Conn.Query(context.Background(), "SELECT distinct name FROM shelfs")
	var shelfs_names []string;
	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			return shelfs_names, err
		}
		shelfs_names = append(shelfs_names, name)
	}
	return shelfs_names, rows.Err()
}


func (db *Database) GetShelfsWithOrderedProducts() ([]custom_types.Shelf, error) {
	rows, _ := db.Conn.Query(context.Background(), "SELECT shelfs.name AS shelf_name, shelfs.product_id AS product_id, orders.id AS order_id, orders.amount, products.name FROM shelfs INNER JOIN orders ON shelfs.product_id=orders.product_id right join products on orders.product_id=products.id WHERE is_main=true")

	var shelfs []custom_types.Shelf
	for rows.Next() {
		var name string
		var productID int
		var orderID int
		var amount int
		var productName string

		err := rows.Scan(&name, &productID, &orderID, &amount, &productName)
		if err != nil {
			return shelfs, err
		}	
		shelfs = append(shelfs, custom_types.Shelf{Name: name, OrderID: orderID, ProductID: productID, Amount: amount, ProductName: productName})
	}
	return shelfs, rows.Err()
}


func (db *Database) GetAdditionalShelfByProductId(product_id int) (string, error) {
	rows, _ := db.Conn.Query(context.Background(), "SELECT name FROM shelfs WHERE product_id=$1 AND is_main=false", product_id)
	
	additional_shelfs := ""
	for rows.Next() {
		var shelf_name string

		err := rows.Scan(&shelf_name)
		if err != nil {
			return additional_shelfs, err
		}
		additional_shelfs = additional_shelfs + shelf_name + " "
	}
	return additional_shelfs, rows.Err()
}

