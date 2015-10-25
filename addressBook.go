
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
  	"gopkg.in/mgo.v2/bson"
  	"io/ioutil"
  	"os"
  	"strings"
)

type AddressInput struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

type AddressUpdate struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

type ResultOutput struct {
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

type ResponseAddOp struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name  		string `json:"name"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State 		string `json:"state"`
	Zip   		string `json:"zip"`
	Coordinate struct {
		Lat float64 
		Lng float64 
	} `json:"coordinate"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/locations", handlePostAdress).Methods("POST")
	router.HandleFunc("/locations/{imdbKey}", handleAdress).Methods("GET", "DELETE", "PUT")
	http.ListenAndServe(":8080", router)
}

func handleAdress(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	imdbKey := vars["imdbKey"]
	sess, err := mgo.Dial("mongodb://kalpana:pass@ds035844.mongolab.com:35844/cmpe273")
	if err != nil {
	    fmt.Printf("Can't connect to mongo, go error %v\n", err)
	    os.Exit(1)
	}
	defer sess.Close()		 
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("AddressBook")

	switch req.Method {
	case "GET":		
		result := ResponseAddOp{}
		
		err = collection.Find(bson.M{"_id": bson.ObjectIdHex(imdbKey)}).Select(bson.M{}).One(&result)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
		    fmt.Printf("Address not found: %v\n", err)
		}
		outputRes,_:= json.Marshal(result)
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res,string(outputRes))
	case "PUT":
		addressVal:= new(AddressUpdate)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&addressVal)
		if error != nil {
			log.Println("ERRR1",error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		inputAdd:=addressVal.Address+"+"+addressVal.City+"+"+addressVal.State+"+"+addressVal.Zip
		inputAdd=strings.Replace(inputAdd, " ", "+", -1)
		var ro ResultOutput
		url :="http://maps.google.com/maps/api/geocode/json?address="+inputAdd+"&sensor=false"
        responseValue, err := http.Get(url)
 		if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        } else {
            defer responseValue.Body.Close()
            reply, err := ioutil.ReadAll(responseValue.Body)
            json.Unmarshal([]byte(reply), &ro)
            if err != nil {
                fmt.Printf("%s", err)
                os.Exit(1)
        	}

	        colQuerier := bson.M{"_id": bson.ObjectIdHex(imdbKey)}
	        change := bson.M{"$set": bson.M{"address": addressVal.Address, "city": addressVal.City,"state":addressVal.State,"zip":addressVal.Zip,"coordinate.lat":ro.Results[0].Geometry.Location.Lat,"coordinate.lng":ro.Results[0].Geometry.Location.Lng}}
			err = collection.Update(colQuerier, change)	
			if err != nil {
				//res.WriteHeader(http.StatusNotFound)
			    fmt.Printf("Address not found: %v\n", err)
			}
		
			result := ResponseAddOp{}
			err = collection.Find(bson.M{"_id": bson.ObjectIdHex(imdbKey)}).Select(bson.M{}).One(&result)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
			    fmt.Printf("Address not found: %v\n", err)
			}
			outputRes,_:= json.Marshal(result)
			//res.WriteHeader(http.StatusNoContent)
			res.WriteHeader(http.StatusCreated)
			fmt.Fprint(res,string(outputRes))
        }
	case "DELETE":
		err = collection.Remove(bson.M{"_id": bson.ObjectIdHex(imdbKey)})
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
			    fmt.Printf("Address not found: %v\n", err)
			}
		fmt.Fprint(res,"Record Deleted.")
		res.WriteHeader(http.StatusOK)
	}
}

func handlePostAdress(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		addressVal:= new(AddressInput)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&addressVal)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		inputAdd:=addressVal.Address+"+"+addressVal.City+"+"+addressVal.State+"+"+addressVal.Zip
		inputAdd=strings.Replace(inputAdd, " ", "+", -1)


		var ro ResultOutput
		url :="http://maps.google.com/maps/api/geocode/json?address="+inputAdd+"&sensor=false"
		//fmt.Println(url)
        responseValue, err := http.Get(url)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        } else {
            defer responseValue.Body.Close()
            reply, err := ioutil.ReadAll(responseValue.Body)
            json.Unmarshal([]byte(reply), &ro)
            if err != nil {
                fmt.Printf("%s", err)
                os.Exit(1)
        	}
        	//fmt.Println("Lat ",ro.Results[0].Geometry.Location.Lat)
        	//fmt.Println("Lng ",ro.Results[0].Geometry.Location.Lng)

			res2D := &ResponseAddOp{
				ID: bson.NewObjectId(),
		        Name: addressVal.Name,
		        Address: addressVal.Address,
		        City:  addressVal.City,
		        State: addressVal.State,
		        Zip: addressVal.Zip,
	        	Coordinate: struct{
						Lat float64 
						Lng float64 
	        	}{ro.Results[0].Geometry.Location.Lat, ro.Results[0].Geometry.Location.Lng},}

		    	sess, err := mgo.Dial("mongodb://kalpana:pass@ds035844.mongolab.com:35844/cmpe273")
				if err != nil {
				    fmt.Printf("Can't connect to mongo, go error %v\n", err)
				    os.Exit(1)
				}
				defer sess.Close()				 
				sess.SetSafe(&mgo.Safe{})
				collection := sess.DB("cmpe273").C("AddressBook")
				err = collection.Insert(res2D)
				if err != nil {
				    fmt.Printf("Can't insert document: %v\n", err)
				    os.Exit(1)
				}
			    outgoingJSON, err := json.Marshal(res2D)
				if err != nil {
				log.Println(error.Error())
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
				}
				res.WriteHeader(http.StatusCreated)
				fmt.Fprint(res, string(outgoingJSON))
        	}
	}

