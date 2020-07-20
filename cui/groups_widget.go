package cui

import (
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
	"github.com/jroimartin/gocui"
	"log"
)

type GroupsWidget struct {
	*log.Logger
	groups          []*gpmodel.Group
	currentGroupIdx int
	selectedGroup   *gpmodel.Group
}

func (gw *GroupsWidget) Layout(v *gocui.View) {
	v.Clear()
	for _, grp := range gw.groups {
		maxItemNameLength := 0
		for _, n := range grp.Children {
			if len(n.Name) > maxItemNameLength {
				maxItemNameLength = len(n.Name)
			}
		}
		maxItemNameLength = maxItemNameLength/2 + 1
		// cx, cy := v.Cursor()
		// fmt.Fprintln(v, cx, cy)
		printGroup(v, "", grp)
	}
}

func printGroup(v *gocui.View, pad string, grp *gpmodel.Group) {
	fmt.Fprint(v, pad)
	if grp.Request != nil {
		//req := gp.RequestSpec(gwn)
		fmt.Fprint(v, "[", colorCyan, grp.Request.Method, colorReset, "] ")
		if grp.Selected() {
			fmt.Fprintln(v, colorGreen, grp.Name, colorReset)
		} else {
			fmt.Fprintln(v, grp.Name)
		}
	} else if len(grp.Children) > 0 {
		chev := "> "
		if grp.Expanded() {
			chev = "\\/"
		}
		label := fmt.Sprintf("%s %s", chev, grp.Name)
		if grp.Selected() {
			fmt.Fprintln(v, colorGreen, label, colorReset)
		} else {
			fmt.Fprintln(v, label)
		}
		if grp.Expanded() {
			for _, child := range grp.Children {
				printGroup(v, pad+" ", child)
			}
		}
	} else {
		fmt.Fprintln(v, "?")
	}
}

func (gw *GroupsWidget) MoveUp() {
	gw.Logger.Println("moveUp")
	if gw.selectedGroup == nil {
		gw.Logger.Println("MoveUp: Nothing selected")
	}
	var nextItem *gpmodel.Group
	prevSib := gw.selectedGroup.PreviousSibling()
	if prevSib != nil {
		if prevSib.Expanded() {
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
}

func (gw *GroupsWidget) MoveDown() {
	gw.Logger.Println("moveDown")
	var nextItem *gpmodel.Group
	if gw.selectedGroup == nil {
		gw.Logger.Println("no current selection")
		return
	}
	if gw.selectedGroup.Expanded() {
		if len(gw.selectedGroup.Children) > 0 {
			nextItem = gw.selectedGroup.Children[0]
		}
	} else {
		nextSib := gw.selectedGroup.NextSibling()
		if nextSib != nil {
			nextItem = nextSib
		} else if gw.selectedGroup.Parent() != nil {
			parentSib := gw.selectedGroup.Parent().NextSibling()
			if parentSib != nil {
				nextItem = parentSib
			}
		}
	}
	if nextItem == nil {
		gw.Logger.Println("Checking for another collection")
		// check for another collection
		if gw.currentGroupIdx+1 < len(gw.groups) {
			gw.currentGroupIdx++
			nextItem = gw.groups[gw.currentGroupIdx]
		}
	}
	if nextItem != nil {
		gw.Logger.Println("Setting next item to ", nextItem.Name)
		gw.selectedGroup.SetSelected(false)
		nextItem.SetSelected(true)
		gw.selectedGroup = nextItem
	} else {
		gw.Logger.Println("moveDown: No nextItem")
	}
}

func (gw *GroupsWidget) CollapseAll() {
	for _, pmColl := range gw.groups {
		pmColl.Expand(false, true)
	}
}

func (gw *GroupsWidget) ExpandAll() {
	for _, pmColl := range gw.groups {
		pmColl.Expand(true, true)
	}
}

func (gw *GroupsWidget) ToggleExpanded() {
	gw.Logger.Println("ToggleExpanded")
	if gw.selectedGroup != nil {
		gw.selectedGroup.ToggleExpanded()
	}
}

func (gw *GroupsWidget) SelectLast() {
	gw.currentGroupIdx = len(gw.groups) - 1
	if gw.currentGroupIdx > -1 {
		rootColl := gw.groups[gw.currentGroupIdx]
		gw.selectedGroup = rootColl.LastExpandedDescendent()
	}
}

func (gw *GroupsWidget) SetGroups(gps []*gpmodel.Group) {
	gw.groups = gps
	if len(gw.groups) > 0 {
		gw.selectedGroup = gw.groups[0]
	}
}