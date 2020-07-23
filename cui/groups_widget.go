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
		fmt.Fprintf(v, "%c ", SelectIcon)
	} else {
		fmt.Fprint(v, "  ")
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
	/*
		var nextItem *gpmodel.Group
		prevSib := gw.selectedNode.PreviousSibling()
		if prevSib != nil {
			if prevSib.expanded {
				nextItem = prevSib.Children[len(prevSib.Children)-1]
			} else {
				nextItem = prevSib
			}
		} else if gw.selectedGroup.Parent() != nil {
			if gw.selectedGroup.Parent().Expanded() {
				nextItem = gw.selectedGroup.Parent()
			} else {
				parentSib := gw.selectedGroup.Parent().PreviousSibling()
				if parentSib != nil {
					nextItem = parentSib
				}
			}
		}
		if nextItem == nil {
			gw.Logger.Println("Checking for a previous collection")
			// check for another collection
			if gw.currentGroupIdx > 0 {
				gw.currentGroupIdx--
				// if the previous collection is expanded and has children
				curr := gw.groups[gw.currentGroupIdx]
				if curr.Expanded() && len(curr.Children) > 0 {
					// select the last child
					nextItem = curr.Children[len(curr.Children)-1]
				} else {
					// otherwise gw is the previous collection gwself
					nextItem = gw.groups[gw.currentGroupIdx]
				}
			} else {
				gw.Logger.Println("Move Up already at first collection")
			}
		}
		if nextItem != nil {
			gw.selectedGroup.SetSelected(false)
			nextItem.SetSelected(true)
			gw.selectedGroup = nextItem
		} else {
			gw.Logger.Println("MoveUp: No nextItem")
		}
	*/
}

func (gw *GroupsWidget) MoveDown() {
	l := gw.Logger
	l.Println("moveDown")
	var nextNode *gwNode
	if gw.selectedNode == nil {
		gw.Logger.Println("no current selection")
		return
	}
	l.Println("MoveDown: Current Selection is ", gw.selectedNode.label)
	if gw.selectedNode.expanded {
		if len(gw.selectedNode.children) > 0 {
			nextNode = gw.selectedNode.children[0]
		}
	} else {
		l.Println("MoveDown: Current not expanded, look for the next sibling")
		nextSib := gw.selectedNode.nextSibling()
		if nextSib != nil {
			nextNode = nextSib
		} else if gw.selectedNode.parent != nil {
			l.Println("Selected Parent is", gw.selectedNode.parent.label)
			parentSib := gw.selectedNode.parent.nextSibling()
			if parentSib != nil {
				l.Println("Found parent next sib")
				nextNode = parentSib
			} else {
				l.Println("No parent next sib")
			}
		} else {
			gw.Logger.Println("No parent and no next sib")
		}
	}
	if nextNode != nil {
		gw.Logger.Println("Setting next item to ", nextNode.label)
		gw.selectedNode.selected = false
		nextNode.selected = true
		gw.selectedNode = nextNode
	} else {
		gw.Logger.Println("moveDown: No nextItem")
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
		gw.currentGroupIdx = len(gw.groups) - 1
		if gw.currentGroupIdx > -1 {
			rootColl := gw.groups[gw.currentGroupIdx]
			gw.selectedGroup = rootColl.LastExpandedDescendent()
		}

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
		t.children[i].expanded = true
	}
	gw.tree = t
	gw.selectedNode = t.children[0]
	gw.selectedNode.selected = true
}

func group2node(group *gpmodel.Group) *gwNode {
	n := gwNode{}
	n.label = group.Name
	if group.Request != nil {
		n.request = group.Request
	} else {
		n.children = make([]*gwNode, len(group.Children))
		for i, grpChild := range group.Children {
			n.children[i] = group2node(grpChild)
		}
	}
	return &n
}

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
