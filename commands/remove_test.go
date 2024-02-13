// Copyright (c) Forge4Flow DAO LLC 2024. All rights reserved.
// Licensed under the MIT license.

package commands

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/forge4flow/forge-cli/test"
)

const testStack = `
provider:
  name: faas
  gateway: http://127.0.0.1:8080

functions:
  fn1:
    lang: go
    handler: ./fn1
`

func Test_remove(t *testing.T) {
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodDelete,
			Uri:                "/system/functions",
			ResponseStatusCode: http.StatusOK,
		},
	})
	defer s.Close()

	// create a yaml stack with the function 'fn1'
	tmpfile, err := ioutil.TempFile("", "stack.*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(testStack)); err != nil {
		tmpfile.Close()
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	resetForTest()

	// run delete with a yaml file and also specify a function to delete.
	// the explicitly specified function should be preferred over the function from the yaml file.
	forgeCmd.SetArgs([]string{
		"remove",
		"--yaml=" + tmpfile.Name(),
		"--gateway=" + s.URL,
		"test-function",
	})
	commandOutput := test.CaptureStdout(func() { forgeCmd.Execute() })

	if !strings.Contains(commandOutput, "Deleting: test-function") {
		t.Error("test-function should be deleted.")
	}
}
