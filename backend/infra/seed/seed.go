package seed

import (
	"SamporDoc/backend/infra/model"
	"SamporDoc/backend/infra/repository"
	"SamporDoc/backend/utils"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

//go:embed customers.csv
var customersSeed []byte


var DefaultShops []model.Shop = []model.Shop{
	{Slug: "HJKmain", Name: "หจก.3พ.รุ่งเรืองศึกษาภัณฑ์ (หลัก)", SortingLevel: 0},
	{Slug: "HJKsec", Name: "หจก.3พ.รุ่งเรืองศึกษาภัณฑ์ (รอง)", SortingLevel: 1},
	{Slug: "BUMmain", Name: "ร้านบำเหน็จศึกษาภัณฑ์ (หลัก)", SortingLevel: 2},
	{Slug: "BUMsec", Name: "ร้านบำเหน็จศึกษาภัณฑ์ (รอง)", SortingLevel: 3},
}


func SeedShops(repo *repository.Repo) {
	fmt.Println("Start seeding shops...")
	for _, shop := range DefaultShops {
		err := repo.CreateShop(&shop)
		if err != nil {
			fmt.Println("Error seeding database -> repo.CreateShop")
			panic(err)
		}
	}
	fmt.Println("Finish seeding shops")
}

func SeedCustomers(repo *repository.Repo) {
	reader := csv.NewReader(strings.NewReader(string(customersSeed)))
	var customers []model.Customer
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error seeding database -> reading customer.csv")
			panic(err)
		}
		if len(record) != 8 {
			fmt.Println("Error seeding database -> reading customer.csv")
			panic(fmt.Errorf("insufficient column in customer.csv"))
		}
		customers = append(customers, model.Customer{
			Name: record[0],
			Address: utils.Ptr(record[1]),
			HeadCheckerName: utils.Ptr(record[2]),
			Checker1Name: utils.Ptr(record[3]),
			Checker2Name: utils.Ptr(record[4]),
			ObjectName: utils.Ptr(record[5]),
			HeadObjectName: utils.Ptr(record[6]),
			BossName: utils.Ptr(record[7]),
		})
	}

	// insert to db
	fmt.Println("Start seeding customers...")
	for _, sch := range customers {
		err := repo.CreateCustomer(&sch)
		if err != nil {
			fmt.Println("Error seeding database -> repo.CreateCustomer")
			panic(err)
		}
	}
	fmt.Println("Finish seeding customers")
}
