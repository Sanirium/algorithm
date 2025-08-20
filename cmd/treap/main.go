package main

import (
	myTreap "algorithms/internal/treap"
	"fmt"
)

func main() {
	treap := myTreap.NewTreap[string]()

	treap.Insert("Beer", 95)
	treap.Insert("Bacon", 77)
	treap.Insert("Eggs", 129)
	treap.Insert("Pork", 56)
	treap.Insert("Milk", 55)
	treap.Insert("Flour", 10)
	treap.Insert("Water", 32)
	treap.Insert("Butter", 76)

	fmt.Println(treap.String())
}
