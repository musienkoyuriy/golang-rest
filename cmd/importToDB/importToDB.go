package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// StrainID is an alias for an integer type
type StrainID int

// Strain structure is a data model for main entity - Strain
// It has ID, Name, Race, Flavors and Effects fields
type Strain struct {
	ID      int                 `json:"id"`
	Name    string              `json:"name"`
	Race    string              `json:"race"`
	Flavors []string            `json:"flavors"`
	Effects map[string][]string `json:"effects"`
}

func main() {
	connectionString := fmt.Sprintf("%s:%s@/%s", "root", "password", "flourishdb")

	var err error
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	InsertDataFromJSON(db)
}

// CreateStrainWithData func iterares over a JSON structure, parse it
// and inserts data in an appropriate database tables
func CreateStrainWithData(db *sql.DB, strain Strain) {
	// ADD STRAIN
	createStrainQuery := fmt.Sprintf("INSERT INTO strains(name, race) VALUES('%s', '%s')",
		strain.Name, strain.Race)

	_, err := db.Exec(createStrainQuery)

	if err != nil {
		fmt.Println("Error while adding a strain")
	}

	var strainID StrainID

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&strainID)

	if err != nil {
		panic("Error while getting an new strain ID")
	}

	// ADD FLAVORS
	for _, flavor := range strain.Flavors {
		createFlavorQuery := fmt.Sprintf("INSERT IGNORE INTO flavors(name) VALUES('%s')",
			flavor)

		_, err = db.Exec(createFlavorQuery)

		if err != nil {
			panic("Error while adding a flavor")
		}

		var flavorID int

		err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&flavorID)

		createStrainFlavorQuery := fmt.Sprintf("INSERT IGNORE INTO strain_flavors(strainId, flavorId) VALUES('%d', '%d')", strainID, flavorID)

		_, err = db.Exec(createStrainFlavorQuery)

		if err != nil {
			panic("Error while adding a strain_flavor")
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

			var effectID int

			err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&effectID)

			createStrainFlavorQuery := fmt.Sprintf("INSERT IGNORE INTO strain_effects(strainId, effectId) VALUES('%d', '%d')", strainID, effectID)

			_, err = db.Exec(createStrainFlavorQuery)

		}
	}
}

// InsertDataFromJSON to dynamic fill database from JSON file
func InsertDataFromJSON(db *sql.DB) {
	jsonFile, err := os.Open("strains.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var strains []Strain

	json.Unmarshal(byteValue, &strains)

	for _, strain := range strains {
		CreateStrainWithData(db, strain)
	}
}
