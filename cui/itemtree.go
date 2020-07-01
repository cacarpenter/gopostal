package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
	"log"
)

// const SELECT_COLOR = colorGreen

type ItemTree struct {
	*log.Logger
	collections []*postman.Collection

	// which collection the cursor is currently on
	currentCollectionIdx int
	// which item in the tree is selected
	selected *postman.Collection
}

/*
func NewItemTree(pc *postman.Collection) *ItemTree {
	// the root item is the selected item
	it := ItemTree{pc, pc}
	return &it
}
 */

func (it *ItemTree) Layout(v *gocui.View) {
	v.Clear()
	for _, coll := range it.collections {
		maxItemNameLength := 0
		for _, n := range coll.Children {
			if len(n.Name) > maxItemNameLength {
				maxItemNameLength = len(n.Name)
			}
		}
		maxItemNameLength = maxItemNameLength/2 + 1
		// cx, cy := v.Cursor()
		// fmt.Fprintln(v, cx, cy)
		printCollection(v, "", coll)
	}
}

func (it *ItemTree) MoveUp() {
	it.Logger.Println("moveUp")
	var nextItem *postman.Collection
	prevSib := it.selected.PreviousSibling()
	if prevSib != nil {
		if prevSib.Expanded() {
			nextItem = prevSib.Children[len(prevSib.Children)-1]
		} else {
			nextItem = prevSib
		}
	} else if it.selected.Parent() != nil {
		if it.selected.Parent().Expanded() {
			nextItem = it.selected.Parent()
		} else {
			parentSib := it.selected.Parent().PreviousSibling()
			if parentSib != nil {
				nextItem = parentSib
			}
		}
	}
	if nextItem != nil {
		it.selected.SetSelected(false)
		nextItem.SetSelected(true)
		it.selected = nextItem
	} else {
		it.Logger.Println("MoveUp: No nextItem")
	}
}

func (it *ItemTree) MoveDown() {
	it.Logger.Println("moveDown")
	var nextItem *postman.Collection
	if it.selected == nil {
		it.Logger.Println("no current selection")
		return
	}
	if it.selected.Expanded() {
		if len(it.selected.Children) > 0 {
			nextItem = it.selected.Children[0]
		}
	} else {
		nextSib := it.selected.NextSibling()
		if nextSib != nil {
			nextItem = nextSib
		} else if it.selected.Parent() != nil {
			parentSib := it.selected.Parent().NextSibling()
			if parentSib != nil {
				nextItem = parentSib
			}
		}
	}
	if nextItem == nil {
		it.Logger.Println("Checking for another collection")
		// check for another collection
		if it.currentCollectionIdx+1 < len(it.collections) {
			it.currentCollectionIdx++
			it.selected = it.collections[it.currentCollectionIdx]
		}
	}
	if nextItem != nil {
		it.Logger.Println("Setting next item to ", nextItem.Label())
		it.selected.SetSelected(false)
		nextItem.SetSelected(true)
		it.selected = nextItem
	} else {
		it.Logger.Println("moveDown: No nextItem")
	}
}

func (it *ItemTree) ExpandAll() {
	for _, pmColl := range it.collections {
		pmColl.Expand(true)
	}
}

func (it *ItemTree) ToggleExpanded() {
	it.Logger.Println("ToggleExpanded")
	if it.selected != nil {
		it.selected.ToggleExpanded()
	}
}

func (it *ItemTree) SetCollections(pcs []*postman.Collection) {
	it.collections = pcs
	if len(it.collections) > 0 {
		it.selected = it.collections[0]
	}
}

func printCollection(v *gocui.View, pad string, pci *postman.Collection) {
	fmt.Fprint(v, pad)
	n := pci.Name
	if len(n) == 0 && pci.Info != nil {
		n = pci.Info.Name
	}
	if pci.Request != nil {
		fmt.Fprint(v, "[", colorCyan, pci.Request.Method, colorReset, "] ")
		if pci.Selected() {
			fmt.Fprintln(v, colorGreen, pci.Name, colorReset)
		} else {
			fmt.Fprintln(v, pci.Name)
		}
	} else if len(pci.Children) > 0 {
		chev := "> "
		if pci.Expanded() {
			chev = "\\/"
		}
		label := fmt.Sprintf("%s %s", chev, n)
		if pci.Selected() {
			fmt.Fprintln(v, colorGreen, label, colorReset)
		} else {
			fmt.Fprintln(v, label)
		}
		if pci.Expanded() {
			for _, child := range pci.Children {
				printCollection(v, pad+" ", child)
			}
		}
	}
}

