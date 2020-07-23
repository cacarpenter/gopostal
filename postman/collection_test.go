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
	wireCollection(&grp,&pc)
	if grp.Name != "Test" {
		t.Errorf("Group name should be Test")
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
	if c.Request != nil {
		t.Fatal("Request unexpected here")
	}
	if len(c.Children) != 2 {
		t.Fatal("Should have 2 children but has ", len(c.Children))
	}
	if c.Children[0].Name != "Folder1" {
		t.Fatal("Should have Folder1 as zeroth child")
	}
	if c.Children[0].Parent() != c {
		t.Fatal("Child 0 does not have expected parent", c.Parent())
	}
}

