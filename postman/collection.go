package postman

import (
	"encoding/json"
	"github.com/cacarpenter/gopostal/gpmodel"
	"io/ioutil"
	"strings"
)

const POSTMAN_COLLECTION_SUFFIX = "postman_collection.json"

type CollectionInfo struct {
	PostmanId string `json:"_postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
}

type Collection struct {
	Name     string          `json:"name"`
	Info     *CollectionInfo `json:"info"`
	Children []*Collection   `json:"item"`
	Events   []Event         `json:"event"`
	Request  *Request        `json:"request"`
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

	group := gpmodel.Group{}
	wireCollection(&group, &coll)
	// group.LinkParent(nil)

	return &group, nil
}

func wireCollection(g *gpmodel.Group, pc *Collection) {
	g.Name = pc.Label()
	for _, childColl := range pc.Children {
		childGroup := gpmodel.Group{}
		childGroup.LinkParent(g)
		wireCollection(&childGroup, childColl)
		g.AddChild(&childGroup)
	}
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
