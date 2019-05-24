package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var strains []Strain

type StrainsAPI struct {
	Router *mux.Router
	DB     *sql.DB
}

func (strainsApi *StrainsAPI) Init(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	strainsApi.DB, err = sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	insertDataFromJSON(strainsApi.DB)

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

func (strain *Strain) createStrain(db *sql.DB) error {
	var err error
	createStrainQuery := fmt.Sprintf("INSERT INTO strains(name, race, flavors, effects) VALUES('%s','%s','%s', '%s')",
		strain.Name, strain.Race, strain.Flavors, strain.Effects)

	_, err = db.Exec(createStrainQuery)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&strain.ID)

	if err != nil {
		return err
	}

	return nil
}

func _getStrainsByCriteria(db *sql.DB, criteria string, criteriaValue string) ([]Strain, error) {
	var err error
	getStrainsByNameQuery := fmt.Sprintf("SELECT * FROM strains WHERE %s = %s", criteria, criteriaValue)

	rows, err := db.Query(getStrainsByNameQuery)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	strains := []Strain{}

	for rows.Next() {
		var strain Strain

		if err := rows.Scan(&strain.ID, &strain.Name, &strain.Race, &strain.Flavors, &strain.Effects); err != nil {
			return nil, err
		}

		strains = append(strains, strain)
	}

	return strains, nil
}

func _deleteStrain(db *sql.DB, id int) error {
	deleteStrainQuery := fmt.Sprintf("DELETE FROM strains WHERE id = %d", id)

	_, err := db.Exec(deleteStrainQuery)
	return err
}

func (strain *Strain) editStrainByID(db *sql.DB) error {
	editStrainQuery := fmt.Sprintf("UPDATE strains SET name='%s', race='%s', flavors='%s', effects='%s' WHERE id='%d'", strain.Name, strain.Race, strain.Flavors, strain.Effects, strain.ID)

	_, err := db.Exec(editStrainQuery)
	return err
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
	// createDatabase()
	strainsApi := StrainsAPI{}
	strainsApi.Init("root", "password", "flourishdb")
}

func createDatabase() {
	// file, err := ioutil.ReadAll("./db.sql")

	// if err != nil {
	// 	// handle error
	// }

	// requests := strings.Split(string(file), ";")

	// for _, request := range requests {
	// 	result, err := db.Exec(request)
	// 	// do whatever you need with result and error
	// }
}

func insertDataFromJSON(db *sql.DB) {
	jsonFile, err := os.Open("strains.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var strains []Strain

	json.Unmarshal(byteValue, &strains)

	for _, strain := range strains {
		// ADD STRAIN
		createStrainQuery := fmt.Sprintf("INSERT INTO strains(id, name, race) VALUES('%d','%s','%s')",
			strain.ID, strain.Name, strain.Race)

		_, err = db.Exec(createStrainQuery)

		if err != nil {
			panic("Error while adding a strain")
		}

		// ADD FLAVORS
		for _, flavor := range strain.Flavors {
			createFlavorQuery := fmt.Sprintf("INSERT IGNORE INTO flavors(name) VALUES('%s')",
				flavor)

			_, err = db.Exec(createFlavorQuery)

			if err != nil {
				panic("Error while adding a flavor")
			}
		}

		fmt.Println(strain.Effects)

		// ADD EFFECTS
		for effectType, effectNames := range strain.Effects {
			fmt.Println(effectType)
			fmt.Println(effectNames)

			for _, effectName := range effectNames {
				createEffectQuery := fmt.Sprintf("INSERT IGNORE INTO effects(name, type) VALUES('%s', '%s')",
					effectName, effectType)

				_, err = db.Exec(createEffectQuery)

				if err != nil {
					panic("Error while adding an effect")
				}
			}
		}
	}
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

	if strains, err := _getStrainsByCriteria(strainsApi.DB, criteria, criteriaValue); err != nil {
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

	if err := _deleteStrain(strainsApi.DB, id); err != nil {
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
