package main

import (
	"database/sql"
	"fmt"
	"strings"
)

func (strain *Strain) createStrain(db *sql.DB) error {
	var err error
	createStrainQuery := fmt.Sprintf("INSERT INTO strains(name, race, flavors, effects) VALUES('%s','%s','%s', '%s')",
		strain.Name, strain.Race, strain.Flavors, strain.Effects)

	_, err = db.Exec(createStrainQuery)

	if err != nil {
		return err
	}

	_ = CreateStrainWithData(db, strain)

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
					on efs.strainId = s.id
					and ef.id = efs.effectId
					and f.id = sf.flavorId
					and s.id = sf.strainId
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

func GetStrainsByCriteria(db *sql.DB, criteria string, criteriaValue string) ([]Strain, error) {
	var err error

	getStrainsByNameQuery := fmt.Sprintf("SELECT strains.id, strains.name, flavors.name, effects.name, effects.type FROM strains INNER JOIN effects INNER JOIN strain_effects INNER JOIN flavors INNER JOIN strain_flavors ON strain_effects.strain_id = strains.id AND effects.id = strain_effects.effectId AND flavors.id = strain_flavors.flavorId AND strains.id = strain_flavors.strainId WHERE %s = %s", criteria, criteriaValue)

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
	deleteStrainFlavorsTableQuery := fmt.Sprintf("DELETE FROM strain_flavors WHERE strainId = %d", strain.ID)
	_, err := db.Exec(deleteStrainFlavorsTableQuery)

	if err != nil {
		return err
	}

	deleteStrainEffectsTableQuery := fmt.Sprintf("DELETE FROM strain_effects WHERE strainId = %d", strain.ID)
	_, err = db.Exec(deleteStrainEffectsTableQuery)

	if err != nil {
		return err
	}

	editStrainQuery := fmt.Sprintf("SELECT id FROM flavors WHERE name IN (%s)", joinFlavors(strain.Flavors))

	flavors, err := db.Query(editStrainQuery)

	if err != nil {
		return err
	}

	flavorIds := make([]int, 100)

	for flavors.Next() {
		var flavorID int

		if err = flavors.Scan(&flavorID); err != nil {
			return err
		}

		flavorIds = append(flavorIds, flavorID)
	}

	// for i := 0; i < len(flavorIds) {
	// 	queryString := fmt.Sprintf("INSERT INTO strain_flavors (strainId, flavorId) VALUES (%d, %d), (%d, %d), (%d, %d)", strain.ID, flavorIds[0], strain.ID, flavorIds[1], strain.ID, flavorIds[2])
	// }

	return err
}

func joinFlavors(flavors []string) string {
	return strings.Join(flavors[:], ", ")
}
