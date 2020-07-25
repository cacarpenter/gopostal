package cui

import (
	"github.com/cacarpenter/gopostal/gpmodel"
	"log"
	"os"
	"testing"
)

func setUp() *GroupsWidget {
	gw := &GroupsWidget{}
	gw.Logger = log.New(os.Stdout, "", 0)
	gw.tree = new(gwNode)
	return gw
}

func TestGroupsWidget_MoveDown_Empty(t *testing.T) {
	gw := setUp()
	gw.MoveDown()
}

func TestGroupsWidget_MoveDown_SingleCollection(t *testing.T) {
	gw := setUp()
	gw.tree.children = make([]*gwNode, 1)
	gw.tree.children[0] = new(gwNode)
	gw.selectedNode = gw.tree.children[0]
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[0] {
		t.Fatal("Selected collection should not be different")
	}
}

func TestGroupsWidget_MoveDown_ThreeCollections(t *testing.T) {
	gw := setUp()
	grps := make([]*gpmodel.Group, 3)
	grps[0] = gpmodel.NewGroup("group 0")
	grps[1] = gpmodel.NewGroup("group 1")
	grps[2] = gpmodel.NewGroup("group 2")
	gw.SetGroups(grps)
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[1] {
		t.Fatal("Selected collection should be 1 not ", gw.selectedNode.label)
	}
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[2] {
		t.Fatal("Selected collection should 2")
	}
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[2] {
		t.Fatal("Selected collection should still be 2")
	}
}

func TestGroupsWidget_MoveDown_SingleManyDecendents(t *testing.T) {
	gw := setUp()
	groups := make([]*gpmodel.Group, 1)
	groups[0] = gpmodel.NewGroup("Only Child")
	groups[0].AddChild(new(gpmodel.Group))
	groups[0].Children[0].AddChild(new(gpmodel.Group))
	gw.SetGroups(groups)

	if gw.selectedNode != gw.tree.children[0] {
		t.Fatal("root not selected")
	}
	gw.ToggleExpanded() // collapse the root. It is expanded by default
	gw.MoveDown()
	// root should still be selected because gw is not expanded
	if gw.selectedNode != gw.tree.children[0] {
		t.Fatal("root still not selected")
	}
	gw.ToggleExpanded()
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[0].children[0] {
		t.Fatal("first descendent not selected")
	}
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[0].children[0] {
		t.Fatal("first descendent still not selected")
	}
	gw.ToggleExpanded() // expand the first decendent
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[0].children[0].children[0] {
		t.Fatal("second descendent not selected")
	}
}

func TestGroupsWidget_MoveDown_ThreeCollsOneDescendentEach(t *testing.T) {
	gw := setUp()
	grps := make([]*gpmodel.Group, 3)
	grps[0] = gpmodel.NewGroup("group 0")
	grps[0].AddChild(gpmodel.NewGroup("group 0 child 0"))
	grps[1] = gpmodel.NewGroup("group 1")
	grps[1].AddChild(new(gpmodel.Group))
	grps[2] = gpmodel.NewGroup("group 2")
	grps[2].AddChild(new(gpmodel.Group))
	gw.SetGroups(grps)

	gw.ToggleExpanded()
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[0] {
		t.Fatal("Not p 0 c 0!", gw.selectedNode.label)
	}
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[1] {
		t.Fatal("Not p 1")
	}
	gw.ToggleExpanded()
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[1].children[0] {
		t.Fatal("Not p 1 c 0")
	}
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[2] {
		t.Fatal("Not p 2")
	}
	gw.ToggleExpanded()
	gw.MoveDown()
	if gw.selectedNode != gw.tree.children[2].children[0] {
		t.Fatal("Not p 2 c 0")
	}
}

/*
func TestGroupsWidget_MoveUp_ThreeCollsOneDescendentEach(t *testing.T) {
	gw := setUp()
	grps := make([]*gpmodel.Group, 3)
	grps[0] = new(gpmodel.Group)
	grps[0].Name = "p0"
	grps[0].AddChild(new(gpmodel.Group))
	grps[1] = new(gpmodel.Group)
	grps[1].Name = "p1"
	grps[1].AddChild(new(gpmodel.Group))
	grps[2] = new(gpmodel.Group)
	grps[2].Name = "p2"
	grps[2].AddChild(new(gpmodel.Group))
	gw.tree.children = grps
	gw.ExpandAll()
	gw.SelectLast()
	if gw.selectedNode != grps[2].Children[0] {
		t.Fatalf("Not at last p2 c0 : %s\n", gw.selectedNode.Name)
	}

	gw.MoveUp()
	if gw.selectedNode != grps[2] {
		t.Fatalf("Not p 2 : %s\n", gw.selectedNode.Name)
	}
	gw.MoveUp()
	if gw.selectedNode != grps[1].Children[0] {
		t.Fatalf("Not p 1 c 0 : %s\n", gw.selectedNode.Name)
	}
	gw.MoveUp()
	if gw.selectedNode != grps[1] {
		t.Fatal("Not p 1")
	}
	gw.MoveUp()
	if gw.selectedNode != grps[0].Children[0] {
		t.Fatal("Not p 0 c 0")
	}
	gw.MoveUp()
	if gw.selectedNode != grps[0] {
		t.Fatal("Not p 0")
	}
	gw.MoveUp()
	if gw.selectedNode != grps[0] {
		t.Fatal("Not still p 0")
	}
}

*/
