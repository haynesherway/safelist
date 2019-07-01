package main

import (
	"fmt"

	"github.com/haynesherway/safelist"
)

func main() {
	safeList := safelist.New()

	fmt.Println("Adding apple")
	safeList.Set("apple")
	fmt.Println("Adding banana")
	safeList.Set("banana")

	fmt.Println("Adding orange and kiwi")
	safeList.SetMultiple([]string{"orange", "kiwi"})

	if safeList.Get("apple") {
		fmt.Println("Map has apple")
	} else {
		fmt.Println("Map does not have apple")
	}

	if safeList.Get("corn") {
		fmt.Println("Map has corn")
	} else {
		fmt.Println("Map does not have corn")
	}

	if safeList.GetOrSet("grape") {
		fmt.Println("Map already had grape")
	} else {
		fmt.Println("Map did not have grape, but it does now")
	}

	fmt.Println(safeList.String())
}
