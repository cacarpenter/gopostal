package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
	"log"
)

var openCollection *postman.Collection

func Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

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

	g.SetManagerFunc(itemLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func itemLayout(g *gocui.Gui) error {
	if openCollection == nil {
		return fmt.Errorf("No open collection")
	}
	//	names := openCollection.Items
	yHeight := len(openCollection.Items)/2 + 1
	maxX, maxY := g.Size()
	if v, err := g.SetView("items", maxX/2-12, maxY/2-yHeight, maxX/2+12, maxY/2+yHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		for _, pci := range openCollection.Items {
			fmt.Fprintln(v, pci.Name)
		}
	}
	return nil
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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
