package cui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

func (ui *ConsoleUI) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, ui.cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, ui.cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, ui.cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'k', gocui.ModNone, ui.cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'e', gocui.ModNone, ui.toggleExpand); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'C', gocui.ModNone, ui.collapseAll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'E', gocui.ModNone, ui.expandAll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, ui.callRequest); err != nil {
		log.Panicln(err)
	}
	// DELETE
	if err := g.SetKeybinding("", gocui.KeyDelete, gocui.ModNone, ui.deleteSelection); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, ui.pageUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, ui.pageDown); err != nil {
		log.Panicln(err)
	}
	/*
		if err := g.SetKeybinding("", gocui.KeyF1, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
			gui.SetCurrentView(treeViewName)
			return nil
		}); err != nil {
			log.Panicln(err)
		}
	*/

	return nil
}

func (ui *ConsoleUI) updateTreeWidget(g *gocui.Gui, f func(*TreeWidget)) error {
	tv, err := g.View(treeViewName)
	if err != nil {
		ui.Logger.Println("Error getting tree view", err)
		return err
	}
	f(ui.treeWidget)
	ui.treeWidget.Layout(tv)
	rv, err := g.View(requestViewName)
	if err != nil {
		return err
	}
	if ui.treeWidget.selectedNode != nil {
		ui.requestWidget.request = ui.treeWidget.selectedNode.request
	}
	ui.requestWidget.Layout(rv)
	return nil
}

func (ui *ConsoleUI) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		fmt.Fprintln(v, "v is not null!")
	}
	if g != nil {
		cv := g.CurrentView()
		if cv != nil {
			fmt.Fprintln(cv, "this is the current non nil view")
		}
	}
	err := ui.updateTreeWidget(g, func(gw *TreeWidget) {
		gw.MoveDown()
	})
	if err != nil {
		return err
	}
	return nil
}

func (ui *ConsoleUI) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		fmt.Fprintf(v, "%d %d\n", cx, cy)
	}
	return ui.updateTreeWidget(g, func(gw *TreeWidget) {
		gw.MoveUp()
	})
}

func (ui *ConsoleUI) toggleExpand(g *gocui.Gui, v *gocui.View) error {
	return ui.updateTreeWidget(g, func(gw *TreeWidget) {
		gw.ToggleExpanded()
	})
}

func (ui *ConsoleUI) expandAll(g *gocui.Gui, v *gocui.View) error {
	return ui.updateTreeWidget(g, func(gw *TreeWidget) {
		gw.ExpandAll()
	})
}

func (ui *ConsoleUI) collapseAll(g *gocui.Gui, v *gocui.View) error {
	return ui.updateTreeWidget(g, func(gw *TreeWidget) {
		gw.CollapseAll()
	})
}

func (ui *ConsoleUI) callRequest(g *gocui.Gui, v *gocui.View) error {
	ui.Logger.Println("callRequest")
	ui.execFunc(v)

	return nil
}

func (ui *ConsoleUI) deleteSelection(g *gocui.Gui, v *gocui.View) error {
	ui.DeleteSelection()
	return nil
}

func (ui *ConsoleUI) pageUp(g *gocui.Gui, v *gocui.View) error {
	tv, err := g.View(treeViewName)
	if err != nil {
		return err
	}

	curx, cury := tv.Cursor()
	ui.Logger.Println("pageUp", curx, cury)
	if cury >= 5 {
		cury -= 5
	}
	ui.ScrollUp()
	tv.SetCursor(curx, cury)
	return nil
}
func (ui *ConsoleUI) pageDown(g *gocui.Gui, v *gocui.View) error {
	tv, err := g.View(treeViewName)
	if err != nil {
		return err
	}

	curx, cury := tv.Cursor()
	ui.Logger.Println("pageDown", curx, cury)
	ui.ScrollDown()
	tv.SetCursor(curx, cury+5)
	return nil
}
