package postman

import (
	"github.com/cacarpenter/gopostal/gpmodel"
	"testing"
)

func TestIsCollectionFile(t *testing.T) {
	if IsCollectionFile("garasabafw") {
		t.Error("garbage should not be a collection file")
	}
	good := "something.postman_collection.json"
	if !IsCollectionFile(good) {
		t.Errorf("%q should return true", good)
	}
}

func TestCollection_wireCollection(t *testing.T) {
	grp := gpmodel.Group{}
	pc := Collection{}
	pc.Name = "Test"
	wireCollection(&grp, &pc)
	if grp.Name != "Test" {
		t.Errorf("Group name should be Test")
	}
	if len(grp.Requests) != 0 {
		t.Error("There should be no requests")
	}
	if len(grp.Children) != 0 {
		t.Error("There should be no children")
	}
}

func TestCollection_wireCollection_SingleTree(t *testing.T) {
	pc := Collection{}
	pc.Children = append(pc.Children, new(Collection))
	pc.Children[0].Children = append(pc.Children[0].Children, new(Collection))
	var grp gpmodel.Group
	wireCollection(&grp, &pc)
	if len(grp.Children) != 1 {
		t.Fatal("Expecting one child")
	}
	if len(grp.Children[0].Children) != 1 {
		t.Fatal("Should be one grandchild")
	}
}

func TestParseCollection(t *testing.T) {
	c, err := ParseCollection("/home/ccarpenter/Documents/postman/Example.postman_collection.json")
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("Collection should not be null")
	}
	if c.Name != "Example" {
		t.Fatal("group name should be not ", c.Name)
	}
	if len(c.Requests) > 0 {
		t.Fatal("Request unexpected here")
	}
	if len(c.Children) != 2 {
		t.Fatal("Should have 2 children but has ", len(c.Children))
	}
	if c.Children[0].Name != "Folder1" {
		t.Fatal("Should have Folder1 as zeroth child")
	}
	if c.Children[0].Parent != c {
		t.Fatal("Child 0 does not have expected parent", c.Parent)
	}
}
