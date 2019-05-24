package main

import (
	"database/sql"
	"fmt"
)

//TODO
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

// FIX
func getStrainsByEffect(db *sql.DB, effectName string) ([]Strain, error) {
	var err error
	getStrainsByEffectQuery := `
			strainsWithEffect = SELECT
			DISTINCT s.id
			FROM strains s
			INNER JOIN effects ef
			INNER JOIN strain_effects efs
			ON ef.id = efs.effectId
			AND s.id = efs.strainId
			WHERE ef.name = 'effect name'

			SELECT
					s.id
					from strains s
					INNER JOIN effects ef
					INNER JOIN strain_effects efs
					INNER JOIN flavors f
					INNER JOIN strain_flavors sf
					INNER JOIN races r
					on efs.strainId = s.id
					and ef.id = efs.effectId
					and f.id = sf.flavorId
					and s.id = sf.strainId
					and r.id = s.raceId
					WHERE s.id IN strainsWithEffect
		`

	rows, err := db.Query(getStrainsByEffectQuery)

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

// TODO
func GetStrainsByCriteria(db *sql.DB, criteria string, criteriaValue string) ([]Strain, error) {
	var err error
	getStrainsByNameQuery := `
		SELECT
			strains.id, strains.name, races.name, flavors.name, effects.name, effects.type, 
		FROM strains
		INNER JOIN effects
		INNER JOIN strain_effects
		INNER JOIN flavors
		INNER JOIN strain_flavors
		INNER JOIN races
		ON strain_effects.strain_id = strains.id
		AND effects.id = strain_effects.effectId
		AND flavors.id = strain_flavors.flavorId
		AND strains.id = strain_flavors.strainId
		AND races.id = strains.raceId
		WHERE name = ?
	`

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

func DeleteStrain(db *sql.DB, id int) error {
	deleteStrainQuery := fmt.Sprintf("DELETE FROM strains WHERE id = %d", id)

	_, err := db.Exec(deleteStrainQuery)
	return err
}

//TODO
func (strain *Strain) editStrainByID(db *sql.DB) error {
	editStrainQuery := fmt.Sprintf("UPDATE strains SET name='%s', race='%s', flavors='%s', effects='%s' WHERE id='%d'", strain.Name, strain.Race, strain.Flavors, strain.Effects, strain.ID)

	_, err := db.Exec(editStrainQuery)
	return err
}
