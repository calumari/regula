package main

import (
	"fmt"

	"github.com/calumari/regula"
	"github.com/calumari/regula/types"
)

// checkPermissions prints permission checks with context info
func checkPermissions(user *types.User) {
	fmt.Printf("\nPermissions for %s:\n", user.Name)

	tests := []struct {
		perm    string
		context regula.Context
	}{
		{"essentials.tpa", regula.Context{}},
		{"essentials.sethome", regula.Context{}},
		{"essentials.fly", regula.Context{}},
		{"essentials.fly", regula.Context{"world": "nether"}},
		{"worldguard.build", regula.Context{"world": "overworld"}},
		{"worldguard.build", regula.Context{"world": "nether"}},
		{"bedwars.join", regula.Context{}},
	}

	for _, test := range tests {
		state := user.HasPermission(test.perm, test.context)
		granted := state.Set && state.Value
		fmt.Printf("  %-25s | Context: %v | Granted: %t (Set:%t Value:%t)\n",
			test.perm, test.context, granted, state.Set, state.Value)
	}
}

func main() {
	// Guest group: can request teleport, but cannot fly
	guest := types.NewGroup("Guest")
	guest.PermissionTree.SetPermission("essentials.tpa", true, regula.Context{})
	guest.PermissionTree.SetPermission("essentials.fly", false, regula.Context{})

	// Member group: inherits Guest, can set home, build in overworld but not nether
	member := types.NewGroup("Member")
	member.Parents = append(member.Parents, guest)
	member.PermissionTree.SetPermission("essentials.sethome", true, regula.Context{})
	member.PermissionTree.SetPermission("worldguard.build", true, regula.Context{"world": "overworld"})
	member.PermissionTree.SetPermission("worldguard.build", false, regula.Context{"world": "nether"})

	// VIP group: inherits Member, can fly except in nether (deny overrides allow)
	vip := types.NewGroup("VIP")
	vip.Parents = append(vip.Parents, member)
	vip.PermissionTree.SetPermission("essentials.fly", true, regula.Context{})
	vip.PermissionTree.SetPermission("essentials.fly", false, regula.Context{"world": "nether"})

	// User with VIP group and a direct override: allow fly in nether despite group deny
	user := types.NewUser("Alex")
	user.Groups = append(user.Groups, vip)
	user.PermissionTree.SetPermission("essentials.fly", true, regula.Context{"world": "nether"})

	checkPermissions(user)
}
