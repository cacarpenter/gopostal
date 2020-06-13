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

	// state
	selection int
}

func NewItemTree(pc *postman.Collection) *ItemTree {
	// the root item is the selected item
	it := ItemTree{pc, pc, 0 }
	return &it
}

func (it *ItemTree) Layout(v *gocui.View) {
	coll := it.collection
	maxItemNameLength := 0
	for _, n := range state.openCollection.Children {
		if len(n.Name) > maxItemNameLength {
			maxItemNameLength = len(n.Name)
		}
	}
	maxItemNameLength = maxItemNameLength/2 + 1

	v.Clear()
	fmt.Fprintf(v, " ~ %d ~\n", it.selection)
	printCollection(v, "", coll)
}

func (it *ItemTree) ArrowUp() {
	prevSib := it.selectedItem.PreviousSibling()
	if prevSib != nil {
		it.selectedItem.SetSelected(false)
		prevSib.SetSelected(true)
		it.selectedItem = prevSib
	}
}

func (it *ItemTree) ArrowDown() {
	if it.selectedItem.Expanded() {
		if len(it.selectedItem.Children) > 0 {
			it.selectedItem.SetSelected(false)
			firstChild := it.selectedItem.Children[0]
			firstChild.SetSelected(true)
			it.selectedItem = firstChild
		}
	} else {
		nextSib := it.selectedItem.NextSibling()
		if nextSib != nil {
			it.selectedItem.SetSelected(false)
			nextSib.SetSelected(true)
			it.selectedItem = nextSib
		}
	}
}

func (it *ItemTree) ToggleExpanded() {
	if it.selectedItem != nil {
		it.selectedItem.ToggleExpanded()
	}
}

func (it *ItemTree) currentCollectionItem() *postman.Collection {
	//var selected *postman.CollectionItem
	selectedCi := it.collection.Children[0]
	for i := 0; i < it.selection; i++ {
		//		selected
		//	if rootPci.Expanded() {}
	}
	return selectedCi
}

func printCollection(v *gocui.View, pad string, pci *postman.Collection) {
	if pci.Selected() {
		fmt.Fprint(v, colorPurple, "*", colorReset)
	} else {
		fmt.Fprint(v, " ")
	}
	fmt.Fprint(v, pad)
	n := pci.Name
	if len(n) == 0 && pci.Info != nil {
		n = pci.Info.Name
	}
	if pci.Request != nil {
		fmt.Fprintln(v, "[", colorCyan, pci.Request.Method, colorReset, "]", pci.Name)
		// fmt.Fprintf(v, "[%s] %s\n", pci.Request.Method, pci.Name)
		// fmt.Fprint(v, string(colorReset))
	} else if len(pci.Children) > 0 {
		chev := "> "
		if pci.Expanded() {
			chev = "\\/"
		}
		fmt.Fprintf(v, "%s %s\n", chev, n)
		if pci.Expanded() {
			for _, child := range pci.Children {
				printCollection(v, pad+" ", child) //it.rootSelection == pciIdx && it.selection == childIdx)
			}
		}
	}
	/*
		else {
			fmt.Fprint(v, "\t\t")
		}
		fmt.Fprintf(v, "%s %t", pci.Name, selected)
		fmt.Fprintln(v)
	*/
}
