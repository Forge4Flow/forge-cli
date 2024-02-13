// Copyright (c) Forge4Flow DAO LLC 2024. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"reflect"
	"testing"
)

var translateLegacyOptsTests = []struct {
	title        string
	inputArgs    []string
	expectedArgs []string
	expectError  bool
}{
	{
		title:        "legacy deploy action with all args, no =",
		inputArgs:    []string{"forge-cli", "-action", "deploy", "-image", "testimage", "-name", "fnname", "-fprocess", `"/usr/bin/faas-img2ansi"`, "-gateway", "https://url", "-handler", "/dir/", "-lang", "python", "-replace"},
		expectedArgs: []string{"forge-cli", "deploy", "--image", "testimage", "--name", "fnname", "--fprocess", `"/usr/bin/faas-img2ansi"`, "--gateway", "https://url", "--handler", "/dir/", "--lang", "python", "--replace"},
		expectError:  false,
	},
	{
		title:        "legacy deploy action with =",
		inputArgs:    []string{"forge-cli", "-action=deploy", "-image=testimage", "-name=fnname", `-fprocess="/usr/bin/faas-img2ansi"`},
		expectedArgs: []string{"forge-cli", "deploy", "--image=testimage", "--name=fnname", `--fprocess="/usr/bin/faas-img2ansi"`},
		expectError:  false,
	},
	{
		title:        "legacy deploy action with -f",
		inputArgs:    []string{"forge-cli", "-action=deploy", "-f", "/dir/file.yml"},
		expectedArgs: []string{"forge-cli", "deploy", "-f", "/dir/file.yml"},
		expectError:  false,
	},
	{
		title:        "legacy deploy action with -yaml",
		inputArgs:    []string{"forge-cli", "-action=deploy", "-yaml", "/dir/file.yml"},
		expectedArgs: []string{"forge-cli", "deploy", "--yaml", "/dir/file.yml"},
		expectError:  false,
	},
	{
		title:        "legacy build action with all args, no =",
		inputArgs:    []string{"forge-cli", "-action", "build", "-image", "testimage", "-name", "fnname", "-handler", "/dir/", "-lang", "python", "-no-cache", "-squash"},
		expectedArgs: []string{"forge-cli", "build", "--image", "testimage", "--name", "fnname", "--handler", "/dir/", "--lang", "python", "--no-cache", "--squash"},
		expectError:  false,
	},
	{
		title:        "legacy delete action (note delete->remove translation)",
		inputArgs:    []string{"forge-cli", "-action", "delete", "-name", "fnname"},
		expectedArgs: []string{"forge-cli", "remove", "fnname"},
		expectError:  false,
	},
	{
		title:        "legacy delete action with yaml",
		inputArgs:    []string{"forge-cli", "-action", "delete", "-f", "/dir/file.yml"},
		expectedArgs: []string{"forge-cli", "remove", "-f", "/dir/file.yml"},
		expectError:  false,
	},
	{
		title:        "legacy version flag",
		inputArgs:    []string{"forge-cli", "-version"},
		expectedArgs: []string{"forge-cli", "version"},
		expectError:  false,
	},
	{
		title:        "version command",
		inputArgs:    []string{"forge-cli", "version"},
		expectedArgs: []string{"forge-cli", "version"},
		expectError:  false,
	},
	{
		title:        "deploy command",
		inputArgs:    []string{"forge-cli", "deploy", "--image", "testimage", "--name", "fnname", "--fprocess", `"/usr/bin/faas-img2ansi"`, "--gateway", "https://url", "--handler", "/dir/", "--lang", "python", "--replace", "--env", "KEY1=VAL1", "--env", "KEY2=VAL2"},
		expectedArgs: []string{"forge-cli", "deploy", "--image", "testimage", "--name", "fnname", "--fprocess", `"/usr/bin/faas-img2ansi"`, "--gateway", "https://url", "--handler", "/dir/", "--lang", "python", "--replace", "--env", "KEY1=VAL1", "--env", "KEY2=VAL2"},
		expectError:  false,
	},
	{
		title:        "build command",
		inputArgs:    []string{"forge-cli", "build", "--image", "testimage", "--name", "fnname", "--handler", "/dir/", "--lang", "python", "--no-cache", "--squash"},
		expectedArgs: []string{"forge-cli", "build", "--image", "testimage", "--name", "fnname", "--handler", "/dir/", "--lang", "python", "--no-cache", "--squash"},
		expectError:  false,
	},
	{
		title:        "remove command",
		inputArgs:    []string{"forge-cli", "remove", "fnname"},
		expectedArgs: []string{"forge-cli", "remove", "fnname"},
		expectError:  false,
	},
	{
		title:        "remove command alias rm",
		inputArgs:    []string{"forge-cli", "rm", "fnname"},
		expectedArgs: []string{"forge-cli", "rm", "fnname"},
		expectError:  false,
	},
	{
		title:        "remove command alias delete",
		inputArgs:    []string{"forge-cli", "delete", "fnname"},
		expectedArgs: []string{"forge-cli", "delete", "fnname"},
		expectError:  false,
	},
	{
		title:        "push command",
		inputArgs:    []string{"forge-cli", "delete", "fnname"},
		expectedArgs: []string{"forge-cli", "delete", "fnname"},
		expectError:  false,
	},
	{
		title:        "bashcompletion command",
		inputArgs:    []string{"forge-cli", "bashcompletion", "/dir/file"},
		expectedArgs: []string{"forge-cli", "bashcompletion", "/dir/file"},
		expectError:  false,
	},
	{
		title:        "legacy flag as value without =",
		inputArgs:    []string{"forge-cli", "-action", "deploy", "-name", `"-name"`},
		expectedArgs: []string{"forge-cli", "deploy", "--name", `"-name"`},
		expectError:  false,
	},
	{
		title:        "legacy flag as value with =",
		inputArgs:    []string{"forge-cli", "-action", "deploy", "-name=-name"},
		expectedArgs: []string{"forge-cli", "deploy", "--name=-name"},
		expectError:  false,
	},
	{
		title:        "unknown legacy flag",
		inputArgs:    []string{"forge-cli", "-action", "deploy", "-fe"},
		expectedArgs: []string{"forge-cli", "deploy", "-fe"},
		expectError:  false,
	},
	{
		title:        "legacy -action missing value",
		inputArgs:    []string{"forge-cli", "-action"},
		expectedArgs: []string{""},
		expectError:  true,
	},
	{
		title:        "legacy -action= missing value",
		inputArgs:    []string{"forge-cli", "-action="},
		expectedArgs: []string{""},
		expectError:  true,
	},
	{
		title:        "legacy -action with unknown value",
		inputArgs:    []string{"forge-cli", "-action", "unknownaction"},
		expectedArgs: []string{""},
		expectError:  true,
	},
	{
		title:        "legacy -action= with unknown value",
		inputArgs:    []string{"forge-cli", "-action=unknownaction"},
		expectedArgs: []string{""},
		expectError:  true,
	},
}

func Test_translateLegacyOpts(t *testing.T) {
	for _, test := range translateLegacyOptsTests {
		t.Run(test.title, func(t *testing.T) {
			actual, err := translateLegacyOpts(test.inputArgs)
			if test.expectError {
				if err == nil {
					t.Errorf("TranslateLegacyOpts test [%s] test failed, expected error not thrown", test.title)
					return
				}
			} else {
				if err != nil {
					t.Errorf("TranslateLegacyOpts test [%s] test failed, unexpected error thrown", test.title)
					return
				}
			}
			if !reflect.DeepEqual(actual, test.expectedArgs) {
				t.Errorf("TranslateLegacyOpts test [%s] test failed, does not match expected result;\n  actual:   [%v]\n  expected: [%v]",
					test.title,
					actual,
					test.expectedArgs,
				)
			}
		})
	}
}
