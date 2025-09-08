package main

import (
	"encoding/json"
	"io"
	"net/http"
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

func main() {

	resp, err := http.Get("https://randomuser.me/api/?results=2&seed=myseed&page=1")
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
	//fmt.Println(data.Results[0].Name.First)

}
