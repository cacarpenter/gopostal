package gpmodel

type Group struct {
	Name     string
	Parent   *Group
	Children []*Group
	Requests []*RequestSpec
}

func (g *Group) PreviousSibling() *Group {
	// this is the root
	if g.Parent == nil {
		return nil
	}
	var prev *Group
	for _, ch := range g.Parent.Children {
		if g == ch {
			break
		}
		prev = ch
	}
	return prev
}

func (g *Group) NextSibling() *Group {
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

func (g *Group) LinkParent(p *Group) {
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
func (g *Group) LastExpandedDescendent() *Group {
	if g.expanded && len(g.Items) > 0 {
		return g.Items[len(g.Items)-1].LastExpandedDescendent()
	}
	return g
}*/

// used by tests
func (g *Group) AddChild(child *Group) {
	child.Parent = g
	g.Children = append(g.Children, child)
}

func NewGroup(name string) *Group {
	g := new(Group)
	g.Name = name
	return g
}
