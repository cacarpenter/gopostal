package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
	"github.com/jroimartin/gocui"
	"log"
)

const (
	SelectIcon     = '\u066D'
	ArrowDownSolid = '\u25BC'
	ArrowDownOpen  = '\u1401'
	ArrowRightOpen = '\u1405'
)

type treeNode struct {
	label              string
	parent             *treeNode
	children           []*treeNode
	request            *gpmodel.RequestSpec
	expanded, selected bool
}

func (n *treeNode) expand(exp, recursive bool) {
	n.expanded = exp
	if recursive {
		for _, ch := range n.children {
			ch.expand(exp, recursive)
		}
	}
}

func (n *treeNode) toggleExpanded() bool {
	n.expanded = !n.expanded
	return n.expanded
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

type TreeWidget struct {
	*log.Logger
	tree         *treeNode
	selectedNode *treeNode
}

func (tw *TreeWidget) Layout(v *gocui.View) {
	v.Clear()
	if tw.tree == nil {
		tw.Logger.Println("No tree set, nothing to layout")
		return
	}
	for _, grp := range tw.tree.children {
		maxItemNameLength := 0
		for _, n := range grp.children {
			if len(n.label) > maxItemNameLength {
				maxItemNameLength = len(n.label)
			}
		}
		maxItemNameLength = maxItemNameLength/2 + 1
		// cx, cy := v.Cursor()
		// fmt.Fprintln(v, cx, cy)
		printNode(v, "", grp)
	}
}

func printNode(v *gocui.View, pad string, node *treeNode) {
	if node.selected {
		fmt.Fprintf(v, "%c", SelectIcon)
	} else {
		fmt.Fprint(v, " ")
	}
	fmt.Fprint(v, pad)
	if node.request != nil {
		fmt.Fprint(v, "[", colorCyan, node.request.Method, colorReset, "] ")
		if node.selected {
			fmt.Fprintf(v, "%s%s%s\n", colorGreen, node.label, colorReset)
		} else {
			fmt.Fprintln(v, node.label)
		}
	} else if len(node.children) > 0 {
		chev := ArrowRightOpen
		if node.expanded {
			chev = ArrowDownOpen
		}
		label := fmt.Sprintf("%c %s", chev, node.label)
		if node.selected {
			fmt.Fprintf(v, "%s%s%s\n", colorGreen, label, colorReset)
		} else {
			fmt.Fprintln(v, label)
		}
		if node.expanded {
			for _, child := range node.children {
				printNode(v, pad+" ", child)
			}
		}
	} else {
		fmt.Fprintln(v, "?")
	}
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
	} else {
		tw.Logger.Println("MoveUp: No nextItem")
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
		tw.selectedNode.toggleExpanded()
	}
}

func (tw *TreeWidget) SelectLast() {
	/*
	if len(tw.tree.children) > 0 {}
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
	for i, g := range gps {
		t.children[i] = group2node(g)
		t.children[i].parent = t
		t.children[i].expanded = true
	}
	tw.tree = t
	tw.tree.label = "root" // won't be shown in UI but useful for testing
	tw.selectedNode = t.children[0]
	tw.selectedNode.selected = true
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
