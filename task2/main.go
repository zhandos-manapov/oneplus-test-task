package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

var apiUrl = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1"

type Currency struct {
	Id                               string  `json:"id"`
	Symbol                           string  `json:"symbol"`
	Name                             string  `json:"name"`
	Image                            string  `json:"image"`
	Current_price                    float32 `json:"current_price"`
	Market_cap                       int     `json:"market_cap"`
	Market_cap_rank                  int     `json:"markey_cap_rank"`
	Fully_diluted_valuation          int     `json:"fully_diluted_valuation"`
	Total_volume                     float32 `json:"Total_volume"`
	High_24h                         float32 `json:"high_24h"`
	Low_24h                          float32 `json:"low_24h"`
	Price_change_24h                 float32 `json:"price_change_24h"`
	Market_cap_change_24h            float32 `json:"market_cap_change_24h"`
	Market_cap_change_percentage_24h float32 `json:"market_cap_change_percentage_24h"`
	Circulating_supply               float32 `json:"circulating_supply"`
	Total_supply                     float32 `json:"total_supply"`
	Max_supply                       float32 `json:"max_supply"`
	Ath                              float32 `json:"ath"`
	Ath_change_percentage            float32 `json:"ath_change_percentage"`
	Ath_date                         string  `json:"ath_date"`
	Atl                              float32 `json:"atl"`
	Atl_change_percentage            float32 `json:"atl_change_percentage"`
	Atl_date                         string  `json:"atl_date"`
	Roi                              struct {
		Times      float32 `json:"times"`
		Currency   string  `json:"currency"`
		Percentage float32 `json:"percentage"`
	} `json:"roi"`
	Last_updated string `json:"last_updated"`
}

type Error struct {
	Message string `json:"message"`
}

func GetCurrency(c *fiber.Ctx) error {
	response, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(Error{Message: "Something went wrong"})
	}
	defer response.Body.Close()

	var currencies []Currency
	if err := json.NewDecoder(response.Body).Decode(&currencies); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(Error{Message: "Something went wrong"})
	}

	id := c.Params("id")
	if len(id) > 0 {
		for _, item := range currencies {
			if item.Id == id {
				return c.Status(fiber.StatusOK).JSON(item)
			}
		}
	}
	return c.Status(fiber.StatusOK).JSON(currencies)
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/api/v1/currency", GetCurrency)
	app.Get("/api/v1/currency/:id", GetCurrency)

	app.Listen(":3000")
}
