# Regula

Regula (pronounced REH-gyoo-lah) is a flexible permission system for Go, inspired by the popular Minecraft plugin [LuckPerms](https://luckperms.net/). It supports context-aware permissions and wildcards for fine-grained control.

## Features

- Context-based permissions (e.g., per world or region)
- Wildcard support (e.g., `admin.*`)
- Simple permission trees for users and groups
- Easy-to-use permission trees for users and groups

## Installation

```bash
go get github.com/calumari/regula
```

## Quick Start

1. Define Groups with specific permissions.
2. Create Users and assign them to groups.
3. Check permissions with context-aware queries.

```go
package main

import (
	"fmt"

	"github.com/calumari/regula"
	"github.com/calumari/regula/types"
)

func main() {
	// Example usage:
	vipGroup := types.NewGroup("VIP")
	vipGroup.PermissionTree.SetPermission("essentials.fly", true, regula.Context{})

	user := types.NewUser("PlayerOne")
	user.Groups = append(user.Groups, vipGroup)

	perm := user.HasPermission("essentials.fly", regula.Context{})
	fmt.Println("Can fly?", perm.Value) // Output: Can fly? true
}
```