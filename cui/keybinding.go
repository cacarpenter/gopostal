package cui

import (
	"bytes"
	"github.com/cacarpenter/gopostal/gp"
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

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, ui.callRequest); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (ui *ConsoleUI) updateTree(g *gocui.Gui, f func(it *ItemTree)) error {
	tv, err := g.View(treeViewName)
	if err != nil {
		return err
	}
	if ui.itemTree != nil {
		f(ui.itemTree)
		ui.itemTree.Layout(tv)
	}
	rv, err := g.View(requestViewName)
	if err != nil {
		return err
	}
	ui.requestWidget.collection = ui.itemTree.selectedItem
	ui.requestWidget.Layout(rv)
	return nil
}

func (ui *ConsoleUI) cursorDown(g *gocui.Gui, v *gocui.View) error {
	err := ui.updateTree(g, func(it *ItemTree) {
		it.MoveDown()
	})
	if err != nil {
		return err
	}
	return nil
}

func (ui *ConsoleUI) cursorUp(g *gocui.Gui, v *gocui.View) error {
	return ui.updateTree(g, func(it *ItemTree) {
		it.MoveUp()
	})
}

func (ui *ConsoleUI) toggleExpand(g *gocui.Gui, v *gocui.View) error {
	return ui.updateTree(g, func(it *ItemTree) {
		it.ToggleExpanded()
	})
}

func (ui *ConsoleUI) callRequest(g *gocui.Gui, v *gocui.View) error {
	debug, err := g.View(debugViewName)
	if err != nil {
		return err
	}
	// TODO move this logic somewhere else
	if ui.itemTree.selectedItem != nil && ui.itemTree.selectedItem.Request != nil {
		response, err := gp.CallRequest(ui.itemTree.selectedItem.Request, debug)

		if err != nil {
			return nil
		}

		for _, ev := range ui.itemTree.selectedItem.Events {
			var buf bytes.Buffer
			for _, l := range ev.Script.Lines {
				buf.WriteString(l)
				buf.WriteString("\n")
			}
			gp.RunJavaScript(buf.String(), *response, debug)
		}
	}
	return nil
}

