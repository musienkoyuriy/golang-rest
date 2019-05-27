package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// createStrain func takes an sql.DB pointer
// insert's a name and race values into strains table, calls updateManyToManyRelation func
// and returns a possible error
func (strain *Strain) createStrain(db *sql.DB) error {
	var err error
	createStrainQuery := fmt.Sprintf("INSERT INTO strains(name, race) VALUES('%s','%s')",
		strain.Name, strain.Race)

	_, err = db.Exec(createStrainQuery)

	if err != nil {
		return err
	}

	updateManyToManyRelation(db, strain)

	return nil
}

// GetStrainsByCriteria func takes a sql.DB pointer and query parameters
// from request. It checks a search criteria and does appropriate SQL request to get the strains
// which meet a search criteria. Returns an array of strains with error
func GetStrainsByCriteria(db *sql.DB, vars map[string][]string) ([]Strain, error) {
	var strainQueryString = ""

	strainNameKeys, ok := vars["name"]
	if ok && len(strainNameKeys[0]) > 1 {
		strainQueryString = fmt.Sprintf("SELECT id FROM strains WHERE name = %s", wrapInQuotes(string(strainNameKeys[0])))
	}

	strainRaceKeys, ok := vars["race"]
	if ok && len(strainRaceKeys[0]) > 1 {
		if !ValidRace(strainRaceKeys[0]) {
			return nil, errors.New("Invalid race name")
		}
		strainQueryString = fmt.Sprintf("SELECT id FROM strains WHERE race = %s", wrapInQuotes(string(strainRaceKeys[0])))
	}

	strainFlavorKeys, ok := vars["flavor"]
	if ok && len(strainFlavorKeys[0]) > 1 {
		strainQueryString = fmt.Sprintf("SELECT DISTINCT s.id FROM flavors f INNER JOIN strain_flavors sf INNER JOIN  strains s ON f.id = sf.flavorId AND s.id = sf.strainId WHERE f.name = %s", wrapInQuotes(string(strainFlavorKeys[0])))
	}

	strainEffectKeys, ok := vars["effect"]
	if ok && len(strainEffectKeys[0]) > 1 {
		strainQueryString = fmt.Sprintf("SELECT DISTINCT s.id FROM effects e INNER JOIN strain_effects se INNER JOIN  strains s ON e.id = se.effectId AND s.id = se.strainId WHERE e.name = %s", wrapInQuotes(string(strainEffectKeys[0])))
	}

	strainRows, err := db.Query(strainQueryString)

	strainIds := []string{}

	if err != nil {
		return nil, err
	}

	for strainRows.Next() {
		var strainID int

		if err = strainRows.Scan(&strainID); err != nil {
			return nil, err
		}

		strainIds = append(strainIds, strconv.Itoa(strainID))
	}

	getStrainsByNameQuery := fmt.Sprintf("SELECT strains.id, strains.name, strains.race, flavors.name, effects.name, effects.type FROM strains INNER JOIN effects INNER JOIN strain_effects INNER JOIN flavors INNER JOIN strain_flavors ON strain_effects.strainId = strains.id AND effects.id = strain_effects.effectId AND flavors.id = strain_flavors.flavorId AND strains.id = strain_flavors.strainId WHERE strains.id IN (%s)", joinByComma(strainIds))

	rows, err := db.Query(getStrainsByNameQuery)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	strains, err = mapStrainsRecordsToJSON(rows)

	if err != nil {
		return nil, err
	}

	return strains, nil
}

// mapStrainsRecordsToJSON func takes a records from SELECT query
// and maps it to a JSON format to send it for client
func mapStrainsRecordsToJSON(rows *sql.Rows) ([]Strain, error) {
	strains := []Strain{}

	for rows.Next() {
		var id int
		var name, race, flavor, effectName, effectType string

		if err := rows.Scan(&id, &name, &race, &flavor, &effectName, &effectType); err != nil {
			return nil, err
		}

		if !strainsHasID(id, strains) {
			var strain Strain

			strain.ID = id
			strain.Name = name
			strain.Race = race
			strain.Flavors = []string{flavor}
			strain.Effects = map[string][]string{
				effectType: []string{effectName}}

			strains = append(strains, strain)
		} else {
			strainIndex, ok := getStrainIndexByID(id, strains)
			strain := strains[strainIndex]

			if ok {
				if !strainHasFlavor(strain, flavor) {
					strain.Flavors = append(strain.Flavors, flavor)
				}

				if !strainHasEffect(strain, effectType, effectName) {
					strain.Effects[effectType] = append(strain.Effects[effectType], effectName)
				}
			}
		}
	}

	return strains, nil
}

