package cui

import (
	"github.com/cacarpenter/gopostal/postman"
	"log"
	"os"
	"testing"
)

func setUp() *ItemTree {
	it := &ItemTree{}
	it.Logger = log.New(os.Stdout, "", 0)
	return it
}

func TestItemTree_MoveDown_Empty(t *testing.T) {
	it := setUp()
	it.MoveDown()
}

func TestItemTree_MoveDown_SingleCollection(t *testing.T) {
	it := setUp()
	it.collections = make([]*postman.Collection, 1)
	it.collections[0] = new(postman.Collection)
	it.selected = it.collections[0]
	it.MoveDown()
	if it.selected != it.collections[0] {
		t.Fatal("Selected collection should not be different")
	}
}

func TestItemTree_MoveDown_ThreeCollections(t *testing.T) {
	it := setUp()
	colls := make([]*postman.Collection, 3)
	colls[0] = new(postman.Collection)
	colls[1] = new(postman.Collection)
	colls[2] = new(postman.Collection)
	it.SetCollections(colls)
	it.MoveDown()
	if it.selected != it.collections[1] {
		t.Fatal("Selected collection should be 1")
	}
	it.MoveDown()
	if it.selected != it.collections[2] {
		t.Fatal("Selected collection should 2")
	}
	it.MoveDown()
	if it.selected != it.collections[2] {
		t.Fatal("Selected collection should still be 2")
	}
}


func TestItemTree_MoveDown_SingleManyDecendents(t *testing.T) {
	it := setUp()
	colls := make([]*postman.Collection, 1)
	colls[0] = new(postman.Collection)
	colls[0].AddChild(new(postman.Collection))
	colls[0].Children[0].AddChild(new(postman.Collection))
	it.SetCollections(colls)

	if it.selected != colls[0] {
		t.Fatal("root not selected")
	}
	it.MoveDown()
	// root should still be selected because it is not expanded
	if it.selected != colls[0] {
		t.Fatal("root still not selected")
	}
	it.ToggleExpanded() // expand the root
	it.MoveDown()
	if it.selected != colls[0].Children[0] {
		t.Fatal("first descendent not selected")
	}
	it.MoveDown()
	if it.selected != colls[0].Children[0] {
		t.Fatal("first descendent still not selected")
	}
	it.ToggleExpanded() // expand the first decendent
	it.MoveDown()
	if it.selected != colls[0].Children[0].Children[0] {
		t.Fatal("second descendent not selected")
	}
}

func TestItemTree_MoveDown_ThreeCollsOneDescendentEach(t *testing.T) {
	it := setUp()
	colls := make([]*postman.Collection, 3)
	colls[0] = new(postman.Collection)
	colls[0].AddChild(new(postman.Collection))
	colls[1] = new(postman.Collection)
	colls[1].AddChild(new(postman.Collection))
	colls[2] = new(postman.Collection)
	colls[2].AddChild(new(postman.Collection))
	it.SetCollections(colls)

	it.ToggleExpanded()
	it.MoveDown()
	if it.selected != colls[0].Children[0] {
		t.Fatal("Not p 0 c 0")
	}
	it.MoveDown()
	if it.selected != colls[1] {
		t.Fatal("Not p 1")
	}
	it.ToggleExpanded()
	it.MoveDown()
	if it.selected != colls[1].Children[0] {
		t.Fatal("Not p 1 c 0")
	}
	it.MoveDown()
	if it.selected != colls[2] {
		t.Fatal("Not p 2")
	}
	it.ToggleExpanded()
	it.MoveDown()
	if it.selected != colls[2].Children[0] {
		t.Fatal("Not p 2 c 0")
	}

}