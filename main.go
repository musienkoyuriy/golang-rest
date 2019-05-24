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

// Init the database and connection
func (strainsApi *StrainsAPI) Init(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	strainsApi.DB, err = sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	InsertDataFromJSON(strainsApi.DB)

	strainsApi.initRoutes()
}

// Strain struct
type Strain struct {
	ID      int                 `json:"id"`
	Name    string              `json:"name"`
	Race    string              `json:"race"`
	Flavors []string            `json:"flavors"`
	Effects map[string][]string `json:"effects"`
}

func (strainsApi *StrainsAPI) initRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/newStrain", strainsApi.createNewStrain).Methods("POST")
	router.HandleFunc("/strains/{criteria}{name}", strainsApi.getStrainsByCriteria).Methods("GET")
	router.HandleFunc("strain/{id}", strainsApi.deleteStrain).Methods("DELETE")
	router.HandleFunc("strain/{id}", strainsApi.editStrain).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8888", router))
}

func main() {
	strainsAPI := StrainsAPI{}
	strainsAPI.Init("root", "password", "flourishdb")
}

func (strainsApi *StrainsAPI) createNewStrain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var strain Strain
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&strain); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}
	defer r.Body.Close()

	if err := strain.createStrain(strainsApi.DB); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}

	response, _ := json.Marshal(strain)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (strainsApi *StrainsAPI) getStrainsByCriteria(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	criteria := vars["criteria"]
	criteriaValue := vars["name"]

	if isValidCriteria(criteria) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid criteria to search strains"))
		return
	}

	if strains, err := GetStrainsByCriteria(strainsApi.DB, criteria, criteriaValue); err != nil {
		print(strains) // TODO remove
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response, _ := json.Marshal(strains)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (strainsApi *StrainsAPI) deleteStrain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please specify correct strain ID"))
		return
	}

	if err := DeleteStrain(strainsApi.DB, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response, _ := json.Marshal(map[string]int{
		"deletedStrain": id})

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (strainsApi *StrainsAPI) editStrain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please specify a valid strain ID"))
		return
	}

	var strain Strain

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&strain); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please specify correct strain data to update"))
		return
	}

	strain.ID = id

	if err := strain.editStrainByID(strainsApi.DB); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response, _ := json.Marshal(strain)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func isValidCriteria(criteria string) bool {
	switch criteria {
	case "race", "flavor", "effect":
		return true
	}
	return false
}
