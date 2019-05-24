package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

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
		// ADD RACE
		createRaceQuery := fmt.Sprintf("INSERT IGNORE INTO races(name) VALUES('%s')", strain.Race)

		_, err = db.Exec(createRaceQuery)

		if err != nil {
			panic("Error while adding a race")
		}

		var raceID int

		err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&raceID)

		if err != nil {
			panic("Error while getting an new race ID")
		}

		// ADD STRAIN
		createStrainQuery := fmt.Sprintf("INSERT INTO strains(name, raceId) VALUES('%s','%d')",
			strain.Name, raceID)

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
