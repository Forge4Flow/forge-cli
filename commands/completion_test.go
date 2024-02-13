package commands

import (
	"testing"
)

func Test_ValidShell(t *testing.T) {

	testArgs := [][]string{
		{"completion", "--shell", "bash"},
		{"completion", "--shell", "zsh"},
	}

	for _, arg := range testArgs {
		forgeCmd.SetArgs(arg)

		err := forgeCmd.Execute()

		if err != nil {
			t.Errorf("err was supposed to be nil but it was: %s", err)
			t.Fail()
		}
	}
}
