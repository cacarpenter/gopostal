package gpmodel

type RequestGroup struct {
	Name     string
	Parent   *RequestGroup
	Children []*RequestGroup
	Requests []*RequestSpec
}

func (g *RequestGroup) PreviousSibling() *RequestGroup {
	// this is the root
	if g.Parent == nil {
		return nil
	}
	var prev *RequestGroup
	for _, ch := range g.Parent.Children {
		if g == ch {
			break
		}
		prev = ch
	}
	return prev
}

func (g *RequestGroup) NextSibling() *RequestGroup {
	// root
	if g.Parent == nil {
		return nil
	}
	numSibs := len(g.Parent.Children)
	for i, ch := range g.Parent.Children {
		if g == ch && i < numSibs-1 {
			return g.Parent.Children[i+1]
		}
	}
	return nil
}

func (g *RequestGroup) LinkParent(p *RequestGroup) {
	if p != nil {
		for _, ch := range p.Children {
			ch.LinkParent(g)
		}
	}
	g.Parent = p
}

/*
// LastExpandedDescendent recursively returns the last child of expanded of a collection.
// Otherwise the collection itself is returned
func (g *RequestGroup) LastExpandedDescendent() *RequestGroup {
	if g.expanded && len(g.Items) > 0 {
		return g.Items[len(g.Items)-1].LastExpandedDescendent()
	}
	return g
}*/

// used by tests
func (g *RequestGroup) AddChild(child *RequestGroup) {
	child.Parent = g
	g.Children = append(g.Children, child)
}

func NewGroup(name string) *RequestGroup {
	g := new(RequestGroup)
	g.Name = name
	return g
}
