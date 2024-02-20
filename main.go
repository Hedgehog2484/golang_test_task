package main

import (
	"os"
	"fmt"
	"strconv"

	"test_golang_task/services"
	_ "test_golang_task/custom_types"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	args := os.Args

	db := services.NewDatabase(os.Getenv("DATABASE_URL"))

	shelfs_names, err := db.GetAllShelfsNames()
	if err != nil {
		panic(err)
	}

	shelfs, err := db.GetShelfsWithOrderedProducts()
	if err != nil {
		panic(err)
	}

	var isShelfPrinted bool
	for i := 0; i < len(shelfs_names); i++ {
		isShelfPrinted = false;
		for j := 0; j < len(shelfs); j++ {
			if shelfs[j].Name == shelfs_names[i]{
				for l := 1; l < len(args); l++ {
					order_id, _ := strconv.Atoi(args[l])
					if shelfs[j].OrderID == order_id {
						if isShelfPrinted == false {
							fmt.Println("==========\nСтеллаж:", shelfs_names[i])
							isShelfPrinted = true
						}

						fmt.Printf("\nТовар: %v (id=%v)\nЗаказ: %v\nКол-во: %v\n", 
						shelfs[j].ProductName, shelfs[j].ProductID, shelfs[j].OrderID, shelfs[j].Amount)
						additional_shelfs, err := db.GetAdditionalShelfByProductId(shelfs[j].ProductID)
						if err != nil {
							panic(err)
						}
						if additional_shelfs != "" {
							fmt.Println("Доп. стелажи:", additional_shelfs)
						}
					}
				}
			}
		}
	}
}

