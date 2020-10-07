package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
	"github.com/jroimartin/gocui"
	"io"
	"log"
)

const (
	SelectIcon     = '\u066D'
	ArrowDownSolid = '\u25BC'
	ArrowDownOpen  = '\u1401'
	ArrowRightOpen = '\u1405'
	selectColor = colorPurple
)

type TreeWidget struct {
	*log.Logger
	tree         *treeNode
	selectedNode *treeNode
	selectedRow  int
}

type treeNode struct {
	label              string
	parent             *treeNode
	children           []*treeNode
	request            *gpmodel.RequestSpec
	expanded, selected bool
}

func (n *treeNode) print(w io.Writer, pad string) {
	if n.selected {
		fmt.Fprintf(w, "%c", SelectIcon)
	} else {
		fmt.Fprint(w, " ")
	}
	fmt.Fprint(w, pad)
	if n.request != nil {
		fmt.Fprint(w, "[", methodColor(n.request.Method), n.request.Method, colorReset, "] ")
		if n.selected {
			fmt.Fprintf(w, "%s%s%s\n", selectColor, n.label, colorReset)
		} else {
			fmt.Fprintln(w, n.label)
		}
	}
}

func methodColor(m string) string {
	switch m {
	case "GET":
		return colorCyan
	case "DELETE":
		return colorRed
	case "PUT":
		return colorPurple
	case "POST":
		return colorGreen
	case "PATCH":
		return colorYellow
	}
	return colorReset
}

func (n *treeNode) expand(exp, recursive bool) int {
	n.expanded = exp
	nodesAdded := len(n.children)
	if recursive {
		for _, ch := range n.children {
			nodesAdded += ch.expand(exp, recursive)
		}
	}
	return nodesAdded
}

func (n *treeNode) toggleExpanded() {
	n.expanded = !n.expanded
}

func (n *treeNode) lastChild() *treeNode {
	if len(n.children) > 0 {
		return n.children[len(n.children)-1]
	}
	return nil
}

func (n *treeNode) lastExpandedDescendant() *treeNode {
	if !n.expanded {
		return n
	}
	lc := n.lastChild()
	if lc != nil {
		return lc.lastExpandedDescendant()
	}
	return nil
}

func (n *treeNode) nextSibling() *treeNode {
	// root
	if n.parent == nil {
		return nil
	}
	numSibs := len(n.parent.children)
	for i, ch := range n.parent.children {
		if n == ch && i < numSibs-1 {
			return n.parent.children[i+1]
		}
	}
	return nil
}

func (n *treeNode) nextRelative() *treeNode {
	if n.parent == nil {
		return nil
	}
	var nextRel *treeNode
	nextSib := n.nextSibling()
	if nextSib != nil {
		nextRel = nextSib
	} else {
		parentNextRel := n.parent.nextRelative()
		if parentNextRel != nil {
			return parentNextRel
		}
	}
	return nextRel
}

// prevSibling returns nil if you are the first child
func (n *treeNode) prevSibling() *treeNode {
	// this is the root
	if n.parent == nil {
		return nil
	}
	var prev *treeNode
	for _, ch := range n.parent.children {
		// if this is the child
		if n == ch {
			break
		}
		// save this child for the next iteration
		prev = ch
	}
	return prev
}

func (tw *TreeWidget) Layout(v *gocui.View) {
	v.Clear()
	if tw.tree == nil {
		tw.Logger.Println("No tree set, nothing to layout")
		return
	}
	ox, oy := v.Origin()
	_, sy := v.Size()
	// tw.Logger.Printf("Layout| Size %d,%d| Origin %d,%d | selected %d\n", sx, sy, ox, oy, tw.selectedRow)

	// handle scrolling
	if tw.selectedRow < oy {
		// selected row is "above" tree, set the origin to the selected row
		v.SetOrigin(ox, tw.selectedRow)
	} else if tw.selectedRow+1-oy > sy { // scroll down
		// TODO really this should be a calculation like selectedRow - oy
		v.SetOrigin(ox, oy+1)
	}
	for _, grp := range tw.tree.children {
		maxItemNameLength := 0
		for _, n := range grp.children {
			if len(n.label) > maxItemNameLength {
				maxItemNameLength = len(n.label)
			}
		}
		maxItemNameLength = maxItemNameLength/2 + 1
		printNode(v, grp, " ")
	}
}

