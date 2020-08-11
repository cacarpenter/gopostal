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
)

type TreeWidget struct {
	*log.Logger
	tree         *treeNode
	selectedNode *treeNode
	selectedRow  int
	maxRows      int // changes based on expanded children
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
			fmt.Fprintf(w, "%s%s%s\n", colorGreen, n.label, colorReset)
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

func (n *treeNode) toggleExpanded() int {
	n.expanded = !n.expanded
	// FIX ME this is not right as it depends on whether the children are expanded, and their children etc. Need another recursive function
	nodeChange := len(n.children)
	// node isn't expanded so the change is negative
	if !n.expanded {
		nodeChange *= -1
	}
	return nodeChange
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

func (n *treeNode) prevSibling() *treeNode {
	// this is the root
	if n.parent == nil {
		return nil
	}
	var prev *treeNode
	for _, ch := range n.parent.children {
		if n == ch {
			break
		}
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
	cx, cy := v.Cursor()
	sx, sy := v.Size()
	tw.Logger.Printf("Layout| Cursor %d,%d | Size %d,%d| Origin %d,%d | selected %d\n", cx, cy, sx, sy, ox, oy, tw.selectedRow)

	// handle scrolling
	if tw.selectedRow+1-oy > sy {
		v.SetOrigin(ox, oy+1)
	} else if tw.selectedRow-oy < 0 {
		v.SetOrigin(ox, oy-1)
	}
	numPrinted := 0
	for _, grp := range tw.tree.children {
		maxItemNameLength := 0
		for _, n := range grp.children {
			if len(n.label) > maxItemNameLength {
				maxItemNameLength = len(n.label)
			}
		}
		maxItemNameLength = maxItemNameLength/2 + 1
		numPrinted += printNode(v, grp, " ")
	}
	tw.Logger.Printf("Layout printed %d rows\n", numPrinted)
}

func printNode(w io.Writer, node *treeNode, pad string) int {
	numPrinted := 0
	node.print(w, pad)
	numPrinted++

	if len(node.children) > 0 {
		chev := ArrowRightOpen
		if node.expanded {
			chev = ArrowDownOpen
		}
		label := fmt.Sprintf("%c %s", chev, node.label)
		if node.selected {
			fmt.Fprintf(w, "%s%s%s\n", colorGreen, label, colorReset)
		} else {
			fmt.Fprintln(w, label)
		}
		if node.expanded {
			for _, child := range node.children {
				numPrinted += printNode(w, child, pad+pad)
			}
		}
	}

	return numPrinted
}

func (tw *TreeWidget) MoveUp() {
	tw.Logger.Println("moveUp")
	if tw.selectedNode == nil {
		tw.Logger.Println("MoveUp: Nothing selected")
		return
	}
	var nextItem *treeNode
	prevSib := tw.selectedNode.prevSibling()
	if prevSib != nil {
		tw.Logger.Println("MoveUp: No previous sibling")
		if prevSib.expanded {
			if len(prevSib.children) > 0 {
				nextItem = prevSib.children[len(prevSib.children)-1]
			}
		} else {
			nextItem = prevSib
		}
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
		tw.Logger.Printf("MoveUp: selected for now %d\n", tw.selectedRow)
	}
}

func (tw *TreeWidget) MoveDown() {
	l := tw.Logger
	var nextNode *treeNode
	if tw.selectedNode == nil {
		tw.Logger.Println("no current selection")
		return
	}
	l.Println("MoveDown: Current Selection is ", tw.selectedNode.label)
	if tw.selectedNode.expanded && len(tw.selectedNode.children) > 0 {
		l.Println("expanding with children, so selecting the first child")
		nextNode = tw.selectedNode.children[0]
	} else {
		nextSib := tw.selectedNode.nextSibling()
		if nextSib != nil {
			nextNode = nextSib
		} else if tw.selectedNode.parent != nil {
			parentSib := tw.selectedNode.parent.nextSibling()
			if parentSib != nil {
				nextNode = parentSib
			}
		}
	}
	if nextNode != nil {
		l.Println("Setting next item to ", nextNode.label)
		tw.selectedNode.selected = false
		nextNode.selected = true
		tw.selectedNode = nextNode
		tw.selectedRow++
		tw.Logger.Printf("MoveDown: selected for now %d\n", tw.selectedRow)

		// we may need to scroll
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
	tw.Logger.Println("ToggleExpanded")
	if tw.selectedNode != nil {
		tw.maxRows = tw.maxRows + tw.selectedNode.toggleExpanded()
		tw.Logger.Printf("ToggleExpanded: MaxRows is now %d\n", tw.maxRows)
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

func (tw *TreeWidget) SetGroups(gps []*gpmodel.Group) {
	tw.Logger.Println("TreeWidget.SetGroups: load", len(gps))
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
	tw.maxRows = numNodes
	tw.Logger.Printf("TreeWidget has %d current max nodes\n", tw.maxRows)
}

func group2node(group *gpmodel.Group) *treeNode {
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
