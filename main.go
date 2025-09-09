package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
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

func workers(w int, jobs <-chan int, results chan<- []User, wg *sync.WaitGroup) {

	// 10 pages × 100 results = 1000
	for page := range jobs {
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

		//total += len(data.Results)

		fmt.Println("Page", page, "users:", len(data.Results))
		results <- data.Results
	}
	wg.Done()
}

func main() {
	start := time.Now()

	var wg sync.WaitGroup
	const numPages = 10
	const numWorkers = 5

	jobs := make(chan int, numPages)
	results := make(chan []User, numPages)

	// start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go workers(w, jobs, results, &wg)
	}

	// start jobs
	for j := 1; j <= numPages; j++ {
		jobs <- j
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	//receive results channel
	for res := range results {
		total += len(res)
	}

	fmt.Printf("Processed %d users\n", total)
	fmt.Println("Elapsed time:", time.Since(start))
}