func printNode(w io.Writer, node *treeNode, pad string) {
	node.print(w, pad)

	if len(node.children) > 0 {
		chev := ArrowRightOpen
		if node.expanded {
			chev = ArrowDownOpen
		}
		label := fmt.Sprintf("%c %s", chev, node.label)
		if node.selected {
			fmt.Fprintf(w, "%s%s%s\n", selectColor, label, colorReset)
		} else {
			fmt.Fprintln(w, label)
		}
		if node.expanded {
			for _, child := range node.children {
				printNode(w, child, pad+pad)
			}
		}
	}

}

func (tw *TreeWidget) MoveUp() {
	if tw.selectedNode == nil {
		return
	}

	var nextItem *treeNode
	prevSib := tw.selectedNode.prevSibling()
	if prevSib != nil {
		nextItem = prevSib.lastExpandedDescendant()
	} else if tw.selectedNode.parent != nil {
		if tw.selectedNode.parent.expanded {
			nextItem = tw.selectedNode.parent
		} else {
			parentSib := tw.selectedNode.parent.prevSibling()
			if parentSib != nil {
				nextItem = parentSib
			}
		}
	}

	if nextItem != nil {
		tw.selectedNode.selected = false
		nextItem.selected = true
		tw.selectedNode = nextItem
		tw.selectedRow--
		// tw.Logger.Printf("MoveUp: selected for now %d\n", tw.selectedRow)
	}
}

func (tw *TreeWidget) MoveDown() {
	var nextNode *treeNode
	if tw.selectedNode == nil {
		tw.Logger.Println("Cannot move: no current selection")
		return
	}
	if tw.selectedNode.expanded && len(tw.selectedNode.children) > 0 {
		nextNode = tw.selectedNode.children[0]
	} else {
		nextNode = tw.selectedNode.nextRelative()
	}
	if nextNode != nil {
		// tw.Logger.Println("Setting next item to ", nextNode.label)
		tw.selectedNode.selected = false
		nextNode.selected = true
		tw.selectedNode = nextNode
		tw.selectedRow++
	}
}

func (tw *TreeWidget) CollapseAll() {
	if tw.tree == nil {
		return
	}
	for _, ch := range tw.tree.children {
		ch.expand(false, true)
	}
}

func (tw *TreeWidget) ExpandAll() {
	if tw.tree == nil {
		return
	}
	for _, ch := range tw.tree.children {
		ch.expand(true, true)
	}
}

func (tw *TreeWidget) ToggleExpanded() {
	if tw.selectedNode != nil {
		tw.selectedNode.toggleExpanded()
	}
}

func (tw *TreeWidget) SelectLast() {
	// nothing to do
	if len(tw.tree.children) < 1 {
		return
	}
	lastRootChild := tw.tree.children[len(tw.tree.children)-1]
	if tw.selectedNode != nil {
		tw.selectedNode.selected = false
	}
	tw.selectedNode = lastRootChild
	tw.selectedNode.selected = true
	/*
		lastChildIdx := len(tw.tree.children) - 1
			rootColl := tw.groups[tw.currentGroupIdx]
			tw.selectedGroup = rootColl.LastExpandedDescendent()
	*/
}

func (tw *TreeWidget) SetRequestGroups(gps []*gpmodel.RequestGroup) {
	tw.Logger.Println("TreeWidget.SetRequestGroups: load", len(gps))
	if len(gps) < 1 {
		return
	}
	t := new(treeNode)
	t.children = make([]*treeNode, len(gps))
	numNodes := len(gps)
	for i, g := range gps {
		t.children[i] = group2node(g)
		t.children[i].parent = t
		numNodes += t.children[i].expand(true, false)
	}
	tw.tree = t
	tw.tree.label = "root" // won't be shown in UI but useful for testing
	tw.selectedNode = t.children[0]
	tw.selectedNode.selected = true
	tw.selectedRow = 0
}

func group2node(group *gpmodel.RequestGroup) *treeNode {
	n := treeNode{}
	n.label = group.Name
	n.children = make([]*treeNode, len(group.Children)+len(group.Requests))
	chIdx := 0
	for _, req := range group.Requests {
		// create a request node for each request
		reqNode := new(treeNode)
		reqNode.label = req.Name
		reqNode.request = req
		reqNode.parent = &n
		n.children[chIdx] = reqNode
		chIdx++
	}
	for _, grpChild := range group.Children {
		n.children[chIdx] = group2node(grpChild)
		n.children[chIdx].parent = &n
		chIdx++
	}
	return &n
}
