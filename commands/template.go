// Copyright (c) Forge4Flow Author(s) 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	forgeCmd.AddCommand(templateCmd)
}

// templateCmd allows access to store and pull commands
var templateCmd = &cobra.Command{
	Use:   `template [COMMAND]`,
	Short: "Forge4Flow template store and pull commands",
	Long:  "Allows browsing templates from store or pulling custom templates",
	Example: `  forge-cli template pull https://github.com/custom/template
  forge-cli template store list
  forge-cli template store ls
  forge-cli template store pull ruby-http
  forge-cli template store pull openfaas-incubator/ruby-http`,
}
