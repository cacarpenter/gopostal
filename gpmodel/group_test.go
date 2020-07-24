package gpmodel

import "testing"

func TestGroup_AddChild(t *testing.T) {
	parent := new(Group)
	child := new(Group)
	parent.AddChild(child)

	if len(parent.Children) != 1 {
		t.Errorf("Child length is not zero but %d", len(parent.Children))
	}
	if parent.Children[0] != child {
		t.Error("parent does not have child set")
	}
	if child.Parent != parent {
		t.Error("child does not have parent set")
	}
}

func TestGroup_NextSibling(t *testing.T) {
	parent := new(Group)
	child1 := new(Group)
	child2 := new(Group)
	parent.AddChild(child1)
	parent.AddChild(child2)

	if child1.NextSibling() != child2 {
		t.Error("child1 next sibling should be child2")
	}
	if child2.NextSibling() != nil {
		t.Error("child2 next sibling is not nil")
	}
}

func TestGroup_PreviousSibling(t *testing.T) {
	parent := new(Group)
	child1 := new(Group)
	child2 := new(Group)
	parent.AddChild(child1)
	parent.AddChild(child2)

	if child2.PreviousSibling() != child1 {
		t.Error("child2 previous sibling should be child1")
	}
	if child1.PreviousSibling() != nil {
		t.Error("child1 previous sibling is not nil")
	}
}

/*
func TestGroup_LastExpandedDescendent(t *testing.T) {
	p := new(Group)
	p.Name = "parent"
	c1 := new(Group)
	c1.Name = "child1"
	c2 := new(Group)
	c2.Name = "child2"
	c3 := new(Group)
	c3.Name = "child3"
	p.AddChild(c1)
	p.AddChild(c2)
	p.AddChild(c3)
	c1g1 := new(Group)
	c1g1.Name = "child1 grand1"
	c1.AddChild(c1g1)
	c1g2 := new(Group)
	c1g2.Name = "child1 grand2"
	c1.AddChild(c1g2)
	c2g1 := new(Group)
	c2.AddChild(c2g1)
	gg1 := new(Group)
	c1g1.AddChild(gg1)
	gg2 := new(Group)
	c1g2.AddChild(gg2)
	gg3 := new(Group)
	c1g1.AddChild(gg3)

	if p.LastExpandedDescendent() != p {
		t.Fatal("unexpanded p should return itself")
	}
	p.ToggleExpanded()
	if p.LastExpandedDescendent() != c3 {
		t.Fatalf("p LastExpandedDescended != c3 but is %s", p.LastExpandedDescendent().Name)
	}
	c2.ToggleExpanded()
	if p.LastExpandedDescendent() != c3 {
		t.Fatalf("p LastExpandedDescended should still be c3 but is %s", p.LastExpandedDescendent().Name)
	}
	c1.ToggleExpanded()
	c1g2.ToggleExpanded()
	if c1.LastExpandedDescendent() != gg2 {
		t.Fatalf("c1 LastExpanded should be gg2 but was %s", c1.LastExpandedDescendent().Name)
	}
}
 */
