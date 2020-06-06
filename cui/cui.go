package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
	"log"
)

var openCollection *postman.Collection
var dirty bool

func Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(goldenLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func Open(filename string) {
	pmColl, err := postman.Parse(filename)
	if err != nil {
		log.Panicln(err)
		return
	}
	openCollection = pmColl
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(goldenLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

/*
A layout based on the golden ratio sort of
*/
func goldenLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	remainder := int(float64(maxX) - float64(maxY)*2.5)

	if openCollection == nil {
		renderError(g, "No File Open")
		return nil
	}

	if leftside, err := g.SetView("leftside", 0, 0, remainder-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// fmt.Fprintf(leftside, "LEFT:  %d %d X %d %d", 0, 0, remainder-1, maxY-1)
		renderCollectionItems(leftside, openCollection)
	}
	if mainView, err := g.SetView("main", remainder, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(mainView, "MAIN:  %d %d X %d %d", remainder, 0, maxX-1, maxY-1)
	}
	return nil
}

func renderError(g *gocui.Gui, msg string) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, msg)
	}
	return nil
}

func renderCollectionItems(v *gocui.View, pc *postman.Collection) {
	fmt.Fprintf(v, "- %s - \n", pc.Info.Name)
	maxItemNameLength := 0
	for _, n := range openCollection.Items {
		if len(n.Name) > maxItemNameLength {
			maxItemNameLength = len(n.Name)
		}
	}
	maxItemNameLength = maxItemNameLength/2 + 1

	for _, pci := range openCollection.Items {
		fmt.Fprintln(v, pci.Name)
		for _, ch := range pci.Children {
			if ch.Request != nil {
				fmt.Fprintf(v, "|%s|\t", ch.Request.Method)
			} else if len(ch.Children) > 0 {
				fmt.Fprintf(v, " > \t")
			} else {
				fmt.Fprint(v, "\t\t")
			}
			fmt.Fprintf(v, "\t%q", ch.Name)
			fmt.Fprintln(v)
		}
	}
}

func itemLayout(g *gocui.Gui) error {
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
