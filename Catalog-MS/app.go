package main

import (
	"encoding/json"
	"log"
	"net/http"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	. "./ccs"
	. "./models"
)

var ccs = CCSDB{}

// Function to fetch all the drinks in the Catalog
func FetchAllDrinksEndPoint(w http.ResponseWriter, r *http.Request) {
	drinks, err := ccs.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, drinks)
}

// GET a Drink by ID
func FetchDrinkEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	drink, err := ccs.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Drink Doesn't Exist")
		return
	}
	respondWithJson(w, http.StatusOK, drink)
}

// POST a new Drink into the Catalog
func AddaDrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var drink Drink
	if err := json.NewDecoder(r.Body).Decode(&drink); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	drink.ID = bson.NewObjectId()
	if err := ccs.Insert(drink); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, drink)
}

// PUT update an existing Drink
func UpdateADrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var drink Drink
	if err := json.NewDecoder(r.Body).Decode(&drink); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := ccs.Update(drink); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "successfully updated the drink"})
}

// DELETE an Drink from the Catalog
func DeleteADrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var drink Drink
	if err := json.NewDecoder(r.Body).Decode(&drink); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := ccs.Delete(drink); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "successfully added a new drink"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func init() {
	ccs.Connect()
	ccs.ConnecttoPrimary()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/menu", FetchAllDrinksEndPoint).Methods("GET")
	r.HandleFunc("/addadrink", AddaDrinkEndPoint).Methods("POST")
	r.HandleFunc("/updatedrink", UpdateADrinkEndPoint).Methods("PUT")
	r.HandleFunc("/delete", DeleteADrinkEndPoint).Methods("DELETE")
	r.HandleFunc("/drink/{id}", FetchDrinkEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3001", r); err != nil {
		log.Fatal(err)
	}
}
