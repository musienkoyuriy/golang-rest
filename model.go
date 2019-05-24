package main

import (
	"database/sql"
	"fmt"
)

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
