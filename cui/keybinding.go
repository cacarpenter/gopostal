package cui

import (
	"github.com/jroimartin/gocui"
	"log"
)

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'k', gocui.ModNone, cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'e', gocui.ModNone, toggleExpand); err != nil {
		log.Panicln(err)
	}

	return nil
}

func cursorMovement(st uiState) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		tv, err := g.View(treeViewName)
		if err != nil {
			return err
		}
		if itemTree != nil {
			itemTree.ArrowDown()
			itemTree.Layout(tv)
		}

		// ox, oy := tv.Origin()
		// cx, cy := tv.Cursor()
		// fmt.Fprintf(tv, "%d %d %d %d\n", ox, oy, cx, cy)
		return nil
	}
}

func updateTree(g *gocui.Gui, f func(it *ItemTree)) error {
	tv, err := g.View(treeViewName)
	if err != nil {
		return err
	}
	if itemTree != nil {
		f(itemTree)
		itemTree.Layout(tv)
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	return updateTree(g, func(it *ItemTree) {
		it.ArrowDown()
	})
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	return updateTree(g, func(it *ItemTree) {
		it.ArrowUp()
	})
}

func toggleExpand(g *gocui.Gui, v *gocui.View) error {
	return updateTree(g, func(it *ItemTree) {
		it.ToggleExpanded()
	})
}
