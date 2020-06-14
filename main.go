package main

import (
	"flag"
	"fmt"
	"github.com/cacarpenter/gopostal/cui"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/cacarpenter/gopostal/util"
)

func main() {
	envFlag := flag.String("env", "local.postman_environment.json", "Specify postman environment")
	flag.Parse()

	if len(flag.Args()) < 2 {
		c := cui.ConsoleUI{}
		c.Run()
		return
	}

	fmt.Println(flag.Args())
	cmd := flag.Arg(0)
	subargs := flag.Args()[1:]

	fmt.Printf("running %q - %q\n", cmd, subargs)

	switch cmd {
	case "diff":
		runDiff(subargs)
	case "print", "show":
		printColl(subargs)
	case "open":
		c := cui.ConsoleUI{}
		c.Open(subargs[0], *envFlag)
	default:
		fmt.Println("Unknown command", cmd)
	}
}

func printColl(subargs []string) {
	if len(subargs) < 1 {
		fmt.Println("gopostal print filename")
		return
	}
	coll, err := postman.ParseCollection(subargs[0])
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
