// Copyright (c) Forge4Flow DAO LLC 2024. All rights reserved.
// Licensed under the MIT license.

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	forgeCmd.AddCommand(bashcompletionCmd)
}

// bashcompletionCmd generates a bash completion file
// TODO split into `completion bash`/`completion zsh`?
var bashcompletionCmd = &cobra.Command{
	Use:   "bashcompletion FILENAME",
	Short: "Generate a bash completion file",
	Long: `Generate a bash completion file for the client.

This currently only works on Bash version 4, and is hidden
pending a merge of https://github.com/spf13/cobra/pull/520.`,
	Hidden:     true,
	Deprecated: `please use the "completion" command`,
	RunE:       runBashcompletion,
}

func runBashcompletion(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide filename for bash completion")
	}
	fileName := args[0]
	err := forgeCmd.GenBashCompletionFile(fileName)
	if err != nil {
		return fmt.Errorf("unable to create bash completion file")
	}

	return nil
}
