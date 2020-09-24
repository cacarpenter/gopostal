package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
	"github.com/cacarpenter/gopostal/util"
	"github.com/jroimartin/gocui"
	"io"
	"log"
)

const (
	treeViewName      = "tree"
	requestViewName   = "request"
	responseViewName  = "response"
	errorViewName     = "error"
	variablesViewName = "variables"
	modalViewName     = "modal"

	// font coloring
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"

	defaultScroll = 5
)

// ConsoleUI contains the various widgets to run the console based user interface
//
type ConsoleUI struct {
	treeWidget      *TreeWidget
	requestWidget   *RequestWidget
	variablesWidget *VariablesWidget
	responseWidget  *ResponseWidget

	modalVisible bool

	*log.Logger

	execFunc func(w io.Writer)

	// left right split that is movable
	verticalSplitX int
}

// NewConsoleUI creates a new ConsoleUI instance
func NewConsoleUI(logger *log.Logger) *ConsoleUI {
	ui := ConsoleUI{}
	ui.treeWidget = new(TreeWidget)
	ui.requestWidget = new(RequestWidget)
	ui.responseWidget = new(ResponseWidget)
	ui.variablesWidget = new(VariablesWidget)
	ui.Logger = logger
	ui.treeWidget.Logger = logger
	ui.verticalSplitX = 30
	return &ui
}

// Run display the console UI
func (ui *ConsoleUI) Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(ui.layout)

	if err := ui.keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (ui *ConsoleUI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	remainder := ui.verticalSplitX
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

	responseX0 := remainder
	responseY0 := maxY - (maxY / 4)
	responseX1 := maxX - 1
	responseY1 := maxY - 1

	modalX0 := 10
	modalY0 := 10
	modalX1 := maxX - 10
	modalY1 := maxY - 10

	if treeView, err := g.SetView(treeViewName, treeX0, treeY0, treeX1, treeY1); err != nil {
		treeView.Title = "Tree"
		treeView.Highlight = false
		treeView.Autoscroll = false
		treeView.SetCursor(0, 0)
		if err != gocui.ErrUnknownView {
			return err
		}
		ui.treeWidget.Layout(treeView)
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
	if responseView, err := g.SetView(responseViewName, responseX0, responseY0, responseX1, responseY1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		responseView.Title = "Response"
	}
	if ui.modalVisible {
		if modalView, err := g.SetView(modalViewName, modalX0, modalY0, modalX1, modalY1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			modalView.Title = "MOD"
			modalView.SelBgColor = gocui.ColorMagenta
			ui.variablesWidget.Layout(modalView)
			g.SetCurrentView(modalViewName)
		}
	} else {
		g.DeleteView(modalViewName)
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

// SetRequestGroups sets the groups to display in the UI
func (ui *ConsoleUI) SetRequestGroups(grps []*gpmodel.RequestGroup) {
	ui.treeWidget.SetRequestGroups(grps)
}

// UpdateVariables changes the variables currently displayed
func (ui *ConsoleUI) UpdateVariables(vars map[string]string) {
	ui.variablesWidget.SetVariables(util.Map2Array(vars))
}

// UpdateVariable update a variable in the display
func (ui *ConsoleUI) UpdateVariable(k, v string) {
	ui.variablesWidget.SetVariable(k, v)
}

// SetOnExec which function to call when a node is executed
func (ui *ConsoleUI) SetOnExec(f func(w io.Writer)) {
	ui.execFunc = f
}

// DeleteSelection
func (ui *ConsoleUI) DeleteSelection(g *gocui.Gui, v *gocui.View) error {
	ui.Logger.Println("Delete not yet implemented")
	return nil
}

// ArrowLeft
func (ui *ConsoleUI) ArrowLeft(g *gocui.Gui, v *gocui.View) error {
	ui.verticalSplitX--
	return nil
}

// ArrowRight
func (ui *ConsoleUI) ArrowRight(g *gocui.Gui, v *gocui.View) error {
	ui.verticalSplitX++
	return nil
}

func (ui *ConsoleUI) ToggleVariablesModal(g *gocui.Gui, v *gocui.View) error {
	ui.modalVisible = !ui.modalVisible
	ui.Logger.Printf("TVM %t\n", ui.modalVisible)
	return ui.layout(g)
}

func (ui *ConsoleUI) ScrollUp(g *gocui.Gui, v *gocui.View) error {
	tv, err := g.View(treeViewName)
	if err != nil {
		return err
	}

	origX, origY := tv.Origin()
	if origY >= defaultScroll {
		origY -= defaultScroll
	}
	return tv.SetOrigin(origX, origY)
}

func (ui *ConsoleUI) ScrollDown(g *gocui.Gui, v *gocui.View) error {
	tv, err := g.View(treeViewName)
	if err != nil {
		return err
	}

	origX, origY := tv.Origin()
	return tv.SetOrigin(origX, origY+defaultScroll)
}

func (ui *ConsoleUI) IsRequestSelected() bool {
	return ui.treeWidget.selectedNode != nil && ui.treeWidget.selectedNode.request != nil
}

func (ui *ConsoleUI) SelectedRequest() *gpmodel.RequestSpec {
	return ui.treeWidget.selectedNode.request
}
