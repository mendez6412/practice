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
	Latitude	float64		`json:"latitude"`
	Longitude	float64		`json:"longitude"`
}

type GeoJson struct {
	Type		string		`json:"type"`
	Features	[]GeoJson	`json:"features,omitempty"`
	Geometry	Geometry	`json:"geometry,omitempty"`
	Properties	Properties	`json:"properties,omitempty"`
}

type Geometry struct {
	Type		string		`json:"type"`
	Coordinates	[]float64	`json:"coordinates"`
}

type Properties struct {
	Name		string		`json:"name"`
}


var addresses []IPAddress

func main() {
	ReadCsv()
	router := mux.NewRouter()
	router.HandleFunc("/getAddress/{id}", GetAddress).Methods("GET", "OPTIONS")
	router.HandleFunc("/getAddressesByBoundary/{swLat}/{swLon}/{neLat}/{neLon}", GetAddressesByBoundary).Methods("GET")
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
		var lat, lon float64
		lat, latErr := strconv.ParseFloat(line[7], 64)
		lon, lonErr := strconv.ParseFloat(line[8], 64)
		if latErr != nil || lonErr != nil {
			// The first row has unparseable strings (cos they're words!), so better error handling here would be good
			// log.Fatal("error on converting float: ", latErr)
		} else {
			addresses = append(addresses, IPAddress{
				Network: line[0],
				Latitude: lat,
				Longitude: lon,
			})
		}
	}
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

func GetAddressesByBoundary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	var swLat, swLon float64
	swLat, sLatErr := strconv.ParseFloat(params["swLat"], 64)
	swLon, sLonErr := strconv.ParseFloat(params["swLon"], 64)
	neLat, nLatErr := strconv.ParseFloat(params["neLat"], 64)
	neLon, nLonErr := strconv.ParseFloat(params["neLon"], 64)

	var filtered []GeoJson
	var responseGeoJson GeoJson


	if sLatErr != nil || sLonErr != nil || nLatErr != nil || nLonErr != nil {
		log.Fatal("error in getting address by boundary")
	} else {
		for _, e := range addresses {
			if (e.Latitude > swLat && e.Latitude < neLat) && (e.Longitude > swLon && e.Longitude < neLon) {
				filtered = append(filtered, GeoJson{
					Type: "Feature",
					Properties: Properties{
						Name: e.Network,
					},
					Geometry: Geometry{
						Type: "Point",
						Coordinates: []float64{e.Longitude, e.Latitude},
					},
				})
			}
		}
		responseGeoJson = GeoJson{
			Type: "FeatureCollection",
			Features: filtered,
		}
		json.NewEncoder(w).Encode(responseGeoJson)
	}

}
