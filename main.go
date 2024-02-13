// Copyright (c) Forge4Flow DAO LLC 2024. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"fmt"
	"os"

	"github.com/forge4flow/forge-cli/commands"
)

func main() {
	customArgs, err := translateLegacyOpts(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	commands.Execute(customArgs)
}
