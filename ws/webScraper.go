package webscraper

import (
	"fmt"
	"strconv"
	"github.com/gocolly/colly/v2"
)

type Snack struct {
	Brand string
	Eaten int
}

type Customer struct {
	Name  string
	Snack Snack
}

func contains(allSnacksPerCustomer []Snack, snackName string) (bool, int) {
	// Check if the snack already exists
    for i, c := range allSnacksPerCustomer {
        if c.Brand == snackName {
			return true, i
		}
    }
    return false, 0
}

func CollectData(url string) map[string][]Snack {
	// Fetch the data from the url, getting specific data about the snack table
	c := colly.NewCollector()

	customersData := make(map[string][]Snack)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting Scraping on:", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response Status:", r.StatusCode)
	})

	c.OnHTML(`#top\.customers > tbody`, func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			eaten, err := strconv.Atoi(el.ChildText("td:nth-child(3)"))

			if err != nil {
				fmt.Println("Error during conversion")
				return
			}

			customer := Customer{
				Name:  el.ChildText("td:nth-child(1)"),
				Snack: Snack{el.ChildText("td:nth-child(2)"), eaten},
			}

			// Check if the customer has already bought the snack
			// If so, add the eaten amount to the existing snack
			// it makes easier later to find which snack is the favorite
			// since the amount itself is not enough to define it.
			if _, ok :=	customersData[customer.Name]; ok {
				if exists, index := contains(customersData[customer.Name], customer.Snack.Brand); exists {
					customersData[customer.Name][index].Eaten += customer.Snack.Eaten
				} else {
					customersData[customer.Name] = append(customersData[customer.Name], customer.Snack)
				}
			} else {
				customersData[customer.Name] = []Snack{customer.Snack}
			}
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "Failed with response:", r, "Error:", err)
	})

	c.Visit(url)
	return customersData
}
