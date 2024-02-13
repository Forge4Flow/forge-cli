// Copyright (c) Forge4Flow Author(s) 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	templateCmd.AddCommand(templateStoreCmd)
}

// templateStoreCmd allows access to pull and list commands from store
var templateStoreCmd = &cobra.Command{
	Use:   `store [COMMAND]`,
	Short: `Command for pulling and listing templates from store`,
	Long:  `This command provides the list of the templates from the official store by default`,
	Example: `  forge-cli template store list --verbose
  forge-cli template store ls -v
  forge-cli template store pull ruby-http
  forge-cli template store pull --url=https://raw.githubusercontent.com/openfaas/store/master/templates.json`,
}
