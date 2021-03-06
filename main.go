package main

import (
	"flag"
	"fmt"
	"github.com/cacarpenter/gopostal/gp"
	"github.com/cacarpenter/gopostal/gpmodel"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/cacarpenter/gopostal/util"
	"log"
)

func main() {
	// envFlag := flag.String("env", "local.postman_environment.json", "Specify postman environment")
	flag.Parse()

	if len(flag.Args()) < 0 {
		fmt.Println("Need to specify a file(s) for now")
		return
	}

	app := gp.New()

	var varGroups []*gpmodel.VarGroup
	var reqGroups []*gpmodel.RequestGroup
	for _, filename := range flag.Args() {
		if postman.IsEnvironmentFile(filename) {
			vg, err := postman.ParseEnvironment(filename)
			if err == nil {
				varGroups = append(varGroups, vg)
			} else {
				log.Panicln("Bad postman env", filename, err)
			}
		} else if postman.IsCollectionFile(filename) {
			pmColl, err := postman.ParseCollection(filename)
			if err == nil {
				reqGroups = append(reqGroups, pmColl)
			} else {
				log.Panicln("Bad postman collection", filename, err)
			}
		}
	}

	app.SetRequestGroups(reqGroups)
	app.SetVarGroups(varGroups)
	app.Run()
	app.Stop()
}

func runDiff(subargs []string) {
	if len(subargs) < 2 {
		log.Println("gopostal diff filename1 filename2")
		return
	}
	coll1, err := postman.ParseCollection(subargs[0])
	if err != nil {
		panic(err)
	}
	coll2, err := postman.ParseCollection(subargs[1])
	if err != nil {
		panic(err)
	}

	tablelength := util.MaxLen(subargs[0], subargs[1])
	hline := util.StringOf('-', tablelength)
	blank := util.StringOf(' ', tablelength)
	space := "|        |"
	max := 0
	if len(coll1.Children) > len(coll2.Children) {
		max = len(coll1.Children)
	} else {
		max = len(coll2.Children)
	}
	fmt.Println(hline, space, hline)
	fmt.Println(subargs[0], space, subargs[1])
	fmt.Println(hline, space, hline)
	for i := 0; i < max; i++ {
		iname1 := blank
		if len(coll1.Children) > i {
			iname1 = util.Pad(coll1.Children[i].Name, tablelength)
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
		if len(coll2.Children) > i {
			iname2 = util.Pad(coll2.Children[i].Name, tablelength)
		}
		fmt.Println(iname1, space, iname2)
	}

	fmt.Println(len(coll1.Children), " vs ", len(coll2.Children))
}
