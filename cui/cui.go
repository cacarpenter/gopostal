package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/cacarpenter/gopostal/util"
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
	*log.Logger
	execFunc func()
}

func NewConsoleUI(logger *log.Logger) *ConsoleUI {
	ui := ConsoleUI{}
	ui.itemTree = &ItemTree{}
	ui.requestWidget = &RequestWidget{}
	ui.variablesWidget = &VariablesWidget{}
	ui.Logger = logger
	ui.itemTree.Logger = logger
	return &ui
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

/*
func (ui *ConsoleUI) Open(collection, environment string) {
	pmColl, err := postman.ParseCollection(collection)
	if err != nil {
		log.Panicln(err)
	}

	// show all the root items by default
	pmColl.ToggleExpanded()
	ui.Init(pmColl)

	if len(environment) > 0 {

	}
	ui.init()
}*/

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
	leftWidthAdd := 6

	// collection view
	treeX0 := 0
	treeY0 := 0
	treeX1 := remainder - 1
	treeY1 := maxY/2 + leftWidthAdd - 1

	variablesX0 := 0
	variablesY0 := maxY/2 + leftWidthAdd
	variablesX1 := remainder - 1
	variablesY1 := maxY - 1

	requestX0 := remainder
	requestY0 := 0
	requestX1 := maxX - 1
	requestY1 := maxY - (maxY / 4) - 1

	debugX0 := remainder
	debugY0 := maxY - (maxY / 4)
	debugX1 := maxX - 1
	debugY1 := maxY - 1

	if treeView, err := g.SetView(treeViewName, treeX0, treeY0, treeX1, treeY1); err != nil {
		treeView.Title = "Tree"
		treeView.Highlight = true
		treeView.Autoscroll = false
		// treeView.SetCursor(0, 0)
		if err != gocui.ErrUnknownView {
			return err
		}
		// investigate how this is happening
		if ui.itemTree != nil {
			ui.itemTree.Layout(treeView)
		}
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
	if debugView, err := g.SetView(debugViewName, debugX0, debugY0, debugX1, debugY1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		debugView.Title = "Debug"
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

func (ui *ConsoleUI) SetPostmanCollections(pcs []*postman.Collection) {
	ui.itemTree.SetCollections ( pcs )
	for i, pc := range pcs {
		if i == 0 {
			ui.itemTree.selected = pc
			pc.SetSelected(true)
		}
		pc.ToggleExpanded()
	}
}

func (ui *ConsoleUI) UpdateVariables(vars map[string]string) {
	ui.variablesWidget.SetVariables(util.Map2Array(vars))
}

func (ui *ConsoleUI) SetOnExec( f func()) {
	ui.execFunc = f
}

func (ui *ConsoleUI) SelectedCollection() *postman.Collection {
	return ui.itemTree.selected
}

func (ui *ConsoleUI) DeleteSelection() {
}