// DeleteStrain func takes a database pointer and strain id
// it deletes a strain record from a strains table by ID parameter
func DeleteStrain(db *sql.DB, id int) error {
	deleteStrainQuery := fmt.Sprintf("DELETE FROM strains WHERE id = %d", id)

	res, err := db.Exec(deleteStrainQuery)

	if err == nil {
		count, err := res.RowsAffected()
		if err == nil && count == 0 {
			return errors.New("Invalid strain ID")
		}
		return nil
	}
	return err
}

// editStrainByID func takes an sql.DB pointer
// updates's a name and race values in strains table,
// delete's strain flavors from strain_flavors table, strain effects strain_effects,
// calls updateManyToManyRelation func
// and returns a possible error
func (strain *Strain) editStrainByID(db *sql.DB) error {
	editStrainQuery := fmt.Sprintf("UPDATE strains SET name = '%s', race = '%s' WHERE id = %d", strain.Name, strain.Race, strain.ID)

	_, err := db.Exec(editStrainQuery)

	if err != nil {
		return err
	}

	deleteStrainFlavorsTableQuery := fmt.Sprintf("DELETE FROM strain_flavors WHERE strainId = %d; ", strain.ID)
	_, err = db.Exec(deleteStrainFlavorsTableQuery)

	if err != nil {
		return err
	}

	deleteStrainEffectsTableQuery := fmt.Sprintf("DELETE FROM strain_effects WHERE strainId = %d", strain.ID)
	_, err = db.Exec(deleteStrainEffectsTableQuery)

	if err != nil {
		return err
	}

	updateManyToManyRelation(db, strain)

	return err
}

// updateManyToManyRelation takes a database pointer and strain pointer
// it updates a relation's tables (strain_effects and strain_flavors)
func updateManyToManyRelation(db *sql.DB, strain *Strain) error {
	var effectsList = []string{}

	for _, effectNames := range strain.Effects {
		for _, effectName := range effectNames {
			effectsList = append(effectsList, "'"+effectName+"'")
		}
	}

	getEffectIdsQuery := fmt.Sprintf("SELECT id FROM effects WHERE name IN (%s)", joinByComma(effectsList))

	effects, err := db.Query(getEffectIdsQuery)

	if err != nil {
		return err
	}

	effectIds := []int{}

	for effects.Next() {
		var effectID int

		if err = effects.Scan(&effectID); err != nil {
			return err
		}

		effectIds = append(effectIds, effectID)
	}

	valuesToInsert := []string{}

	for _, effectID := range effectIds {
		valuesToInsert = append(valuesToInsert, fmt.Sprintf("(%d, %d)", strain.ID, effectID))
	}

	queryString := fmt.Sprintf("INSERT INTO strain_effects (strainId, effectId) VALUES %s", joinByComma(valuesToInsert))

	_, err = db.Exec(queryString)

	if err != nil {
		return err
	}

	flavorsInQuotes := WrapElementsInQuotes(strain.Flavors)

	editStrainQuery := fmt.Sprintf("SELECT id FROM flavors WHERE name IN (%s)", joinByComma(flavorsInQuotes))

	flavors, err := db.Query(editStrainQuery)

	if err != nil {
		return err
	}

	flavorIds := []int{}

	for flavors.Next() {
		var flavorID int

		if err = flavors.Scan(&flavorID); err != nil {
			return err
		}

		flavorIds = append(flavorIds, flavorID)
	}

	valuesToInsert = []string{}

	for _, flavorID := range flavorIds {
		valuesToInsert = append(valuesToInsert, fmt.Sprintf("(%d, %d)", strain.ID, flavorID))
	}

	queryString = fmt.Sprintf("INSERT INTO strain_flavors (strainId, flavorId) VALUES %s", joinByComma(valuesToInsert))

	_, err = db.Exec(queryString)

	if err != nil {
		return err
	}

	return nil
}
