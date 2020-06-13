package postman

import (
	"encoding/json"
	"io/ioutil"
)

type CollectionInfo struct {
	PostmanId string `json:"_postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
}

type Collection struct {
	Name               string         `json:"name"`
	Info               CollectionInfo `json:"info"`
	Children           []Collection   `json:"item"`
	Events             []Event        `json:"event"`
	Request            *Request       `json:"request"`
	expanded, selected bool
}

func (ci *Collection) Expanded() bool {
	return ci.expanded
}

func (ci *Collection) ToggleExpanded() bool {
	ci.expanded = !ci.expanded
	return ci.expanded
}

func (ci *Collection) Selected() bool {
	return ci.selected
}

func (ci *Collection) ToggleSelected() bool {
	ci.selected = !ci.selected
	return ci.selected
}

func ParseCollection(filename string) (*Collection, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var coll Collection
	err = json.Unmarshal(data, &coll)
	if err != nil {
		return nil, err
	}
	return &coll, nil
}
