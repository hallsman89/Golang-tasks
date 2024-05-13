package comparer

import (
	"fmt"
	reader "s21compare/ex00/db-reader"
)

func CompareDB(oldDB, newDB reader.DataBase) {
	for _, oldCake := range oldDB.Cakes {
		found := false
		for _, newCake := range newDB.Cakes {
			if oldCake.Name == newCake.Name {
				found = true
				if oldCake.Time != newCake.Time {
					fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", oldCake.Name, newCake.Time, oldCake.Time)
				}
				for _, oldIngredient := range oldCake.Ingredients {
					ingredientFound := false
					for _, newIngredient := range newCake.Ingredients {
						if oldIngredient.Name == newIngredient.Name {
							ingredientFound = true
							if oldIngredient.Count != newIngredient.Count {
								fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", oldIngredient.Name, oldCake.Name, newIngredient.Count, oldIngredient.Count)
							}
							if oldIngredient.Unit != newIngredient.Unit {
								fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", oldIngredient.Name, oldCake.Name, newIngredient.Unit, oldIngredient.Unit)
							}
							break
						}
					}
					if !ingredientFound {
						fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", oldIngredient.Name, oldCake.Name)
					}
				}
				for _, newIngredient := range newCake.Ingredients {
					ingredientFound := false
					for _, oldIngredient := range oldCake.Ingredients {
						if oldIngredient.Name == newIngredient.Name {
							ingredientFound = true
							break
						}
					}
					if !ingredientFound {
						fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", newIngredient.Name, oldCake.Name)
					}
				}
				break
			}
		}
		if !found {
			fmt.Printf("REMOVED cake \"%s\"\n", oldCake.Name)
		}
	}
	for _, newCake := range newDB.Cakes {
		found := false
		for _, oldCake := range oldDB.Cakes {
			if oldCake.Name == newCake.Name {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("ADDED cake \"%s\"\n", newCake.Name)
		}
	}
}
