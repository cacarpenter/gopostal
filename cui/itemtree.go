package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
)

type ItemTree struct {
	collection *postman.Collection

	// state
	selection int

	// changing size of the tree depending on what is expanded and collapsed
	treeSize int
}

func NewItemTree(pc *postman.Collection) *ItemTree {
	it := ItemTree{pc, 0, len(pc.Children)}
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
	fmt.Fprintf(v, " ~ %d of (%d) ~\n", it.selection, it.treeSize)
	for _, pci := range coll.Children {
		printCollectionItem(v, "", &pci, false)
	}
}

func (it *ItemTree) ArrowUp() {
	it.selection--
	/*
		if it.selection < 1 {
			if it.rootSelection > 0 {
				it.rootSelection--
				it.selection = len(it.collection.Items[it.rootSelection].Children) - 1
			}
		} else {
			it.selection--
		}
	*/
}

func (it *ItemTree) ArrowDown() {
	it.selection++
	/*
		numSel := len(it.collection.Items[it.rootSelection].Children)
		if it.selection < numSel-1 {
			it.selection++
		} else {
			it.selection = 0
			it.rootSelection++
		}
	*/
}

func (it *ItemTree) ToggleExpend() {
	ci := it.currentCollectionItem()
	numChild := len(ci.Children)
	if ci.ToggleExpanded() {
		it.treeSize += numChild
	} else {
		it.treeSize -= numChild
	}
}

func (it *ItemTree) currentCollectionItem() *postman.Collection {
	//var selected *postman.CollectionItem
	selectedCi := &it.collection.Children[0]
	for i := 0; i < it.selection; i++ {
		//		selected
		//	if rootPci.Expanded() {}
	}
	return selectedCi
}

func printCollectionItem(v *gocui.View, pad string, pci *postman.Collection, selected bool) {
	if selected {
		fmt.Fprintf(v, "%s* ", pad)
	} else {
		fmt.Fprintf(v, "%s  ", pad)
	}
	if pci.Request != nil {
		fmt.Fprintf(v, "[%s] %s\n", pci.Request.Method, pci.Name)
	} else if len(pci.Children) > 0 {
		chev := "> "
		if pci.Expanded() {
			chev = "\\/"
		}
		fmt.Fprintf(v, "%s %s\n", chev, pci.Name)
		if pci.Expanded() {
			for _, child := range pci.Children {
				printCollectionItem(v, pad+" ", &child, false) //it.rootSelection == pciIdx && it.selection == childIdx)
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
