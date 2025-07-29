// Package types provides basic models like User and Group that work with regula
// permissions. These models show how to build user and group permissions using
// regula.Node.
package types

import (
	"github.com/calumari/regula"
)

// Group holds permissions and can inherit from parent groups.
type Group struct {
	// Name is the group's identifier or display name.
	Name string
	// PermissionTree stores this group's own permissions.
	PermissionTree *regula.Node
	// Parents are other groups this group inherits permissions from.
	Parents []*Group
}

// NewGroup creates a new Group with the given name.
func NewGroup(name string) *Group {
	return &Group{
		Name:           name,
		PermissionTree: regula.NewNode(),
		Parents:        []*Group{},
	}
}

// GetPermission checks this group’s permissions, then checks parent groups if
// not found.
func (g *Group) GetPermission(perm string, context regula.Context) regula.PermissionState {
	if val := g.PermissionTree.GetPermission(perm, context); val.Set {
		return val
	}
	for _, parent := range g.Parents {
		if val := parent.GetPermission(perm, context); val.Set {
			return val
		}
	}
	return regula.PermissionState{}
}

// User represents an individual with permissions and group memberships.
type User struct {
	// Name is the user's identifier or display name.
	Name string
	// PermissionTree stores the user's own permissions.
	PermissionTree *regula.Node
	// Groups the user belongs to.
	Groups []*Group
}

// NewUser creates a new User with the given name.
func NewUser(name string) *User {
	return &User{
		Name:           name,
		PermissionTree: regula.NewNode(),
		Groups:         []*Group{},
	}
}

// HasPermission checks if the user has a permission. Checks the user's own
// permissions first, then group permissions.
func (u *User) HasPermission(perm string, context regula.Context) regula.PermissionState {
	if val := u.PermissionTree.GetPermission(perm, context); val.Set {
		return val
	}
	for _, group := range u.Groups {
		if val := group.GetPermission(perm, context); val.Set {
			return val
		}
	}
	return regula.PermissionState{}
}
