package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/gorilla/mux"
)

type IPAddress struct {
	Network		string		`json:"id"`
	Latitude	string		`json:"latitude"`
	Longitude	string		`json:"longitude"`
}

var addresses []IPAddress

func main() {
	ReadCsv()
	router := mux.NewRouter()
	router.HandleFunc("/getAddress/{id}", GetAddress).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func ReadCsv() {
	csvFile, _ := os.Open("ipv4.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
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
}

func GetAddress(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id, err = strconv.Atoi(params["id"])
	if err!= nil {
		log.Fatal("error on converting int")
	}
	var address = FindAddressByIndex(id)
	json.NewEncoder(w).Encode(address)
}

func FindAddressByIndex(index int) IPAddress {
	return addresses[index]
}
