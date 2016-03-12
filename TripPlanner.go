package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LongLat struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

type PriceEst struct {
	Prices []struct {
		//	CurrencyCode         string  `json:"currency_code"` //optional
		//	DisplayName          string  `json:"display_name"` //optional
		Distance float64 `json:"distance"`
		Duration int     `json:"duration"`
		//	Estimate             string  `json:"estimate"` //optional
		//	HighEstimate         int     `json:"high_estimate"` //optional
		//	LocalizedDisplayName string  `json:"localized_display_name"` //optional
		LowEstimate int `json:"low_estimate"`
		Visited     bool
		//	Minimum              int     `json:"minimum"` //optional
		ProductID string `json:"product_id"`
		//		SurgeMultiplier      int     `json:"surge_multiplier"` //optional
	} `json:"prices"`
}

type BestRoute struct {
	BestRouteLocationIds   []string `json:"best_route_location_ids" bson:"best_route_location_ids"`
	ID                     string   `json:"id" bson:"id"`
	StartingFromLocationID string   `json:"starting_from_location_id" bson:"starting_from_location_id"`
	Status                 string   `json:"status" bson:"status"`
	TotalDistance          float64  `json:"total_distance" bson:"total_distance"`
	TotalUberCosts         int      `json:"total_uber_costs" bson:"total_uber_costs"`
	TotalUberDuration      int      `json:"total_uber_duration" bson:"total_uber_duration"`
}

type EstimateTime struct {
	//Driver          interface{} `json:"driver"` //optional
	Eta int `json:"eta"`
	//Location        interface{} `json:"location"` //optional
	//	RequestID       string      `json:"request_id"` //optional
	//	Status          string      `json:"status"` //optional
	//SurgeMultiplier int         `json:"surge_multiplier"` //optional
	//Vehicle         interface{} `json:"vehicle"` //optional
}

type TripRequest struct {
	BestRouteLocationIds      []string `json:"best_route_location_ids" bson:"best_route_location_ids"`
	ID                        string   `json:"id" bson:"id"`
	NextDestinationLocationID string   `json:"next_destination_location_id" bson:"next_destination_location_id"`
	StartingFromLocationID    string   `json:"starting_from_location_id" bson:"starting_from_location_id"`
	Status                    string   `json:"status" bson:"status"`
	TotalDistance             float64  `json:"total_distance" bson:"total_distance"`
	TotalUberCosts            int      `json:"total_uber_costs" bson:"total_uber_costs"`
	TotalUberDuration         int      `json:"total_uber_duration" bson:"total_uber_duration"`
	UberWaitTimeEta           int      `json:"uber_wait_time_eta" bson:"uber_wait_time_eta"`
}

type Info struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Zip     string `json:"zip" bson:"zip"`
}

type InfoReturn struct {
	Name        string     `json:"name" bson:"name"`
	Address     string     `json:"address" bson:"address"`
	City        string     `json:"city" bson:"city"`
	State       string     `json:"state" bson:"state"`
	Zip         string     `json:"zip" bson:"zip"`
	ID          string     `json:"id" bson:"id"`
	Coordinates Coordinate `json:"coordinate" bson:"coordinate"`
	Visited     bool
}
type InfoReturns struct {
	Name        string     `json:"name" bson:"name"`
	Address     string     `json:"address" bson:"address"`
	City        string     `json:"city" bson:"city"`
	State       string     `json:"state" bson:"state"`
	Zip         string     `json:"zip" bson:"zip"`
	ID          string     `json:"id" bson:"id"`
	Coordinates Coordinate `json:"coordinate" bson:"coordinate"`
}
type Coordinate struct {
	Lng float64 `json:"lng" bson:"lng"`
	Lat float64 `json:"lat" bson:"lat"`
}

type JsonCoordinates struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	StartLat  float64 `json:"start_latitude" bson:"start_latitude"`
	StartLng  float64 `json:"start_longitude" bson:"start_longitude"`
	EndLat    float64 `json:"end_latitude" bson:"end_latitude"`
	EndLng    float64 `json:"end_longitude" bson:"end_longitude"`
}

type Locations struct {
	Start       string   `json:"starting_from_location_id" bson:"starting_from_location_id"`
	LocationIDs []string `json:"location_ids" bson:"location_ids"`
}

