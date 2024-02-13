package commands

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/forge4flow/forge-cli/test"
)

func Test_storeDeploy_withNameFlag(t *testing.T) {
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodPut,
			Uri:                "/system/functions",
			ResponseStatusCode: http.StatusOK,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		forgeCmd.SetArgs([]string{
			"store",
			"deploy",
			"figlet",
			"--gateway=" + s.URL,
			"--name=foo",
		})
		forgeCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:Deployed)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}

	if found, err := regexp.MatchString(`(?m:200 OK)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}

	if found, err := regexp.MatchString(`(?m:/function/foo)`, stdOut); err != nil || !found {
		t.Fatalf("Wrong function name (should be `foo`):\n%s", stdOut)
	}

	// cleaning after test
	functionName = ""
}

func Test_storeDeploy_withoutNameFlag(t *testing.T) {
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodPut,
			Uri:                "/system/functions",
			ResponseStatusCode: http.StatusOK,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		forgeCmd.SetArgs([]string{
			"store",
			"deploy",
			"figlet",
			"--gateway=" + s.URL,
		})
		forgeCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:Deployed)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}

	if found, err := regexp.MatchString(`(?m:200 OK)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}

	if found, err := regexp.MatchString(`(?m:/function/foo)`, stdOut); err != nil || found {
		t.Fatalf("Wrong function name (should not be `foo`):\n%s", stdOut)
	}
}
