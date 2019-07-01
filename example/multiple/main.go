package main

import (
	"fmt"

	"github.com/haynesherway/safelist"
)

func main() {
	safeList := safelist.New()

	fmt.Println("Adding apple to fruit list")
	safeList.SetList("fruit", "apple")
	fmt.Println("Adding banana to fruit list")
	safeList.SetList("fruit", "banana")
	fmt.Println("Adding orange to fruit list")
	safeList.SetList("fruit", "orange")

	fmt.Println("Adding carrot to vegetable list")
	safeList.SetList("vegetable", "carrot")

	fruits := safeList.GetList("fruit")
	fmt.Println("Fruits:", fruits)

	fmt.Println("Adding celery to vegetables if it is not already added")
	if safeList.GetOrSetList("vegetable", "celery") {
		fmt.Println("Celery was already in vegetable")
	} else {
		fmt.Println("Celery was added to vegetable")
	}

	fmt.Println("Adding cucumber to vegetables")
	vegs, existed := safeList.GetAndSetList("vegetable", "cucumber")
	if existed {
		fmt.Println("Cucumber already in vegetables")
	} else {
		fmt.Println("Cucumber added to vegetables")
	}
	fmt.Println("Vegetables:", vegs)

	fmt.Println("Deleting banana from fruit list")
	safeList.DeleteListValue("fruit", "banana")
	fmt.Println(safeList.String())

	fmt.Println("Deleting vegetable list")
	safeList.DeleteList("vegetable")
	fmt.Println(safeList.String())

}
