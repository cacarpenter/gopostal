package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/jroimartin/gocui"
)

type ItemTree struct {
	collection *postman.Collection

	// state
	rootSelection, selection int
	expendedRoots            []bool
}

func NewItemTree(pc *postman.Collection) *ItemTree {
	it := ItemTree{pc, 0, 0, make([]bool, len(pc.Items))}
	return &it
}

func (it *ItemTree) Layout(v *gocui.View) {
	coll := it.collection
	maxItemNameLength := 0
	for _, n := range state.openCollection.Items {
		if len(n.Name) > maxItemNameLength {
			maxItemNameLength = len(n.Name)
		}
	}
	maxItemNameLength = maxItemNameLength/2 + 1

	v.Clear()
	fmt.Fprintf(v, " ~ %d %d ~\n", it.rootSelection, it.selection)
	for pciIdx, pci := range coll.Items {
		fmt.Fprintf(v, " > %s\n", pci.Name)
		if it.expendedRoots[pciIdx] {
			for childIdx, child := range pci.Children {
				printCollectionItem(v, &child, it.rootSelection == pciIdx && it.selection == childIdx)
			}
		}
	}
}

func (it *ItemTree) ArrowUp() {
	it.selection = it.selection - 1
}

func (it *ItemTree) ArrowDown() {
	numSel := len(it.collection.Items[it.rootSelection].Children)
	if it.selection < numSel+1 {
		it.selection++
	} else {
		it.selection = 0
		it.rootSelection++
	}
}

func (it *ItemTree) ToggleExpend() {
	// just roots right now
	it.expendedRoots[it.rootSelection] = !it.expendedRoots[it.rootSelection]
}

func printCollectionItem(v *gocui.View, pci *postman.CollectionItem, selected bool) {
	if pci.Request != nil {
		fmt.Fprintf(v, "[%s]\t", pci.Request.Method)
	} else if len(pci.Children) > 0 {
		fmt.Fprintf(v, " > \t")
	} else {
		fmt.Fprint(v, "\t\t")
	}
	fmt.Fprintf(v, "\t%q %t", pci.Name, selected)
	fmt.Fprintln(v)
}
