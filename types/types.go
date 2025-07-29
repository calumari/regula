package types

import (
	"github.com/calumari/regula"
)

// Group holds permissions and can inherit from parent groups.
type Group struct {
	name           string
	permissionTree *regula.Node
	parents        []*Group
}

// NewGroup creates a new Group with the given name.
func NewGroup(name string) *Group {
	return &Group{
		name:           name,
		permissionTree: regula.NewNode(),
		parents:        []*Group{},
	}
}

// GetPermission checks this group’s permissions, then checks parent groups if
// not found.
func (g *Group) GetPermission(perm string, context regula.Context) regula.PermissionState {
	if val := g.permissionTree.GetPermission(perm, context); val.Set {
		return val
	}
	for _, parent := range g.parents {
		if val := parent.GetPermission(perm, context); val.Set {
			return val
		}
	}
	return regula.PermissionState{}
}

// SetPermission adds or updates a permission in this group's permission tree.
func (g *Group) SetPermission(perm string, value bool, context regula.Context) {
	g.permissionTree.SetPermission(perm, value, context)
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) Parents() []*Group {
	return g.parents
}

func (g *Group) AddParent(parent *Group) {
	g.parents = append(g.parents, parent)
}

// User represents an individual with permissions and group memberships.
type User struct {
	name           string
	permissionTree *regula.Node
	groups         []*Group
}

// NewUser creates a new User with the given name.
func NewUser(name string) *User {
	return &User{
		name:           name,
		permissionTree: regula.NewNode(),
		groups:         []*Group{},
	}
}

// GetPermission checks if the user has a permission. Checks the user's own
// permissions first, then group permissions.
func (u *User) GetPermission(perm string, context regula.Context) regula.PermissionState {
	if val := u.permissionTree.GetPermission(perm, context); val.Set {
		return val
	}
	for _, group := range u.groups {
		if val := group.GetPermission(perm, context); val.Set {
			return val
		}
	}
	return regula.PermissionState{}
}

// SetPermission adds or updates a permission in the user's permission tree.
func (u *User) SetPermission(perm string, value bool, context regula.Context) {
	u.permissionTree.SetPermission(perm, value, context)
}

func (u *User) HasPermission(perm string, context regula.Context) bool {
	state := u.GetPermission(perm, context)
	return state.Set && state.Value
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Groups() []*Group {
	return u.groups
}

func (u *User) AddGroup(group *Group) {
	u.groups = append(u.groups, group)
}
