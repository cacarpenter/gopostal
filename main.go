package main

import (
	"fmt"
	"github.com/cacarpenter/gopostal/cui"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/cacarpenter/gopostal/util"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		cui.Run()
		return
	}

	cmd := os.Args[1]
	subargs := os.Args[2:]

	fmt.Println("running ", cmd, subargs)

	switch cmd {
	case "diff":
		runDiff(subargs)
	case "print", "show":
		printColl(subargs)
	case "open":
		cui.Open(subargs[0])
	default:
		fmt.Println("Unknown command", cmd)
	}
}

func printColl(subargs []string) {
	if len(subargs) < 1 {
		fmt.Println("gopostal print filename")
		return
	}
	coll, err := postman.Parse(subargs[0])
	if err != nil {
		panic(err)
	}
	postman.Print(coll)
}

func runDiff(subargs []string) {
	if len(subargs) < 2 {
		fmt.Println("gopostal diff filename1 filename2")
		return
	}
	coll1, err := postman.Parse(subargs[0])
	if err != nil {
		panic(err)
	}
	coll2, err := postman.Parse(subargs[1])
	if err != nil {
		panic(err)
	}

	tablelength := util.MaxLen(subargs[0], subargs[1])
	hline := util.StringOf('-', tablelength)
	blank := util.StringOf(' ', tablelength)
	space := "|        |"
	max := 0
	if len(coll1.Items) > len(coll2.Items) {
		max = len(coll1.Items)
	} else {
		max = len(coll2.Items)
	}
	fmt.Println(hline, space, hline)
	fmt.Println(subargs[0], space, subargs[1])
	fmt.Println(hline, space, hline)
	for i := 0; i < max; i++ {
		iname1 := blank
		if len(coll1.Items) > i {
			iname1 = util.Pad(coll1.Items[i].Name, tablelength)
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
		if len(coll2.Items) > i {
			iname2 = util.Pad(coll2.Items[i].Name, tablelength)
		}
		fmt.Println(iname1, space, iname2)
	}

	fmt.Println(len(coll1.Items), " vs ", len(coll2.Items))
}
