package postman

import (
	"encoding/json"
	"fmt"
	"github.com/cacarpenter/gopostal/gpmodel"
	"io/ioutil"
	"log"
	"strings"
)

const POSTMAN_COLLECTION_SUFFIX = "postman_collection.json"

type CollectionInfo struct {
	PostmanId string `json:"_postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
}

type Collection struct {
	Name    string          `json:"name"`
	Info    *CollectionInfo `json:"info"`
	Items   []*Collection   `json:"item"`
	Events  []Event         `json:"event"`
	Request *Request        `json:"request"`
}

func ParseCollection(filename string) (*gpmodel.Group, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var coll Collection
	err = json.Unmarshal(data, &coll)
	if err != nil {
		return nil, err
	}

	return wireCollection(nil, &coll), nil
}

func wireCollection(parent *gpmodel.Group, pc *Collection) *gpmodel.Group {
	fmt.Println("Setting group name to", pc.Label())
	if pc.Request != nil {
		fmt.Println("Found request", pc.Request)
		if parent == nil {
			log.Panicln("Cannot set request on nil parent")
		}
		fmt.Println("Set request label to parent name", parent.Name)
		reqSpec := NewRequestSpec(pc.Request)
		reqSpec.Name = parent.Name
		parent.Requests = append(parent.Requests, reqSpec)
		return parent
	}
	fmt.Println("No request on ", pc.Label(), " look for items")
	var p *gpmodel.Group
	if parent == nil {
		fmt.Println("Creating new parent")
		p = new(gpmodel.Group)
	} else {
		fmt.Println("Using parent ", parent.Name)
		p = parent
	}
	p.Name = pc.Label()
	fmt.Println("Recur on ", len(pc.Items))
	for _, childColl := range pc.Items {
		if childColl.Request != nil {
			fmt.Println("Found request on ", childColl.Label())
			wireCollection(p, childColl)
		} else {
			fmt.Println("New Group for ", childColl.Name)
			grp := new(gpmodel.Group)
			wireCollection(grp, childColl)
			p.AddChild(grp)
		}
	}
	return p
}

func (c *Collection) Label() string {
	if c.Info != nil {
		return c.Info.Name
	}
	return c.Name
}

func IsCollectionFile(filename string) bool {
	return strings.HasSuffix(filename, POSTMAN_COLLECTION_SUFFIX)
}
