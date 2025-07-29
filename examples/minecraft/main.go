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
		state := user.GetPermission(test.perm, test.context)
		granted := state.Set && state.Value
		fmt.Printf("  %-25s | Context: %v | Granted: %t (Set:%t Value:%t)\n", test.perm, test.context, granted, state.Set, state.Value)
	}
}

func main() {
	// Guest group: can request teleport, but cannot fly
	guest := types.NewGroup("Guest")
	guest.SetPermission("essentials.tpa", true, regula.Context{})
	guest.SetPermission("essentials.fly", false, regula.Context{})

	// Member group: inherits Guest, can set home, build in overworld but not nether
	member := types.NewGroup("Member")
	member.AddParent(guest)
	member.SetPermission("essentials.sethome", true, regula.Context{})
	member.SetPermission("worldguard.build", true, regula.Context{"world": "overworld"})
	member.SetPermission("worldguard.build", false, regula.Context{"world": "nether"})

	// VIP group: inherits Member, can fly except in nether (deny overrides allow)
	vip := types.NewGroup("VIP")
	vip.AddParent(member)
	vip.SetPermission("essentials.fly", true, regula.Context{})
	vip.SetPermission("essentials.fly", false, regula.Context{"world": "nether"})

	// User with VIP group and a direct override: allow fly in nether despite group deny
	user := types.NewUser("Alex")
	user.AddGroup(vip)
	user.SetPermission("essentials.fly", true, regula.Context{"world": "nether"})

	checkPermissions(user)
}
