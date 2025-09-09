package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RandomUserResponse struct {
	Results []User `json:"results"`
}

type User struct {
	Name  Name   `json:"name"`
	Email string `json:"email"`
}

type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

func main() {
	start := time.Now()
	total := 0

	// 10 pages × 5000 results = 50,000
	for page := 1; page <= 10; page++ {
		url := fmt.Sprintf("https://randomuser.me/api/?results=100&seed=myseed&page=%d", page)

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var data RandomUserResponse
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("Error on page", page, ":", err)
			continue
		}

		total += len(data.Results)

		fmt.Println("Page", page, "users:", len(data.Results))
	}

	fmt.Printf("Processed %d users\n", total)
	fmt.Println("Elapsed time:", time.Since(start))
}
