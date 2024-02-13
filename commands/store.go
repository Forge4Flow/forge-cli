// Copyright (c) Forge4Flow DAO LLC 2024. All rights reserved.
// Licensed under the MIT license.

package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/forge4flow/forge-cli/proxy"
	storeV2 "github.com/forge4flow/forge-cli/schema/store/v2"
	"github.com/spf13/cobra"
)

var (
	storeAddress     string
	verbose          bool
	storeDeployFlags DeployFlags
	//Platform platform variable set at build time
	Platform string
	// if the CLI is built using buildx, then the Platform value needs to be mapped to
	// one of the supported values used in the store.
	shortPlatform = map[string]string{
		"linux/arm/v6": "armhf",
		"linux/amd64":  "x86_64",
		"linux/arm64":  "arm64",
	}
)

const (
	defaultStore      = "https://raw.githubusercontent.com/openfaas/store/master/functions.json"
	maxDescriptionLen = 40
)

var platformValue string

func init() {
	storeCmd.PersistentFlags().StringVarP(&storeAddress, "url", "u", defaultStore, "Alternative Store URL starting with http(s)://")
	storeCmd.PersistentFlags().StringVarP(&platformValue, "platform", "p", Platform, "Target platform for store")

	forgeCmd.AddCommand(storeCmd)
}

var storeCmd = &cobra.Command{
	Use:   `store`,
	Short: "Forge4Flow store commands",
	Long:  "Allows browsing and deploying Forge4Flow functions from a store",
}

func storeList(store string) ([]storeV2.StoreFunction, error) {

	var storeData storeV2.Store

	store = strings.TrimRight(store, "/")

	timeout := 60 * time.Second
	tlsInsecure := false

	client := proxy.MakeHTTPClient(&timeout, tlsInsecure)

	res, err := client.Get(store)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Forge4Flow store at URL: %s", store)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK:
		bytesOut, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read result from Forge4Flow store at URL: %s", store)
		}

		jsonErr := json.Unmarshal(bytesOut, &storeData)
		if jsonErr != nil {
			return nil, fmt.Errorf("cannot parse result from Forge4Flow store at URL: %s\n%s", store, jsonErr.Error())
		}
	default:
		bytesOut, err := io.ReadAll(res.Body)
		if err == nil {
			return nil, fmt.Errorf("server returned unexpected status code: %d - %s", res.StatusCode, string(bytesOut))
		}
	}

	return storeData.Functions, nil
}

func filterStoreList(functions []storeV2.StoreFunction, platform string) []storeV2.StoreFunction {
	var filteredList []storeV2.StoreFunction

	for _, function := range functions {

		_, ok := getValueIgnoreCase(function.Images, platform)

		if ok {
			filteredList = append(filteredList, function)
		}
	}

	return filteredList
}

// getValueIgnoreCase get a key value from map by ignoring case for key
func getValueIgnoreCase(kv map[string]string, key string) (string, bool) {
	for k, v := range kv {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}
	return "", false
}

func storeFindFunction(functionName string, storeItems []storeV2.StoreFunction) *storeV2.StoreFunction {
	var item storeV2.StoreFunction

	for _, item = range storeItems {
		if item.Name == functionName || item.Title == functionName {
			return &item
		}
	}

	return nil
}

func getPlatform() string {
	if len(Platform) == 0 {
		return mainPlatform
	}
	return Platform
}

func getTargetPlatform(inputPlatform string) string {
	if len(inputPlatform) == 0 {
		currentPlatform := getPlatform()
		target, ok := shortPlatform[currentPlatform]
		if ok {
			return target
		}
		return currentPlatform
	}
	return inputPlatform
}

func getStorePlatforms(functions []storeV2.StoreFunction) []string {
	var distinctPlatformMap = make(map[string]bool)
	var result []string

	for _, function := range functions {
		for key := range function.Images {
			_, exists := distinctPlatformMap[key]

			if !exists {
				distinctPlatformMap[key] = true
				result = append(result, key)
			}
		}
	}

	return result
}
