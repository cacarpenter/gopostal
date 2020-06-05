package postman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Parse(filename string) (*Collection, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var coll Collection
	err = json.Unmarshal(data, &coll)
	if err != nil {
		return nil, err
	}
	return &coll, nil
}

func printCollectionItem(pad string, ci *CollectionItem) {
	fmt.Println(pad+"* Name :", ci.Name)
	if ci.Request != nil {
		fmt.Println(pad+"* Req  : ", ci.Request)
	}
	for _, child := range ci.Children {
		printCollectionItem(pad+pad, &child)
	}
}

func Print(coll *Collection) {
	fmt.Println("Name:", coll.Info.Name)
	for _, i := range coll.Items {
		printCollectionItem("  ", &i)
	}
}
