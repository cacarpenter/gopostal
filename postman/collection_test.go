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

