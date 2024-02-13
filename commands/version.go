// Copyright (c) Forge4Flow DAO LLC 2024. All rights reserved.
// Licensed under the MIT license.

package commands

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"os"

	"github.com/alexellis/arkade/pkg/get"
	"github.com/forge4flow/forge-cli/proxy"
	"github.com/forge4flow/forge-cli/stack"
	"github.com/forge4flow/forge-cli/version"
	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

// GitCommit injected at build-time
var (
	shortVersion bool
	warnUpdate   bool
)

func init() {
	versionCmd.Flags().BoolVar(&shortVersion, "short-version", false, "Just print Git SHA")
	versionCmd.Flags().StringVarP(&gateway, "gateway", "g", defaultGateway, "Gateway URL starting with http(s)://")
	versionCmd.Flags().BoolVar(&tlsInsecure, "tls-no-verify", false, "Disable TLS validation")
	versionCmd.Flags().BoolVar(&envsubst, "envsubst", true, "Substitute environment variables in functions.yml file")

	versionCmd.Flags().BoolVar(&warnUpdate, "warn-update", true, "Check for new version and warn about updating")

	versionCmd.Flags().StringVarP(&token, "token", "k", "", "Pass a JWT token to use instead of basic auth")
	forgeCmd.AddCommand(versionCmd)
}

// versionCmd displays version information
var versionCmd = &cobra.Command{
	Use:   "version [--short-version] [--gateway GATEWAY_URL]",
	Short: "Display the clients version information",
	Long: fmt.Sprintf(`The version command returns the current clients version information.

This currently consists of the GitSHA from which the client was built.
- https://github.com/forge4flow/forge-cli/tree/%s`, version.GitCommit),
	Example: `  forge-cli version
  forge-cli version --short-version`,
	RunE: runVersionE,
}

func runVersionE(cmd *cobra.Command, args []string) error {
	if shortVersion {
		fmt.Println(version.BuildVersion())
		return nil
	}

	printLogo()
	fmt.Printf(`CLI:
 commit:  %s
 version: %s
`, version.GitCommit, version.BuildVersion())
	printServerVersions()

	if warnUpdate {
		version := version.Version
		latest, err := get.FindGitHubRelease("openfaas", "forge-cli")
		if err != nil {
			return fmt.Errorf("unable to find latest version online error: %s", err.Error())
		}

		if version != "" && version != latest {
			fmt.Printf("Your forge-cli version (%s) may be out of date. Version: %s is now available on GitHub.\n", version, latest)
		}
	}

	return nil
}

func printServerVersions() error {

	var services stack.Services
	var gatewayAddress string
	var yamlGateway string
	if len(yamlFile) > 0 {
		parsedServices, err := stack.ParseYAMLFile(yamlFile, regex, filter, envsubst)
		if err == nil && parsedServices != nil {
			services = *parsedServices
			yamlGateway = services.Provider.GatewayURL
		}
	}

	gatewayAddress = getGatewayURL(gateway, defaultGateway, yamlGateway, os.Getenv(openFaaSURLEnvironment))

	versionTimeout := 5 * time.Second
	cliAuth, err := proxy.NewCLIAuth(token, gatewayAddress)
	if err != nil {
		return err
	}
	transport := GetDefaultCLITransport(tlsInsecure, &versionTimeout)
	cliClient, err := proxy.NewClient(cliAuth, gatewayAddress, transport, &versionTimeout)
	if err != nil {
		return err
	}
	gatewayInfo, err := cliClient.GetSystemInfo(context.Background())
	if err != nil {
		return err
	}

	printGatewayDetails(gatewayAddress, gatewayInfo.Version.Release, gatewayInfo.Version.SHA)

	fmt.Printf(`
Provider
 name:          %s
 orchestration: %s
 version:       %s 
 sha:           %s
`, gatewayInfo.Provider.Name, gatewayInfo.Provider.Orchestration, gatewayInfo.Provider.Version.Release, gatewayInfo.Provider.Version.SHA)
	return nil
}

func printGatewayDetails(gatewayAddress, version, sha string) {
	fmt.Printf(`
Gateway
 uri:     %s`, gatewayAddress)

	if version != "" {
		fmt.Printf(`
 version: %s
 sha:     %s
`, version, sha)
	}

	fmt.Println()
}

// printLogo prints an ASCII logo, which was generated with figlet
func printLogo() {
	figletColoured := aec.GreenF.Apply(figletStr)
	if runtime.GOOS == "windows" {
		figletColoured = aec.GreenF.Apply(figletStr)
	}
	fmt.Printf(figletColoured)
}

const figletStr = `______                      ___ ______ _               
|  ___|                    /   ||  ___| |              
| |_ ___  _ __ __ _  ___  / /| || |_  | | _____      __
|  _/ _ \| '__/ _  |/ _ \/ /_| ||  _| | |/ _ \ \ /\ / /
| || (_) | | | (_| |  __/\___  || |   | | (_) \ V  V / 
\_| \___/|_|  \__, |\___|    |_/\_|   |_|\___/ \_/\_/  
               __/ |                                   
              |___/                                    

`
