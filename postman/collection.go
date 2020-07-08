package postman

import (
	"encoding/json"
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
	Name               string          `json:"name"`
	Info               *CollectionInfo `json:"info"`
	Children           []*Collection   `json:"item"`
	Events             []Event         `json:"event"`
	Request            *Request        `json:"request"`
	expanded, selected bool
	parent             *Collection
}

func (c *Collection) Expanded() bool {
	return len(c.Children) > 0 && c.expanded
}

func (c *Collection) SetExpanded(e bool) {
	c.expanded = e
}

func (c *Collection) Expand(exp, recursive bool) {
	c.expanded = exp
	if recursive {
		for _, ch := range c.Children {
			ch.Expand(exp, recursive)
		}
	}
}

func (c *Collection) ToggleExpanded() bool {
	c.expanded = !c.expanded
	return c.expanded
}

func (c *Collection) Selected() bool {
	return c.selected
}

func (c *Collection) ToggleSelected() bool {
	c.selected = !c.selected
	return c.selected
}

func (c *Collection) SetSelected(b bool) {
	c.selected = b
}

func (c *Collection) PreviousSibling() *Collection {
	// this is the root
	if c.parent == nil {
		return nil
	}
	var prev *Collection
	for _, ch := range c.parent.Children {
		if c == ch {
			break
		}
		prev = ch
	}
	return prev
}

func (c *Collection) NextSibling() *Collection {
	// root
	if c.parent == nil {
		return nil
	}
	numSibs := len(c.parent.Children)
	for i, ch := range c.parent.Children {
		if c == ch && i < numSibs-1 {
			return c.parent.Children[i+1]
		}
	}
	return nil
}

func (c *Collection) Parent() *Collection {
	return c.parent
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

	// link up the parents
	coll.linkParent(nil)
	coll.ToggleSelected()

	return &coll, nil
}

func (c *Collection) linkParent(p *Collection) {
	for _, ch := range c.Children {
		ch.linkParent(c)
	}
	c.parent = p
}

func (c *Collection) ParentName() string {
	if c.parent != nil {
		return c.parent.Name
	}
	return ""
}

func (c *Collection) Label() string {
	if c.Info != nil {
		return c.Info.Name
	}
	return c.Name
}

func (c *Collection) AddChild(child *Collection) {
	child.parent = c
	c.Children = append(c.Children, child)
}

// LastExpandedDescendent recursively returns the last child of expanded of a collection.
// Otherwise the collection itself is returned
func (c *Collection) LastExpandedDescendent() *Collection {
	if c.expanded && len(c.Children) > 0 {
		return c.Children[len(c.Children)-1].LastExpandedDescendent()
	}
	return c
}

func IsCollectionFile(filename string) bool {
	return strings.HasSuffix(filename, POSTMAN_COLLECTION_SUFFIX)
}
