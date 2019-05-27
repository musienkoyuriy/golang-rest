package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var strains []Strain

// StrainsAPI struct
type StrainsAPI struct {
	Router *mux.Router
	DB     *sql.DB
}

// Init method belongs to StrainsAPI structure and
// intended for database connection and routes initialization
func (strainsApi *StrainsAPI) Init(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	strainsApi.DB, err = sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	strainsApi.initRoutes()
}

// Strain structure is a data model for main entity - Strain
// It has ID, Name, Race, Flavors and Effects fields
type Strain struct {
	ID      int                 `json:"id"`
	Name    string              `json:"name"`
	Race    string              `json:"race"`
	Flavors []string            `json:"flavors"`
	Effects map[string][]string `json:"effects"`
}

// initRoutes iniializes routes for searching strains, adding, editing and deleting
func (strainsApi *StrainsAPI) initRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/newStrain", strainsApi.createNewStrain).Methods("POST")
	router.HandleFunc("/strains", strainsApi.getStrainsByCriteria).Methods("GET")
	router.HandleFunc("/strains/{id}", strainsApi.deleteStrain).Methods("DELETE")
	router.HandleFunc("/editStrain", strainsApi.editStrain).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8888", router))
}

// main is an entry point for app
// calls route's initialization method
func main() {
	strainsAPI := StrainsAPI{}
	strainsAPI.Init("root", "password", "flourishdb")
}

// createNewStrain is /createNewStrain route handler
// takes a response and request objects, decode's request body and do appropriate
// POST request for adding a new strain
func (strainsApi *StrainsAPI) createNewStrain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var strain Strain
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&strain); err != nil {
		sendHTTPError(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := strain.createStrain(strainsApi.DB); err != nil {
		sendHTTPError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(strain)

	sendJSONResponse(w, response)
}

// getStrainsByCriteria is /strains route handler
// takes a response and request objects, get's a query parameter to determine a search criteria
// and returns a JSON with matched strai items
func (strainsApi *StrainsAPI) getStrainsByCriteria(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := r.URL.Query()

	if strains, err := GetStrainsByCriteria(strainsApi.DB, vars); err != nil {
		print(strains)
		sendHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(strains)

	sendJSONResponse(w, response)
}

func sendHTTPError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func sendJSONResponse(w http.ResponseWriter, response []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// deleteStrain is /strain/{id} route handler
// it deletes an existing strain by specified ID parameter
func (strainsApi *StrainsAPI) deleteStrain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendHTTPError(w, "Please specify correct strain ID", http.StatusBadRequest)
		return
	}

	if err := DeleteStrain(strainsApi.DB, id); err != nil {
		sendHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(map[string]int{
		"deletedStrain": id})

	sendJSONResponse(w, response)
}

// editStrain is /editStrain route handler
// takes a response and request objects, decode's request body and
// edits an existing strain record
func (strainsApi *StrainsAPI) editStrain(w http.ResponseWriter, r *http.Request) {
	var strain Strain

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&strain); err != nil {
		sendHTTPError(w, "Please specify correct strain data to update", http.StatusBadRequest)
		return
	}

	if err := strain.editStrainByID(strainsApi.DB); err != nil {
		sendHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(strain)

	sendJSONResponse(w, response)
}
