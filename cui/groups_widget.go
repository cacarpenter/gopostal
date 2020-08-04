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

type gwNode struct {
	label              string
	parent             *gwNode
	children           []*gwNode
	request            *gpmodel.RequestSpec
	expanded, selected bool
}

func (g *gwNode) expand(exp, recursive bool) {
	g.expanded = exp
	if recursive {
		for _, ch := range g.children {
			ch.expand(exp, recursive)
		}
	}
}

func (g *gwNode) toggleExpanded() bool {
	g.expanded = !g.expanded
	return g.expanded
}

func (g *gwNode) nextSibling() *gwNode {
	// root
	if g.parent == nil {
		return nil
	}
	numSibs := len(g.parent.children)
	for i, ch := range g.parent.children {
		if g == ch && i < numSibs-1 {
			return g.parent.children[i+1]
		}
	}
	return nil
}

func (g *gwNode) prevSibling() *gwNode {
	// this is the root
	if g.parent == nil {
		return nil
	}
	var prev *gwNode
	for _, ch := range g.parent.children {
		if g == ch {
			break
		}
		prev = ch
	}
	return prev
}

type GroupsWidget struct {
	*log.Logger
	tree         *gwNode
	selectedNode *gwNode
}

func (gw *GroupsWidget) Layout(v *gocui.View) {
	v.Clear()
	if gw.tree == nil {
		gw.Logger.Println("No tree set, nothing to layout")
		return
	}
	for _, grp := range gw.tree.children {
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

func printNode(v *gocui.View, pad string, node *gwNode) {
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

func (gw *GroupsWidget) MoveUp() {
	gw.Logger.Println("moveUp")
	if gw.selectedNode == nil {
		gw.Logger.Println("MoveUp: Nothing selected")
		return
	}
	var nextItem *gwNode
	prevSib := gw.selectedNode.prevSibling()
	if prevSib != nil {
		gw.Logger.Println("MoveUp: No previous sibling")
		if prevSib.expanded {
			if len(prevSib.children) > 0 {
				nextItem = prevSib.children[len(prevSib.children)-1]
			}
		} else {
			nextItem = prevSib
		}
	} else if gw.selectedNode.parent != nil {
		if gw.selectedNode.parent.expanded {
			nextItem = gw.selectedNode.parent
		} else {
			parentSib := gw.selectedNode.parent.prevSibling()
			if parentSib != nil {
				nextItem = parentSib
			}
		}
	}
	if nextItem != nil {
		gw.selectedNode.selected = false
		nextItem.selected = true
		gw.selectedNode = nextItem
	} else {
		gw.Logger.Println("MoveUp: No nextItem")
	}
}

func (gw *GroupsWidget) MoveDown() {
	l := gw.Logger
	var nextNode *gwNode
	if gw.selectedNode == nil {
		gw.Logger.Println("no current selection")
		return
	}
	l.Println("MoveDown: Current Selection is ", gw.selectedNode.label)
	if gw.selectedNode.expanded && len(gw.selectedNode.children) > 0 {
		nextNode = gw.selectedNode.children[0]
	} else {
		nextSib := gw.selectedNode.nextSibling()
		if nextSib != nil {
			nextNode = nextSib
		} else if gw.selectedNode.parent != nil {
			parentSib := gw.selectedNode.parent.nextSibling()
			if parentSib != nil {
				nextNode = parentSib
			}
		}
	}
	if nextNode != nil {
		l.Println("Setting next item to ", nextNode.label)
		gw.selectedNode.selected = false
		nextNode.selected = true
		gw.selectedNode = nextNode
	}
}

func (gw *GroupsWidget) CollapseAll() {
	if gw.tree == nil {
		return
	}
	for _, ch := range gw.tree.children {
		ch.expand(false, true)
	}
}

func (gw *GroupsWidget) ExpandAll() {
	if gw.tree == nil {
		return
	}
	for _, ch := range gw.tree.children {
		ch.expand(true, true)
	}
}

func (gw *GroupsWidget) ToggleExpanded() {
	gw.Logger.Println("ToggleExpanded")
	if gw.selectedNode != nil {
		gw.selectedNode.toggleExpanded()
	}
}

func (gw *GroupsWidget) SelectLast() {
	/*
	if len(gw.tree.children) > 0 {}
	lastChildIdx := len(gw.tree.children) - 1
	rootColl := gw.groups[gw.currentGroupIdx]
	gw.selectedGroup = rootColl.LastExpandedDescendent()
	 */
}

func (gw *GroupsWidget) SetGroups(gps []*gpmodel.Group) {
	gw.Logger.Println("GroupsWidget.SetGroups: load", len(gps))
	if len(gps) < 1 {
		return
	}
	t := new(gwNode)
	t.children = make([]*gwNode, len(gps))
	for i, g := range gps {
		t.children[i] = group2node(g)
		t.children[i].parent = t
		t.children[i].expanded = true
	}
	gw.tree = t
	gw.tree.label = "root" // won't be shown in UI but useful for testing
	gw.selectedNode = t.children[0]
	gw.selectedNode.selected = true
}

func group2node(group *gpmodel.Group) *gwNode {
	n := gwNode{}
	n.label = group.Name
	n.children = make([]*gwNode, len(group.Children)+len(group.Requests))
	chIdx := 0
	for _, req := range group.Requests {
		// create a request node for each request
		reqNode := new(gwNode)
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
