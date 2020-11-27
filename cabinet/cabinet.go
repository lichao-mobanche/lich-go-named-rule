package cabinet

import (
	"github.com/lichao-mobanche/lich-go-named-rule/pkg/unionset"
	"sync"
)

type groupName = string
type tageName = string
type arrow = string
type rule interface{}
type tage struct {
	rule      rule
	nextTages map[arrow]tageName
}

// NewGroup TODO
func NewGroup(gname groupName) *Group {
	if gname == "" {
		return nil
	}
	return &Group{gname, make(map[tageName]tage), unionset.NewUnionset(), Complete, &sync.RWMutex{}}
}

// Group TODO
type Group struct {
	gname  groupName
	tags   map[tageName]tage
	graph  *unionset.UnionSet
	status Status
	*sync.RWMutex
}

// LoadTage creates a new Group
func (g *Group) LoadTage(t tageName, r interface{}) (tageName, interface{}) {
	g.Lock()
	defer g.Unlock()
	if _, ok := g.tags[t]; !ok {
		g.tags[t] = tage{
			r,
			make(map[arrow]tageName),
		}
		g.graph.Join(t, t)
		g.setStatus(false)
		return t, r
	}
	old := g.tags[t]
	g.tags[t] = tage{
		r,
		old.nextTages,
	}
	return t, old.rule
}

func (g Group) GetTage(t tageName) (r interface{}) {
	g.RLock()
	defer g.RUnlock()
	if tag, ok := g.tags[t]; ok {
		r = tag.rule
	}
	return
}

func (g *Group) rebuild() {
	new := unionset.NewUnionset()
	defects := false
	for tn, tag := range g.tags {
		for _, st := range tag.nextTages {
			new.Join(tn, st)
			if _, ok := g.tags[st]; !ok {
				defects = true
			}
		}
	}
	g.setStatus(defects)
	g.graph = new
}

func (g *Group) RemoveTage(t tageName) {
	g.Lock()
	defer g.Unlock()
	if _, ok := g.tags[t]; ok {
		delete(g.tags, t)
		g.rebuild()
	}
}
func (g *Group) setStatus(detective bool) {
	if detective {
		g.status = Incomplete
		return
	}
	l := g.graph.GetGroupNumber()
	if l == 1 {
		g.status = Complete
	} else if l > 1 {
		g.status = Incomplete
	}
}
func (g *Group) LoadSubTage(t tageName, a arrow, st tageName) (tageName, arrow, tageName) {
	g.Lock()
	defer g.Unlock()
	if tag, ok := g.tags[t]; ok {
		resst := st
		if _, ok := tag.nextTages[a]; ok {
			resst = tag.nextTages[a]
		}
		tag.nextTages[a] = st
		g.graph.Join(t, st)
		_, ok := g.tags[st]
		g.setStatus(!ok)
		return t, a, resst
	}
	return "", "", ""
}

func (g *Group) RemoveSubTage(t tageName, a arrow) {
	g.Lock()
	defer g.Unlock()
	if tag, ok := g.tags[t]; ok {
		delete(tag.nextTages, a)
		g.rebuild()
	}
}

func (g Group) GetSubTag(t tageName) interface{} {
	g.RLock()
	defer g.RUnlock()
	res := Result{}
	if t != "" {
		tag, ok := g.tags[t]
		if !ok {
			return nil
		}
		res[t] = tag.nextTages
	} else {
		for k, tag := range g.tags {
			res[k] = tag.nextTages
		}
	}
	return res
}

func (g Group) CheckNextTage(t tageName, a arrow) tageName {
	if tag, ok := g.tags[t]; ok {
		if next, ok := tag.nextTages[a]; ok {
			return next
		}
	}
	return ""
}

func (g *Group) LoadGroupRule(r interface{}) interface{} {
	g.Lock()
	defer g.Unlock()
	if _, ok := g.tags[GroupTag]; !ok {
		g.tags[GroupTag] = tage{
			r,
			nil,
		}
		return r
	}
	old := g.tags[GroupTag]
	g.tags[GroupTag] = tage{
		r,
		nil,
	}
	return old.rule
}

func (g Group) GetGroupRule() (r interface{}) {
	return g.GetTage(GroupTag)
}

func (g *Group) RemoveGroupRule() {
	g.Lock()
	defer g.Unlock()
	if _, ok := g.tags[GroupTag]; ok {
		delete(g.tags, GroupTag)
		g.rebuild()
	}
}

func (g *Group) GroupInfo() Result {
	g.RLock()
	defer g.RUnlock()
	res := Result{}
	res["status"] = g.status
	res["graph"] = g.graph.GetGroups()
	return res
}

func NewCabinet() *Cabinet {
	return &Cabinet{make(map[groupName]*Group), 0, &sync.RWMutex{}}
}

// Cabinet TODO
type Cabinet struct {
	unit map[groupName]*Group
	size int
	*sync.RWMutex
}

func (c Cabinet) Size() int {
	return c.size
}

func (c *Cabinet) LoadGroup(group *Group) groupName {
	c.Lock()
	defer c.Unlock()
	c.unit[group.gname] = group
	return group.gname
}

func (c *Cabinet) RemoveGroup(g groupName) groupName {
	c.Lock()
	defer c.Unlock()
	res := ""
	if _, ok := c.unit[g]; ok {
		delete(c.unit, g)
		c.size--
		res = g
	}
	return res
}

func (c Cabinet) GetGroup(g groupName) *Group {
	c.RLock()
	defer c.RUnlock()
	if group, ok := c.unit[g]; ok {
		return group
	}
	return nil
}
