package main

import (
	"encoding/json"
	"fmt"
	"github.com/cacarpenter/gopostal/gp"
	"github.com/cacarpenter/gopostal/util"
	"io/ioutil"
	"os"
)

const SCHEMA_2_1_0 = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("gopostal command")
		return
	}

	cmd := os.Args[1]
	subargs := os.Args[2:]

	fmt.Println("running ", cmd, subargs)

	switch cmd {
	case "diff":
		runDiff(subargs)
	default:
		fmt.Println("Unknown command", cmd)
	}
	/*

		fmt.Println("Name:", coll.Info.Name)
		for _, i := range coll.Item {
			fmt.Println(i.Name)
			for _, i2 := range i.Item {
				fmt.Println("\t -", i2.Name)
				fmt.Println("\t *", i2.Request)
			}
		}
	*/
}

func parseCollection(filename string) (*gp.PostmanCollection, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var coll gp.PostmanCollection
	err = json.Unmarshal(data, &coll)
	if err != nil {
		return nil, err
	}
	return &coll, nil
}

func runDiff(subargs []string) {
	if len(subargs) < 2 {
		fmt.Println("gopostal diff filename1 filename2")
		return
	}
	coll1, err := parseCollection(subargs[0])
	if err != nil {
		panic(err)
	}
	coll2, err := parseCollection(subargs[1])
	if err != nil {
		panic(err)
	}

	tablelength := util.MaxLen(subargs[0], subargs[1])
	hline := util.StringOf('-', tablelength)
	blank := util.StringOf(' ', tablelength)
	space := "|        |"
	max := 0
	if len(coll1.Item) > len(coll2.Item) {
		max = len(coll1.Item)
	} else {
		max = len(coll2.Item)
	}
	fmt.Println(hline, space, hline)
	fmt.Println(subargs[0], space, subargs[1])
	fmt.Println(hline, space, hline)
	for i := 0; i < max; i++ {
		iname1 := blank
		if len(coll1.Item) > i {
			iname1 = util.Pad(coll1.Item[i].Name, tablelength)
			/*
				if len(coll1.Item[i].Name) < tablelength {
					fill := util.StringOf(' ', len(coll1.Item[i].Name)-tablelength)
					iname1 = coll1.Item[i].Name + fill
				} else {
					iname1 = coll1.Item[i].Name
				}
			*/
		}
		iname2 := blank
		if len(coll2.Item) > i {
			iname2 = util.Pad(coll2.Item[i].Name, tablelength)
		}
		fmt.Println(iname1, space, iname2)
	}

	fmt.Println(len(coll1.Item), " vs ", len(coll2.Item))
}
