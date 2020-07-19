package gpmodel

import "testing"

func TestGroup_AddChild(t *testing.T) {
	parent := new(Group)
	child := new(Group)
	parent.AddChild(child)

	if len(parent.Children) != 1 {
		t.Errorf("Child length is not zero but %d", len(parent.Children))
	}
	if &parent.Children[0] != child {
		t.Error("parent does not have child set")
	}
	if child.Parent() != parent {
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