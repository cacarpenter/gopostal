package cui

import (
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
	if err := g.SetKeybinding("", gocui.KeyDelete, gocui.ModNone, ui.DeleteSelection); err != nil {
		log.Panicln(err)
	}
	// PAGING
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, ui.ScrollUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, ui.ScrollDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, ui.ArrowLeft); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, ui.ArrowRight); err != nil {
		log.Panicln(err)
	}

	// VARIABLES MODAL
	if err := g.SetKeybinding("", 'v', gocui.ModNone, ui.ToggleVariablesModal); err != nil {
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

	if err := g.SetKeybinding("", gocui.KeyF2, gocui.ModNone, ui.SaveCollection); err != nil {
		log.Panicln(err)
	}

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
	// shouldn't happen, there should always be a current view
	if v == nil {
		return nil
	}
	ui.Logger.Println("cursorDown: Got a view", v.Name())
	if v.Name() == treeViewName {
		return ui.updateTreeWidget(g, func(gw *TreeWidget) {
			gw.MoveDown()
		})
	}
	return nil
}

func (ui *ConsoleUI) cursorUp(g *gocui.Gui, v *gocui.View) error {
	// shouldn't happen, there should always be a current view
	if v == nil {
		return nil
	}
	if v.Name() == treeViewName {
		return ui.updateTreeWidget(g, func(gw *TreeWidget) {
			gw.MoveUp()
		})
	}
	return nil
}

func (ui *ConsoleUI) toggleExpand(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	if v.Name() == treeViewName {
		return ui.updateTreeWidget(g, func(gw *TreeWidget) {
			gw.ToggleExpanded()
		})
	}
	return nil
}

func (ui *ConsoleUI) expandAll(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	if v.Name() == treeViewName {
		return ui.updateTreeWidget(g, func(gw *TreeWidget) {
			gw.ExpandAll()
		})
	}
	return nil
}

func (ui *ConsoleUI) collapseAll(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	if v.Name() == treeViewName {
		return ui.updateTreeWidget(g, func(gw *TreeWidget) {
			gw.CollapseAll()
		})
	}
	return nil
}

func (ui *ConsoleUI) callRequest(g *gocui.Gui, v *gocui.View) error {
	ui.Logger.Println("callRequest")
	g.Update(func(g2 *gocui.Gui) error {
		responseView, err := g2.View(responseViewName)
		if err != nil {
			ui.Logger.Println("ERROR: No response view found")
			return err
		}
		ui.execFunc(responseView)
		varsView, err := g2.View(variablesViewName)
		if err != nil {
			ui.Logger.Printf("ERROR: No view %s found: %q\n", variablesViewName, err)
			return err
		}
		ui.variablesWidget.Layout(varsView)
		return nil
	})
	return nil
}
