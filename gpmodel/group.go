package gpmodel

type Group struct {
	Name               string
	Children           []*Group
	parent             *Group
	Request            *RequestSpec
	expanded, selected bool
}

func (g *Group) Expanded() bool {
	return len(g.Children) > 0 && g.expanded
}

func (g *Group) SetExpanded(e bool) {
	g.expanded = e
}

func (g *Group) Expand(exp, recursive bool) {
	g.expanded = exp
	if recursive {
		for _, ch := range g.Children {
			ch.Expand(exp, recursive)
		}
	}
}

func (g *Group) ToggleExpanded() bool {
	g.expanded = !g.expanded
	return g.expanded
}

func (g *Group) Selected() bool {
	return g.selected
}

func (g *Group) ToggleSelected() bool {
	g.selected = !g.selected
	return g.selected
}

func (g *Group) SetSelected(b bool) {
	g.selected = b
}

func (g *Group) PreviousSibling() *Group {
	// this is the root
	if g.parent == nil {
		return nil
	}
	var prev *Group
	for _, ch := range g.parent.Children {
		if g == ch {
			break
		}
		prev = ch
	}
	return prev
}

func (g *Group) NextSibling() *Group {
	// root
	if g.parent == nil {
		return nil
	}
	numSibs := len(g.parent.Children)
	for i, ch := range g.parent.Children {
		if g == ch && i < numSibs-1 {
			return g.parent.Children[i+1]
		}
	}
	return nil
}

func (g *Group) Parent() *Group {
	return g.parent
}

func (g *Group) LinkParent(p *Group) {
	for _, ch := range p.Children {
		ch.LinkParent(g)
	}
	g.parent = p
}

/*
// LastExpandedDescendent recursively returns the last child of expanded of a collection.
// Otherwise the collection itself is returned
func (g *Group) LastExpandedDescendent() *Group {
	if g.expanded && len(g.Children) > 0 {
		return g.Children[len(g.Children)-1].LastExpandedDescendent()
	}
	return g
}

*/

/*
func (g *Group) ParentName() string {
	if g.parent != nil {
		return g.parent.Name
	}
	return ""
}
*/

// used by tests
func (g *Group) AddChild(child *Group) {
	child.parent = g
	g.Children = append(g.Children, child)
}

