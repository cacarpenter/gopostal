package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
)

// const SELECT_COLOR = colorGreen

type ItemTree struct {
	collection *postman.Collection

	// which item in the tree is selected
	selectedItem *postman.Collection
}

func NewItemTree(pc *postman.Collection) *ItemTree {
	// the root item is the selected item
	it := ItemTree{pc, pc}
	return &it
}

func (it *ItemTree) Layout(v *gocui.View) {
	coll := it.collection
	maxItemNameLength := 0
	for _, n := range it.collection.Children {
		if len(n.Name) > maxItemNameLength {
			maxItemNameLength = len(n.Name)
		}
	}
	maxItemNameLength = maxItemNameLength/2 + 1

	v.Clear()
	printCollection(v, "", coll)
}

func (it *ItemTree) MoveUp() {
	var nextItem *postman.Collection
	prevSib := it.selectedItem.PreviousSibling()
	if prevSib != nil {
		if prevSib.Expanded() {
			nextItem = prevSib.Children[len(prevSib.Children)-1]
		} else {
			nextItem = prevSib
		}
	} else if it.selectedItem.Parent() != nil {
		if it.selectedItem.Parent().Expanded() {
			nextItem = it.selectedItem.Parent()
		} else {
			parentSib := it.selectedItem.Parent().PreviousSibling()
			if parentSib != nil {
				nextItem = parentSib
			}
		}
	}
	if nextItem != nil {
		it.selectedItem.SetSelected(false)
		nextItem.SetSelected(true)
		it.selectedItem = nextItem
	}
}

func (it *ItemTree) MoveDown() {
	var nextItem *postman.Collection
	if it.selectedItem.Expanded() {
		if len(it.selectedItem.Children) > 0 {
			nextItem = it.selectedItem.Children[0]
		}
	} else {
		nextSib := it.selectedItem.NextSibling()
		if nextSib != nil {
			nextItem = nextSib
		} else if it.selectedItem.Parent() != nil {
			parentSib := it.selectedItem.Parent().NextSibling()
			if parentSib != nil {
				nextItem = parentSib
			}
		}
	}
	if nextItem != nil {
		it.selectedItem.SetSelected(false)
		nextItem.SetSelected(true)
		it.selectedItem = nextItem
	}
}

func (it *ItemTree) ToggleExpanded() {
	if it.selectedItem != nil {
		it.selectedItem.ToggleExpanded()
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
