package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
	"log"
)

const (
	treeViewName    = "tree"
	requestViewName = "request"
	errorViewName   = "error"
)

type uiState struct {
	openCollection *postman.Collection
	dirty          bool
	variables      map[string]string
}

var state uiState

var itemTree *ItemTree

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

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func Open(collection, environment string) {
	pmColl, err := postman.ParseCollection(collection)
	if err != nil {
		log.Panicln(err)
		return
	}
	state = uiState{pmColl, false, make(map[string]string)}

	itemTree = NewItemTree(pmColl)

	if len(environment) > 0 {
		env, err := postman.ParseEnv(environment)
		if err == nil {
			fmt.Println("Loading environment", env.Name)
			for _, ev := range env.Values {
				if ev.Enabled {
					state.variables[ev.Key] = ev.Value
				}
			}
		}
	}
	Run()
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

	// golden-ish ratio
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

	mainWidth := 80

	treeX0 := 0
	treeY0 := 0
	treeX1 := maxX - mainWidth - 1
	treeY1 := maxY - 1

	requestX0 := maxX - mainWidth
	requestY0 := 0
	requestX1 := maxX - 1
	requestY1 := maxY - 1

	if treeView, err := g.SetView(treeViewName, treeX0, treeY0, treeX1, treeY1); err != nil {
		treeView.Title = "Tree"
		// treeView.Highlight = true
		// treeView.Autoscroll = true
		// treeView.SetCursor(0, 0)
		if err != gocui.ErrUnknownView {
			return err
		}
		if itemTree != nil {
			itemTree.Layout(treeView)
		}
	}
	if requestView, err := g.SetView(requestViewName, requestX0, requestY0, requestX1, requestY1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(requestView, "MAIN:  %d %d X %d %d", requestX0, requestY0, requestX1, requestY1)
	}
	return nil
}

func renderError(g *gocui.Gui, msg string) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("error", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, msg)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
