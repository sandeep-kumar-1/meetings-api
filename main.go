package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Participant struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
	RSVP  string `json:"RSVP"`
}

type meetings struct {
	ID               string       	`json:"ID"`
	Title            string 		`json:"Title"`
	Participant      []Participant 	`json:"Participants"`
	StartTime        string			`json:"StartTime"`
	EndTime          string			`json:"EndTime"`
	Creationtimestap string 	  	`json:"CreateTime"`
}

type allEvents []meetings

var events = allEvents{}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func ScheduleAMeet(w http.ResponseWriter, r *http.Request) {
	var newMeet meetings
	
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event id, title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newMeet)

	events = append(events, newMeet)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newMeet)
}

func GetAMeet(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func GetAllMeet(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func ListofTime(w http.ResponseWriter, r *http.Request) {
	start_time := mux.Vars(r)["id_1"]
	end_time := mux.Vars(r)["id_2"]
	var output = allEvents{}
	for _, singleEvent := range events {
		if singleEvent.StartTime >= start_time && singleEvent.EndTime <= end_time {
			output = append(output, singleEvent)
		}
	}
	json.NewEncoder(w).Encode(output)
}

func ListEmail(w http.ResponseWriter, r *http.Request) {
	email_id:= mux.Vars(r)["id"]
	var output = allEvents{}
	for _, singleEvent := range events {
		for _, mail := range singleEvent.Participant {
			if mail.Email == email_id {
				output = append(output, singleEvent)
			}
		}
	}
	json.NewEncoder(w).Encode(output)

}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	
	router.HandleFunc("/meetings", ScheduleAMeet).Methods("POST")
	
	router.HandleFunc("/meetings", GetAllMeet).Methods("GET")
	router.HandleFunc("/meetings/{id}", GetAMeet).Methods("GET")

	router.HandleFunc("/meetings?start={id_1}&end={id_2}", ListofTime).Methods("GET")
	router.HandleFunc("/meetings?participant={id}", ListEmail).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
