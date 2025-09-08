package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

/*
 simple sequential Go example that fetches a few batches from RandomUser.me
 so you can test the logic before scaling up to 50,000 users.
 https://randomuser.me/api/?results=5000&seed=myseed&page=1
https://randomuser.me/api/?results=5000&seed=myseed&page=2
...
https://randomuser.me/api/?results=5000&seed=myseed&page=10

*/

type RandomUserResponse struct {
	Results []User `json:"results"`
	Info    Info   `json:"info"`
}

type User struct {
	Name    Name    `json:"name"`
	Email   string  `json:"email"`
	Login   Login   `json:"login"`
	Phone   string  `json:"phone"`
	Cell    string  `json:"cell"`
	Nat     string  `json:"nat"`
	Picture Picture `json:"picture"`
}

type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type Login struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

type Picture struct {
	Large string `json:"large"`
}

type Info struct {
	Seed    string `json:"seed"`
	Results int    `json:"results"`
	Page    int    `json:"page"`
}

func worker(jobs <-chan User, results chan<- User, wg *sync.WaitGroup, counter *uint64) {
	for jb := range jobs {
		results <- jb
		wg.Done()

		// increment counter atomically
		atomic.AddUint64(counter, 1)
	}
}

func main() {

	start := time.Now()

	resp, err := http.Get("https://randomuser.me/api/?results=5000&seed=myseed&page=1")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data RandomUserResponse
	err = json.Unmarshal(body, &data)

	// create workerpool

	var wg sync.WaitGroup

	const numWorkers = 5
	var counter uint64
	jobs := make(chan User, len(data.Results))
	results := make(chan User, len(data.Results))

	// start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(jobs, results, &wg, &counter)
	}

	// send all data as jobs
	for _, d := range data.Results {
		wg.Add(1)
		jobs <- d
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println(r.Name.First)
	}

	fmt.Printf("Processed %d users\n", counter)
	fmt.Println("elapsed time: ", time.Since(start))

}
