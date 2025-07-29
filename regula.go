// Package regula manages hierarchical permissions with support for context and
// wildcards.
package regula

import (
	"strings"
)

// Context holds extra info (key-value pairs) to refine permission checks (e.g.,
// {"world": "overworld"}).
type Context map[string]string

// PermissionState shows if a permission was set and its true/false value.
type PermissionState struct {
	// Set is true if the permission exists.
	Set bool
	// Value is the permission's allow/deny state.
	Value bool
}

// PermissionEntry is a permission rule in a Node.
type PermissionEntry struct {
	// Wildcard is true if this permission uses a wildcard like "admin.*".
	Wildcard bool
	// Context limits when this permission applies.
	Context Context
	// Value holds the allow/deny state and if it's set.
	Value PermissionState
}

// Node is a point in the permission tree, holding entries and child nodes.
type Node struct {
	// Entries are permission rules for this node.
	Entries []PermissionEntry
	// Children nodes represent sub-permissions.
	Children map[string]*Node
}

// NewNode creates a new empty permission node.
func NewNode() *Node {
	return &Node{
		Children: make(map[string]*Node),
	}
}

// SetPermission adds or updates a permission in the tree.
//
//	'perm' is the permission string like "admin.command" or "admin.*".
//	'value' is true (allow) or false (deny).
//	'context' is optional extra info to narrow when it applies.
func (n *Node) SetPermission(perm string, value bool, context Context) {
	parts := strings.Split(perm, ".")
	wildcard := len(parts) > 0 && parts[len(parts)-1] == "*"
	if wildcard {
		parts = parts[:len(parts)-1]
	}

	current := n
	for _, part := range parts {
		if current.Children[part] == nil {
			current.Children[part] = NewNode()
		}
		current = current.Children[part]
	}

	entry := PermissionEntry{
		Wildcard: wildcard,
		Context:  context,
		Value:    PermissionState{Set: true, Value: value},
	}
	current.Entries = append(current.Entries, entry)
}

// GetPermission looks up a permission with context and returns its state.
//
// Checks deeper (more specific) nodes first, prefers later entries, and
// respects context matches.
func (n *Node) GetPermission(perm string, currentContext Context) PermissionState {
	parts := strings.Split(perm, ".")
	partsLength := len(parts)
	if partsLength == 0 {
		return PermissionState{}
	}

	path := []*Node{n}
	current := n

	for _, part := range parts {
		next, ok := current.Children[part]
		if !ok {
			break
		}
		path = append(path, next)
		current = next
	}

	// Check from most specific node back to root
	for i := len(path) - 1; i >= 0; i-- {
		node := path[i]
		// Check entries in reverse to respect overrides
		for j := len(node.Entries) - 1; j >= 0; j-- {
			entry := node.Entries[j]

			// Match exact or wildcard at this level, and check context
			//  (partsLength == i && !entry.Wildcard): Exact match at this level (e.g., "admin.command" at the "command" node).
			//  (partsLength > i && entry.Wildcard)  : Wildcard match for a parent path (e.g., "admin.*" for "admin.command").
			if ((partsLength == i && !entry.Wildcard) || (partsLength > i && entry.Wildcard)) && isContextMatch(entry.Context, currentContext) {
				return entry.Value
			}
		}
	}

	// No permission set found
	return PermissionState{}
}

// isContextMatch returns true if all key-values in permContext exist and match
// currentContext. Empty permContext means match anything.
func isContextMatch(permContext, currentContext Context) bool {
	if len(permContext) == 0 {
		return true
	}
	for k, v := range permContext {
		if currV, ok := currentContext[k]; !ok || currV != v {
			return false
		}
	}
	return true
}