func hello1(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	//	fmt.Fprintf(rw, "Hello, %s!\n", p.ByName("name"))

	var g InfoReturn
	//var g InfoReturn

	//fmt.Fprintf(rw, "Hello, %s!\n", p.ByName("name"))

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&g)
	if err != nil {
		fmt.Println("error in decoding")
	}

	StartQuery := "http://maps.google.com/maps/api/geocode/json?address="
	WhereQuery := g.Address + " " + g.City + " " + g.State
	WhereQuery = strings.Replace(WhereQuery, " ", "+", -1)

	EndQuery := "&sensor=false"
	Url := StartQuery + WhereQuery + EndQuery

	resp, err := http.Get(Url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	var LL LongLat

	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &LL); err != nil {
		panic(err)
	}

	g.Coordinates.Lat = LL.Results[0].Geometry.Location.Lat
	g.Coordinates.Lng = LL.Results[0].Geometry.Location.Lng

	maxWait := time.Duration(20 * time.Second)
	session, err := mgo.DialWithTimeout("Your MongoDB URL Here", maxWait)
	if err != nil {
		fmt.Println(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("locations").C("peoples")

	g.ID = p.ByName("name")
	err = c.Update(bson.M{"id": p.ByName("name")}, bson.M{"$set": bson.M{"address": g.Address, "city": g.City, "state": g.State, "zip": g.Zip, "coordinate": bson.M{"lat": g.Coordinates.Lat, "lng": g.Coordinates.Lng}}})
	if err != nil {
		log.Fatal(err)
	}

	err = c.Find(bson.M{"id": p.ByName("name")}).One(&g)
	if err != nil {
		fmt.Println("Record not found")
	}

	mapB, _ := json.Marshal(g)
	fmt.Fprintf(rw, string(mapB))

}

func trip(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	maxWait := time.Duration(20 * time.Second)
	session, err := mgo.DialWithTimeout("Your MongoDB URL Here", maxWait)

	if err != nil {
		fmt.Println(err)
	}

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("locations").C("peoples")

	Results := BestRoute{}
	err = c.Find(bson.M{"id": p.ByName("name")}).One(&Results)
	if err != nil {
		fmt.Println("Record not found")

	}
	rw.WriteHeader(http.StatusOK)
	mapB, _ := json.Marshal(Results)
	fmt.Fprintf(rw, string(mapB))

}

func greeting1(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	maxWait := time.Duration(20 * time.Second)
	session, err := mgo.DialWithTimeout("mongodb://sachingunari:hpprotego1@ds039504.mongolab.com:39504/locations", maxWait)
	if err != nil {
		fmt.Println(err)
	} //else {
	//fmt.Println("Session created")
	//}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("locations").C("peoples")

	result := InfoReturn{}
	err = c.Find(bson.M{"name": "John Smith"}).One(&result)
	if err != nil {
		fmt.Println("ID not found")
		//log.Fatal(err)
	}

	_, err = c.RemoveAll(bson.M{"id": p.ByName("name")})
	if err != nil {
		fmt.Println("ID not found")
		//log.Fatal(err)
	}

}

func greeting(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var t Info
	var g InfoReturns
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println("error in decoding")
	}

	g.Name = t.Name
	g.Address = t.Address
	g.City = t.City
	g.State = t.State
	g.Zip = t.Zip

	StartQuery := "http://maps.google.com/maps/api/geocode/json?address="
	WhereQuery := t.Address + " " + t.City + " " + t.State
	WhereQuery = strings.Replace(WhereQuery, " ", "+", -1)

	//+ " " myjson3.City + " " + myjson3.State

	EndQuery := "&sensor=false"
	Url := StartQuery + WhereQuery + EndQuery

	resp, err := http.Get(Url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	var LL LongLat

	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &LL); err != nil {
		panic(err)
	}

	g.Coordinates.Lat = LL.Results[0].Geometry.Location.Lat
	g.Coordinates.Lng = LL.Results[0].Geometry.Location.Lng
	g.ID = strconv.Itoa(rand.Intn(100000))
	maxWait := time.Duration(20 * time.Second)
	session, err := mgo.DialWithTimeout("Your MongoDB Cloud URL", maxWait)

	if err != nil {
		fmt.Println(err)
	} else {
		//	fmt.Println("Session created")
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("locations").C("peoples")
	err = c.Insert(&g)
	if err != nil {
		fmt.Println("could not insert to DB")
		log.Fatal(err)
	}
	rw.WriteHeader(http.StatusCreated)
	mapB, _ := json.Marshal(g)
	fmt.Fprintf(rw, string(mapB))

}

func trips(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	var L Locations
	var IDs []string
	var Start InfoReturn
	var Len int
	var results []InfoReturn
	var Esti []PriceEst
	var EstiL PriceEst
	var min float64
	var Route BestRoute
	var LIndex int
	min = 100

	maxWait := time.Duration(20 * time.Second)
	session, err := mgo.DialWithTimeout("mongodb://sachingunari:hpprotego1@ds039504.mongolab.com:39504/locations", maxWait)

	decoder := json.NewDecoder(req.Body)
	errr := decoder.Decode(&L)
	if errr != nil {
		fmt.Println("error in decoding")
	}

	Len = len(L.LocationIDs)
	IDs = make([]string, Len)
	results = make([]InfoReturn, Len)
	for index, _ := range L.LocationIDs {
		IDs[index] = L.LocationIDs[index]
		// index is the index where we are
		// element is the element from someSlice for where we are
	}

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("locations").C("peoples")

	err = c.Find(bson.M{"id": L.Start}).One(&Start)
	if err != nil {
		fmt.Println("Record not found")

	}
	StartV := Start
	for index, _ := range L.LocationIDs {
		err = c.Find(bson.M{"id": IDs[index]}).One(&results[index])
		if err != nil {
			fmt.Println("Record not found")
			//	log.Fatal(err)
		}

	}
	Esti = make([]PriceEst, Len)
	Route.BestRouteLocationIds = make([]string, Len)
	Route.Status = "planning"
	Route.StartingFromLocationID = StartV.ID
	Route.ID = strconv.Itoa(rand.Intn(10000))

	for n, _ := range results {

		for m, _ := range results {

			if results[m].Visited == false {

				URL := "https://sandbox-api.uber.com/v1/estimates/price?server_token=" + " Your server token here" + "&start_latitude=" + strconv.FormatFloat(StartV.Coordinates.Lat, 'f', 7, 64) + "&start_longitude=" + strconv.FormatFloat(StartV.Coordinates.Lng, 'f', 7, 64) + "&end_latitude=" + strconv.FormatFloat(results[m].Coordinates.Lat, 'f', 7, 64) + "&end_longitude=" + strconv.FormatFloat(results[m].Coordinates.Lng, 'f', 7, 64)
				resp, errs := http.Get(URL)

				if errs != nil {
					fmt.Println("error fetching URL")
					//	log.Fatal(err)
				}

				body, errs := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				if errs != nil {
					fmt.Println("UBER Decoding error")
					//	log.Fatal(err)
				}
				errs = json.Unmarshal(body, &Esti[m])
				if errs != nil {
					fmt.Println("UBER Decoding error")
					//	log.Fatal(err)
				}

				if min > Esti[m].Prices[0].Distance {
					min = Esti[m].Prices[0].Distance
					LIndex = m

				}

			}
		}

		Route.TotalDistance = Route.TotalDistance + Esti[LIndex].Prices[0].Distance
		Route.TotalUberCosts = Route.TotalUberCosts + Esti[LIndex].Prices[0].LowEstimate
		Route.TotalUberDuration = Route.TotalUberDuration + Esti[LIndex].Prices[0].Duration
		results[LIndex].Visited = true
		Route.BestRouteLocationIds[n] = results[LIndex].ID
		StartV = results[LIndex]

		min = 100

	}

	URL := "https://sandbox-api.uber.com/v1/estimates/price?server_token=" + "Server Token Here" + "&start_latitude=" + strconv.FormatFloat(results[LIndex].Coordinates.Lat, 'f', 7, 64) + "&start_longitude=" + strconv.FormatFloat(results[LIndex].Coordinates.Lng, 'f', 7, 64) + "&end_latitude=" + strconv.FormatFloat(Start.Coordinates.Lat, 'f', 7, 64) + "&end_longitude=" + strconv.FormatFloat(Start.Coordinates.Lng, 'f', 7, 64)
	resp, errs := http.Get(URL)

	if errs != nil {
		fmt.Println("error fetchin URL")
		//	log.Fatal(err)
	}

	body, errs := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if errs != nil {
		fmt.Println("UBER Decoding error")
		//	log.Fatal(err)
	}
	errs = json.Unmarshal(body, &EstiL)
	if errs != nil {
		fmt.Println("UBER Decoding error")
		//	log.Fatal(err)
	}

	Route.TotalDistance = Route.TotalDistance + EstiL.Prices[0].Distance
	Route.TotalUberCosts = Route.TotalUberCosts + EstiL.Prices[0].LowEstimate
	Route.TotalUberDuration = Route.TotalUberDuration + EstiL.Prices[0].Duration

	err = c.Insert(&Route)
	if err != nil {
		fmt.Println("Could not insert in DB found")

	}
	rw.WriteHeader(http.StatusCreated)
	mapB, _ := json.Marshal(Route)
	fmt.Fprintf(rw, string(mapB))

}

func tripStart(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	Estimate := PriceEst{}

	maxWait := time.Duration(20 * time.Second)
	session, err := mgo.DialWithTimeout("mongodb://sachingunari:hpprotego1@ds039504.mongolab.com:39504/locations", maxWait)

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("locations").C("peoples")

	Results := TripRequest{}
	err = c.Find(bson.M{"id": p.ByName("name")}).One(&Results)
	if err != nil {
		fmt.Println("Record not found")

	}

	TripStatus := TripRequest{}
	var FirstTrip bool
	var NextLoc int
	FirstTrip = true
	var StartLocationID string

	if Results.StartingFromLocationID != Results.NextDestinationLocationID {

		for _, _ = range Results.BestRouteLocationIds {
			if Results.NextDestinationLocationID != "" {
				FirstTrip = false
			}
		}

		if FirstTrip == false {
			for a, _ := range Results.BestRouteLocationIds {
				if Results.BestRouteLocationIds[a] == Results.NextDestinationLocationID {
					NextLoc = a + 1
					break
				}
			}
		}

		TripStatus.StartingFromLocationID = Results.StartingFromLocationID
		if FirstTrip == true {
			StartLocationID = Results.StartingFromLocationID
			TripStatus.Status = "requesting"
			TripStatus.NextDestinationLocationID = Results.BestRouteLocationIds[0]
		} else {
			StartLocationID = Results.NextDestinationLocationID
			TripStatus.Status = "requesting"
			if NextLoc < len(Results.BestRouteLocationIds) {
				TripStatus.NextDestinationLocationID = Results.BestRouteLocationIds[NextLoc]
			} else {
				StartLocationID = Results.NextDestinationLocationID
				TripStatus.Status = "completed"
				TripStatus.NextDestinationLocationID = TripStatus.StartingFromLocationID

			}

		}

		TripStatus.BestRouteLocationIds = Results.BestRouteLocationIds
		TripStatus.ID = Results.ID

		TripStatus.TotalDistance = Results.TotalDistance
		TripStatus.TotalUberCosts = Results.TotalUberCosts
		TripStatus.TotalUberDuration = Results.TotalUberDuration

		Start := InfoReturn{}
		End := InfoReturn{}
		err = c.Find(bson.M{"id": StartLocationID}).One(&Start)
		if err != nil {
			fmt.Println("Record not found Start")
		}
		err = c.Find(bson.M{"id": TripStatus.NextDestinationLocationID}).One(&End)
		if err != nil {
			fmt.Println("Record not found End")
		}
		URL := "https://sandbox-api.uber.com/v1/estimates/price?server_token=nOIdFtPeZHT8922iF-6HLMllWlLhHKZA0YTM6UC2&start_latitude=" + strconv.FormatFloat(Start.Coordinates.Lat, 'f', 7, 64) + "&start_longitude=" + strconv.FormatFloat(Start.Coordinates.Lng, 'f', 7, 64) + "&end_latitude=" + strconv.FormatFloat(End.Coordinates.Lat, 'f', 7, 64) + "&end_longitude=" + strconv.FormatFloat(End.Coordinates.Lng, 'f', 7, 64)
		resp, errs := http.Get(URL)

		if errs != nil {
			fmt.Println("error fetchin URL")
			//	log.Fatal(err)
		}

		body, errs := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if errs != nil {
			fmt.Println("UBER Decoding error")
			//	log.Fatal(err)
		}
		errs = json.Unmarshal(body, &Estimate)
		if errs != nil {
			fmt.Println("UBER Decoding error")
			//	log.Fatal(err)
		}
		ProductID := Estimate.Prices[0].ProductID

		apiUrl := "https://sandbox-api.uber.com/v1/requests"

		JsonParam := JsonCoordinates{}
		JsonParam.ProductID = ProductID
		JsonParam.StartLat = Start.Coordinates.Lat
		JsonParam.StartLng = Start.Coordinates.Lng
		JsonParam.EndLat = End.Coordinates.Lat
		JsonParam.EndLng = End.Coordinates.Lng

		JsonStr, err := json.Marshal(JsonParam)
		if err != nil {
			fmt.Println("UBER Error")
			//	log.Fatal(err)
		}

		req, err = http.NewRequest("POST", apiUrl, bytes.NewBuffer(JsonStr))
		if err != nil {
			fmt.Println("UBER Error")
			//	log.Fatal(err)
		}

		req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiI2MzI4NGZkZS1jMDE1LTQ3YTktYWM1Zi01NmM3ZTY4ODExYzgiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6IjYwNGNkMDZjLTcxMTctNDUzMy04YjA1LTFmZGE2MmFhYTlhMiIsImV4cCI6MTQ1MDMzOTE3MiwiaWF0IjoxNDQ3NzQ3MTcyLCJ1YWN0IjoiNVp0T21CajFhcUp4bGkyR3FOZ3dLc1cyc1VRZkJhIiwibmJmIjoxNDQ3NzQ3MDgyLCJhdWQiOiJFb1EtWTFqTDNwM3pxMmZPQmppeThYYkpyRVVGOVRLeSJ9.M3Fri4WgULN15HrlH0RyUNhCApPntFvMp2RpfS-Zq9wTC7srykO3O_mbYnqo1h3FDhbB6Mkpx2OaR_K-w_SSRUVkHS9rKQoh_06WRkXlucKPjXM-ZztDFmAzSngM7-8EN61pO3tbYY1ZcAVhK03SumurEbxaEZXD8cOH_4bmQPT47wR6DI2qCcQHVuPPVcAefL9YV7A8sZZF6O9h9ByAcdg2sK-cYoziWae6_o0Dr8MJki9uhJ1YMoqBfQR3kRVGqCvFO_qo_4pxHjGFhZ6SlzEYv51UtjPcaHUJLaqTDLYtP-kVjkWbFWgL_eZm2dv9DUNBQjCXxac9uuc0C49_jQ")
		req.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("UBER Error")
			//	log.Fatal(err)
		}

		bodys, errs := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if errs != nil {
			fmt.Println("UBER Error")
			//	log.Fatal(err)
		}

		var TimeEst EstimateTime
		errs = json.Unmarshal(bodys, &TimeEst)
		if errs != nil {
			//	log.Fatal(err)
		}
		TripStatus.UberWaitTimeEta = TimeEst.Eta

		err = c.Update(bson.M{"id": p.ByName("name")}, TripStatus)
		if err != nil {
			log.Fatal(err)
		}
		rw.WriteHeader(http.StatusOK)
		mapB, _ := json.Marshal(TripStatus)
		fmt.Fprintf(rw, string(mapB))
	} else {
		rw.WriteHeader(http.StatusOK)
		Results.UberWaitTimeEta = 0
		Results.NextDestinationLocationID = ""
		mapB, _ := json.Marshal(Results)
		fmt.Fprintf(rw, string(mapB))
	}

}

func main() {

	mux := httprouter.New()
	mux.POST("/Location", greeting)
	mux.POST("/trips", trips)
	mux.GET("/trips/:name", trip)
	mux.PUT("/Location/:name", hello1)
	mux.PUT("/trips/:name/request", tripStart)
	mux.DELETE("/Location/:name", greeting1)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
