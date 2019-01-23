package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

type IPAddress struct {
	Network		string		`json:"id"`
	Latitude	string		`json:"latitude"`
	Longitude	string		`json:"longitude"`
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	csvFile, _ := os.Open("ipv4.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var addresses []IPAddress
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		addresses = append(addresses, IPAddress{
			Network: line[0],
			Latitude: line[7],
			Longitude: line[8],
		})
	}
	// addressJson, _ := json.Marshal(addresses)
	json.NewEncoder(w).Encode(addresses)
}

type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
    City  string `json:"city,omitempty"`
    State string `json:"state,omitempty"`
}

var people []Person
