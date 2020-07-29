package postman

import (
	"encoding/json"
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
	if pc.Request != nil {
		if parent == nil {
			log.Panicln("Cannot set request on nil parent")
		}
		reqSpec := NewRequestSpec(pc.Request)
		reqSpec.Name = pc.Label()
		parent.Requests = append(parent.Requests, reqSpec)
		return parent
	}
	var p *gpmodel.Group
	if parent == nil {
		p = new(gpmodel.Group)
	} else {
		p = parent
	}
	p.Name = pc.Label()
	for _, childColl := range pc.Items {
		if childColl.Request != nil {
			wireCollection(p, childColl)
		} else {
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
