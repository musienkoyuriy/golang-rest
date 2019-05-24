package importscript

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Strain struct
type Strain struct {
	ID      int                 `json:"id"`
	Name    string              `json:"name"`
	Race    string              `json:"race"`
	Flavors []string            `json:"flavors"`
	Effects map[string][]string `json:"effects"`
}

func connectToStrainsDB(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	InsertDataFromJSON(db)
}

// InsertDataFromJSON func
func InsertDataFromJSON(db *sql.DB) {
	connectToStrainsDB("root", "password", "flourishdb")

	jsonFile, err := os.Open("strains.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var strains []Strain

	json.Unmarshal(byteValue, &strains)

	for i := 0; i < len(strains); i++ {
		// ADD STRAIN
		createStrainQuery := fmt.Sprintf("INSERT INTO strains(id, name, race) VALUES('%d','%s','%s')",
			strains[i].ID, strains[i].Name, strains[i].Race)

		_, err = db.Exec(createStrainQuery)

		if err != nil {
			panic("Error while adding a strain")
		}

		// ADD FLAVORS
		for j := 0; j < len(strains[i].Flavors); j++ {
			flavor := strains[i].Flavors[j]

			createFlavorQuery := fmt.Sprintf("INSERT IGNORE INTO flavors(name) VALUES('%s')",
				flavor)

			_, err = db.Exec(createFlavorQuery)

			if err != nil {
				panic("Error while adding a flavor")
			}
		}

		fmt.Println(strains[i].Effects)

		// ADD EFFECTS
		for effectType, effectNames := range strains[i].Effects {
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
