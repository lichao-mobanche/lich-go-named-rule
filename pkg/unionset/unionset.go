package unionset

// Deleting an element is not supported

// Element TODO
type Element = interface{}

// NewUnionset TODO
func NewUnionset() *UnionSet {
	return &UnionSet{make(map[Element]Element, 0),
		make(map[Element]int, 0), make(map[Element][]Element, 0)}
}

// UnionSet TODO
type UnionSet struct {
	parent map[Element]Element
	sz     map[Element]int
	groups map[Element][]Element
}

func (u *UnionSet) init(key Element) {
	if _, ok := u.parent[key]; !ok {
		u.parent[key] = key
		u.sz[key] = 0
		u.groups[key] = []Element{key}
	}
}

func (u *UnionSet) findP(key Element) (res Element) {
	u.init(key)
	res = key
	for u.parent[res] != res {
		u.parent[res] = u.parent[u.parent[res]]
		res = u.parent[res]
	}
	return
}

// Join a key-value pair to UnionSet
func (u *UnionSet) Join(key1, key2 Element) {
	k1p := u.findP(key1)
	k2p := u.findP(key2)
	if k1p == k2p {
		return
	}
	if u.sz[k1p] < u.sz[k2p] {
		u.sz[k2p] += u.sz[k1p]
		u.parent[k1p] = k2p
		u.groups[k2p] = append(u.groups[k2p], u.groups[k1p]...)
		delete(u.groups, k1p)
	} else {
		u.sz[k1p] += u.sz[k2p]
		u.parent[k2p] = k1p
		u.groups[k1p] = append(u.groups[k1p], u.groups[k2p]...)
		delete(u.groups, k2p)
	}
	return
}

// GetGroupNumber return the number of groups
func (u UnionSet) GetGroupNumber() int {
	return len(u.groups)
}

// GetGroups return the struct of groups
func (u UnionSet) GetGroups() [][]Element {
	i := 0
	g := make([][]Element, len(u.groups))
	for _, group := range u.groups {
		g[i] = group
		i++
	}
	return g
}
