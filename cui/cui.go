package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
	"log"
)

type uiState struct {
	openCollection   *postman.Collection
	dirty            bool
	level, selection int
}

var state uiState

func Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(goldenLayout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	state = uiState{nil, false, 0, 0}

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
	state.openCollection = pmColl
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

/*
A layout based on the golden ratio sort of
*/
func goldenLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if state.openCollection == nil {
		renderError(g, "No File Open")
		return nil
	}

	// goldenish ratio
	// remainder := int(float64(maxX) - float64(maxY)*2.5)
	/*
	menuX0 := 0
	menuY0 := 0
	menuX1 := remainder - 1
	menuY1 := maxY - 1

	mainX0 := remainder
	mainY0 := 0
	mainX1 := maxX - 1
	mainY1 := maxY - 1
	 */

	menuX0 := 0
	menuY0 := 0
	menuX1 := maxX - 9
	menuY1 := maxY - 1

	mainX0 := maxX - 10
	mainY0 := 0
	mainX1 := maxX - 1
	mainY1 := maxY - 1

	if leftside, err := g.SetView("menu", menuX0, menuY0, menuX1, menuY1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		renderCollectionItems(leftside, state)
	}
	if mainView, err := g.SetView("main", mainX0, mainY0, mainX1, mainY1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(mainView, "MAIN:  %d %d X %d %d", mainX0, mainY0, mainX1, mainY1)
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

func renderCollectionItems(v *gocui.View, st uiState) {
	coll := st.openCollection
	fmt.Fprintf(v, "- %s - \n", coll.Info.Name)
	maxItemNameLength := 0
	for _, n := range state.openCollection.Items {
		if len(n.Name) > maxItemNameLength {
			maxItemNameLength = len(n.Name)
		}
	}
	maxItemNameLength = maxItemNameLength/2 + 1

	for pciIdx, pci := range coll.Items {
		fmt.Fprintf(v, " > %s\n", pci.Name)
		for childIdx, child := range pci.Children {
			printCollectionItem(v, &child, st.level == pciIdx && st.selection == childIdx)
		}
	}
}

func printCollectionItem(v *gocui.View, pci *postman.CollectionItem, selected bool) {
	if pci.Request != nil {
		fmt.Fprintf(v, "|%s|\t", pci.Request.Method)
	} else if len(pci.Children) > 0 {
		fmt.Fprintf(v, " > \t")
	} else {
		fmt.Fprint(v, "\t\t")
	}
	fmt.Fprintf(v, "\t%q %t", pci.Name, selected)
	fmt.Fprintln(v)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
