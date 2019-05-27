package main

import "strings"

// WrapElementsInQuotes takes an array of strings, iterates over array
// and returns modified array with qouted elements
func WrapElementsInQuotes(list []string) []string {
	var itemsInQuotes = []string{}

	for _, val := range list {
		itemsInQuotes = append(itemsInQuotes, wrapInQuotes(val))
	}

	return itemsInQuotes
}

// wrapInQuotes takes string and returns a new string
// wrapped in single quotes
func wrapInQuotes(value string) string {
	return "'" + value + "'"
}

// joinByComma takes an array of strings
// and joins array values by comma separator
func joinByComma(flavors []string) string {
	return strings.Join(flavors[:], ", ")
}

// strainsHasId takes a strain's ID and array of strains
// returns true if ID is already in array
// otherwise it return false
func strainsHasID(id int, strains []Strain) bool {
	for _, strain := range strains {
		if strain.ID == id {
			return true
		}
	}
	return false
}

// getStrainIndexById takes a strain's ID and array of strains
// returns index and true if there is strain with specified ID in array
// otherwise it returns zero-index with false
func getStrainIndexByID(id int, strains []Strain) (int, bool) {
	for index, strain := range strains {
		if strain.ID == id {
			return index, true
		}
	}
	return 0, false
}

// strainHasFlavor takes a strain and string with flavor name
// returns true if flavors of strain contains specified flavor name
// otherwise it returns false
func strainHasFlavor(strain Strain, strainFlavor string) bool {
	for _, flavor := range strain.Flavors {
		if flavor == strainFlavor {
			return true
		}
	}
	return false
}

// strainHasEffect takes a strain and string with effect name
// returns true if effects of strain contains specified effect name in effect type
// otherwise it returns false
func strainHasEffect(strain Strain, strainEffectType string, strainEffectName string) bool {
	for _, effectNames := range strain.Effects {
		for _, effectName := range effectNames {
			if effectName == strainEffectName {
				return true
			}
		}
	}
	return false
}

// ValidRace take a race's name
// and checks it for validness (only "sativa", "indica" and "hybrid" values are available)
func ValidRace(raceName string) bool {
	raceName = strings.ToLower(raceName)
	switch raceName {
	case "sativa", "indica", "hybrid":
		return true
	default:
		return false
	}
}
