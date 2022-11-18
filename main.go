package main

import (
	"fmt"
	"flag"
	"sort"
	"encoding/json"
	"candystore/ws"
)

type TopCustomer struct {
	Name string `json:"name"`
	FavoriteSnack string `json:"favoriteSnack"`
	TotalSnacks int `json:"totalSnacks"`
}

func metrics(snacks []webscraper.Snack) (string, int) {
	totalSnacks := 0
	favoriteSnackEaten := 0
	favoriteSnack := ""

	for _, snack := range snacks {
		if snack.Eaten > favoriteSnackEaten {
			favoriteSnackEaten = snack.Eaten
			favoriteSnack = snack.Brand
		}
		totalSnacks += snack.Eaten
	}

	return favoriteSnack, totalSnacks
}

const defaultURL = "https://candystore.zimpler.net/"

func main() {
	url := flag.String("url", defaultURL, "url to fetch data from")
    flag.Parse()

	CustomersTable := webscraper.CollectData(*url)
	topCustomers := make([]TopCustomer, 0)

	for name, snacks := range CustomersTable {
		favoriteSnack, totalSnacks := metrics(snacks)

		topCustomer :=	TopCustomer{
						Name: name,
						FavoriteSnack: favoriteSnack,
						TotalSnacks: totalSnacks,
					}
		topCustomers = append(topCustomers, topCustomer)
	}

	sort.Slice(topCustomers, func(i, j int) bool {
		return topCustomers[i].TotalSnacks > topCustomers[j].TotalSnacks
	})

	topCustomersOutput, err := json.MarshalIndent(topCustomers, "", "  ")

	if err != nil {
		fmt.Println("Error during marshalling")
		return
	}

	fmt.Println(string(topCustomersOutput))
}
