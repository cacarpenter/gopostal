package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gp"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
	"log"
)

const (
	treeViewName      = "tree"
	requestViewName   = "request"
	debugViewName     = "debug"
	errorViewName     = "error"
	variablesViewName = "variables"

	// font coloring
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

type ConsoleUI struct {
	itemTree        *ItemTree
	requestWidget   *RequestWidget
	variablesWidget *VariablesWidget
}

func (ui *ConsoleUI) Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(ui.goldenLayout)

	if err := ui.keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (ui *ConsoleUI) Open(collection, environment string) {
	pmColl, err := postman.ParseCollection(collection)
	if err != nil {
		log.Panicln(err)
		return
	}

	// show all the root items by default
	pmColl.ToggleExpanded()
	ui.itemTree = NewItemTree(pmColl)
	ui.requestWidget = &RequestWidget{pmColl}
	ui.variablesWidget = &VariablesWidget{}

	if len(environment) > 0 {
		env, err := postman.ParseEnv(environment)
		if err == nil {
			sess := gp.CurrentSession()
			sess.Update(env, true)
		} else {
			fmt.Println("Cant load env", err)
		}
	}
	ui.Run()
}

/*
A layout based on the golden ratio sort of
*/
func (ui *ConsoleUI) goldenLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	/*
		if state.openCollection == nil {
			renderError(g, "No File Open")
			return nil
		}
	*/

	// golden-ish ratio
	remainder := int(float64(maxX) - float64(maxY)*2.5)
	treeX0 := 0
	treeY0 := 0
	treeX1 := remainder - 1
	treeY1 := maxY - 1

	requestX0 := remainder
	requestY0 := 0
	requestX1 := maxX - 1
	requestY1 := maxY - 1

	// test values
	/*
		mainWidth := 100

		treeX0 := 0
		treeY0 := 0
		treeX1 := maxX - mainWidth - 1
		treeY1 := maxY - 1

		requestX0 := maxX - mainWidth
		requestY0 := 0
		requestX1 := maxX - 1
		requestY1 := maxY/2 - 1
	*/

	variablesX0 := maxX - remainder
	variablesY0 := maxY / 2
	variablesX1 := maxX - 1
	variablesY1 := maxY - 1

	if treeView, err := g.SetView(treeViewName, treeX0, treeY0, treeX1, treeY1); err != nil {
		treeView.Title = "Tree"
		// treeView.Highlight = true
		// treeView.Autoscroll = true
		// treeView.SetCursor(0, 0)
		if err != gocui.ErrUnknownView {
			return err
		}
		ui.itemTree.Layout(treeView)
	}
	if requestView, err := g.SetView(requestViewName, requestX0, requestY0, requestX1, requestY1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		requestView.Title = "Request"
		fmt.Fprintf(requestView, "MAIN:  %d %d X %d %d", requestX0, requestY0, requestX1, requestY1)
	}
	if variablesView, err := g.SetView(variablesViewName, variablesX0, variablesY0, variablesX1, variablesY1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		variablesView.Title = "Variables"
		// fmt.Fprintf(variablesView, "v:  %d %d X %d %d", variablesX0, variablesY0, variablesX1, variablesY1)
		ui.variablesWidget.Layout(variablesView)
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
