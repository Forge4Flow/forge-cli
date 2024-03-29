// Copyright (c) Forge4Flow Author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Test_GenerateRegistryAuth(t *testing.T) {
	registryURL := "https://index.docker.io/v1/"
	username := "docker_user"
	password := "docker_password"
	gotBytes, _ := generateRegistryAuth(registryURL, username, password)
	got, err := bytesToRegistryStruct(gotBytes)

	if err != nil {
		t.Errorf("Error converting bytes to struct, %q", err)
		t.Fail()
	}

	want := RegistryAuth{
		AuthConfigs: map[string]Auth{
			registryURL: {Base64AuthString: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Structs were not equal, want: %q\ngot:%q", want, got)
		t.Fail()
	}
}

func Test_GenerateRegistryAuthNoRegistryURL(t *testing.T) {
	registryURL := ""
	username := "docker_user"
	password := "docker_password"
	_, err := generateRegistryAuth(registryURL, username, password)

	if err == nil {
		t.Error("Want err, got nil")
		t.Fail()
	}
}

func Test_GenerateRegistryAuthNoUsername(t *testing.T) {
	registryURL := "https://index.docker.io/v1"
	username := ""
	password := "docker_password"
	_, err := generateRegistryAuth(registryURL, username, password)

	if err == nil {
		t.Error("We were expecting an error as 'username' was empty, but we didnt get one")
		t.Fail()
	}
}

func Test_GenerateRegistryAuthNoPassword(t *testing.T) {
	registryURL := "https://index.docker.io/v1"
	username := "docker_user"
	password := ""
	_, err := generateRegistryAuth(registryURL, username, password)

	if err == nil {
		t.Error("We were expecting an error as 'password' was empty, but we didnt get one")
		t.Fail()
	}
}

func Test_GenerateRegistryAuthNoInputs(t *testing.T) {
	registryURL := ""
	username := ""
	password := ""
	_, err := generateRegistryAuth(registryURL, username, password)

	if err == nil {
		t.Error("We were expecting an error as all inputs were empty, but we didnt get one")
		t.Fail()
	}
}

func Test_GenerateECRRegistryAuth(t *testing.T) {
	region := "eu-west-2"
	accountId := "1234567"
	gotBytes, _ := generateECRRegistryAuth(accountId, region)
	got, err := bytesToECRStruct(gotBytes)

	if err != nil {
		t.Errorf("Error converting bytes to struct, %q", err)
		t.Fail()
	}

	want := ECRRegistryAuth{
		CredsStore: "ecr-login",
		CredHelpers: map[string]string{
			fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", accountId, region): "ecr-login",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("Structs were not equal, the generated bytes did not match the expected output")
		t.Fail()
	}
}

func Test_GenerateECRRegistryAuthNoRegion(t *testing.T) {
	region := ""
	accountId := "1234567"
	_, err := generateECRRegistryAuth(accountId, region)

	if err == nil {
		t.Error("We were expecting an error as 'region' was empty, but we didnt get one")
		t.Fail()
	}
}

func Test_GenerateECRRegistryAuthNoAccountId(t *testing.T) {
	region := "eu-west-2"
	accountId := ""
	_, err := generateECRRegistryAuth(accountId, region)

	if err == nil {
		t.Error("We were expecting an error as 'accountId' was empty, but we didnt get one")
		t.Fail()
	}
}

func Test_GenerateECRRegistryAuthNoAccountIdOrRegion(t *testing.T) {
	region := ""
	accountId := ""
	_, err := generateECRRegistryAuth(accountId, region)

	if err == nil {
		t.Error("We were expecting an error as 'accountId' and 'region' were empty, but we didnt get one")
		t.Fail()
	}
}

func bytesToECRStruct(bytes []byte) (ECRRegistryAuth, error) {
	obj := ECRRegistryAuth{}
	err := json.Unmarshal(bytes, &obj)

	return obj, err
}

func bytesToRegistryStruct(bytes []byte) (RegistryAuth, error) {
	obj := RegistryAuth{}
	err := json.Unmarshal(bytes, &obj)

	return obj, err
}
