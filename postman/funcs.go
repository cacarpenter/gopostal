package postman

import "fmt"

func printCollectionItem(pad string, ci *Collection) {
	fmt.Println(pad+"* Name :", ci.Name)
	if ci.Request != nil {
		fmt.Println(pad+"* Req  : ", ci.Request)
	}
	for _, child := range ci.Children {
		printCollectionItem(pad+pad, child)
	}
}

func Print(coll *Collection) {
	fmt.Println("Name:", coll.Info.Name)
	for _, i := range coll.Children {
		printCollectionItem("  ", i)
	}
}
