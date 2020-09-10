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
	treeWidget      *TreeWidget
	requestWidget   *RequestWidget
	variablesWidget *VariablesWidget
	responseWidget  *ResponseWidget
	*log.Logger
	execFunc func(w io.Writer)
}

func NewConsoleUI(logger *log.Logger) *ConsoleUI {
	ui := ConsoleUI{}
	ui.treeWidget = new(TreeWidget)
	ui.requestWidget = new(RequestWidget)
	ui.responseWidget = new(ResponseWidget)
	ui.variablesWidget = new(VariablesWidget)
	ui.Logger = logger
	ui.treeWidget.Logger = logger
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
A layout based on the golden ratio sort of
*/
func (ui *ConsoleUI) goldenLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

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

	responseX0 := remainder
	responseY0 := maxY - (maxY / 4)
	responseX1 := maxX - 1
	responseY1 := maxY - 1

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

func (ui *ConsoleUI) SetGroups(grps []*gpmodel.Group) {
	ui.treeWidget.SetGroups(grps)
}

func (ui *ConsoleUI) UpdateVariables(vars map[string]string) {
	ui.variablesWidget.SetVariables(util.Map2Array(vars))
}

func (ui *ConsoleUI) UpdateVariable(k, v string) {
	ui.variablesWidget.SetVariable(k, v)
}

func (ui *ConsoleUI) SetOnExec(f func(w io.Writer)) {
	ui.execFunc = f
}

func (ui *ConsoleUI) DeleteSelection() {
}

func (ui *ConsoleUI) ScrollUp() {
	ui.Logger.Println("ScrollUp")
}

func (ui *ConsoleUI) ScrollDown() {
	ui.Logger.Println("ScrollDown")
}

func (ui *ConsoleUI) IsRequestSelected() bool {
	return ui.treeWidget.selectedNode != nil && ui.treeWidget.selectedNode.request != nil
}

func (ui *ConsoleUI) SelectedRequest() *gpmodel.RequestSpec {
	return ui.treeWidget.selectedNode.request
}
