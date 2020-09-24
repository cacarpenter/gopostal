package postman

import (
	"fmt"
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
	grp := gpmodel.RequestGroup{}
	pc := Collection{}
	pc.Name = "Test"
	wireCollection(&grp, &pc)
	if grp.Name != "Test" {
		t.Errorf("RequestGroup name should be Test")
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
	pc.Items = append(pc.Items, new(Collection))
	pc.Items[0].Items = append(pc.Items[0].Items, new(Collection))
	var grp gpmodel.RequestGroup
	wireCollection(&grp, &pc)
	if len(grp.Children) != 1 {
		t.Fatal("Expecting one child")
	}
	if len(grp.Children[0].Children) != 1 {
		t.Fatal("Should be one grandchild")
	}
}

func TestCollection_wireCollection_simple(t *testing.T) {
	pc := Collection{}
	pc.Info = &CollectionInfo{}
	pc.Info.Name = "simple"
	gp := wireCollection(nil, &pc)
	if gp.Name != "simple" {
		t.Fatal("Collect name isnt simple")
	}
	if len(gp.Children) != 0 {
		t.Fatal("Len children should be zero not", len(gp.Children))
	}
	if len(gp.Requests) != 0 {
		t.Fatal("Len request should be zero not", len(gp.Requests))
	}
}

func TestCollection_wireCollection_TwoReqs(t *testing.T) {
	pc := Collection{}
	pc.Info = new(CollectionInfo)
	pc.Info.Name = "coll name"
	pc.Items = make([]*Collection, 2)
	pc.Items[0] = new(Collection)
	r1 := new(Request)
	r1.Method = "GET"
	pc.Items[0].Name = "Get req 1"
	pc.Items[0].Request = r1

	r2 := new(Request)
	r2.Method = "GET"
	pc.Items[1] = new(Collection)
	pc.Items[1].Name = "Get req 2"
	pc.Items[1].Request = r2

	gp := wireCollection(nil, &pc)
	if gp.Name != "coll name" {
		t.Fatal("RequestGroup name should be 'coll name' not", gp.Name)
	}
	if len(gp.Children) != 0 {
		t.Fatal("Len children should be zero not", len(gp.Children))
	}
	if len(gp.Requests) != 2 {
		t.Fatal("Len request should be zero not", len(gp.Requests))
	}
	if gp.Requests[0].Method != "GET" {
		t.Fatal("Req 0 method should be GET")
	}
	if gp.Requests[0].Name != "Get req 1" {
		t.Fatal("Req 0 not named properly", gp.Requests[0].Name)
	}
	if gp.Requests[1].Name != "Get req 2" {
		t.Fatal("Req 1 not named properly", gp.Requests[1].Name)
	}
}

func TestParseCollection_Example(t *testing.T) {
	c, err := ParseCollection("../test_assets/Example.postman_collection.json")
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("Collection should not be null")
	}
	if c.Name != "Example" {
		t.Fatal("group name should be Example and not", c.Name)
	}
	if len(c.Requests) > 0 {
		t.Fatal("Request unexpected here")
	}
	if len(c.Children) != 2 {
		t.Fatal("Should have 2 children but has ", len(c.Children))
	}
	for i, ch := range c.Children {
		fmt.Println(i, ch.Name)
	}
	if c.Children[0].Name != "Folder1" {
		t.Fatal("Should have Folder1 as zeroth child", c.Children[0].Name)
	}
	if c.Children[0].Parent != c {
		t.Fatal("Child 0 does not have expected parent", c.Parent)
	}
}

func TestParseCollection_Example2(t *testing.T) {
	c, err := ParseCollection("../test_assets/Example2.postman_collection.json")
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("Collection should not be null")
	}
	if c.Name != "Example2" {
		t.Fatal("group name should be Example and not", c.Name)
	}
	if len(c.Requests) > 0 {
		t.Fatal("Request unexpected here")
	}
	if len(c.Children) != 2 {
		t.Fatal("Should have 2 children but has ", len(c.Children))
	}
	f1 := c.Children[0]
	if f1.Name != "Folder1" {
		t.Fatal("Should have Folder1 as zeroth child")
	}
	if f1.Parent != c {
		t.Fatal("Child 0 does not have expected parent", c.Parent)
	}
	if len(f1.Requests) != 2 {
		t.Fatal("Folder1 should have 2 requests but found", len(f1.Requests))
	}
	f2 := c.Children[1]
	if f2.Name != "Folder2" {
		t.Fatal("child1 not named Folder2")
	}
	if len(f2.Children) != 1 {
		t.Fatal("Folder2 should have one child")
	}
	if f2.Children[0].Name != "Folder 2 : 1" {
		t.Fatal("Folder 2:1 not found", f2.Children[0].Name)
	}
}
