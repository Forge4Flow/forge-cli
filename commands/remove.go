// Copyright (c) Forge4Flow DAO LLC 2024. All rights reserved.
// Licensed under the MIT license.

package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/forge4flow/forge-cli/proxy"
	"github.com/forge4flow/forge-cli/stack"
	"github.com/spf13/cobra"
)

func init() {
	// Setup flags that are used by multiple commands (variables defined in faas.go)
	removeCmd.Flags().StringVarP(&gateway, "gateway", "g", defaultGateway, "Gateway URL starting with http(s)://")
	removeCmd.Flags().BoolVar(&tlsInsecure, "tls-no-verify", false, "Disable TLS validation")
	removeCmd.Flags().BoolVar(&envsubst, "envsubst", true, "Substitute environment variables in functions.yml file")
	removeCmd.Flags().StringVarP(&token, "token", "k", "", "Pass a JWT token to use instead of basic auth")
	removeCmd.Flags().StringVarP(&functionNamespace, "namespace", "n", "", "Namespace of the function")

	forgeCmd.AddCommand(removeCmd)
}

// removeCmd deletes/removes Forge4Flow function containers
var removeCmd = &cobra.Command{
	Use: `remove FUNCTION_NAME [--gateway GATEWAY_URL]
  forge-cli remove -f YAML_FILE [--regex "REGEX"] [--filter "WILDCARD"]`,
	Aliases: []string{"rm", "delete"},
	Short:   "Remove deployed Forge4Flow functions",
	Long: `Removes/deletes deployed Forge4Flow functions either via the supplied YAML config
using the "--yaml" flag (which may contain multiple function definitions), or by
explicitly specifying a function name.`,
	Example: `  forge-cli remove -f https://domain/path/myfunctions.yml
  forge-cli remove -f ./functions.yml
  forge-cli remove -f ./functions.yml --filter "*gif*"
  forge-cli remove -f ./functions.yml --regex "fn[0-9]_.*"
  forge-cli remove url-ping
  forge-cli remove img2ansi --gateway==http://remote-site.com:8080`,
	RunE: runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	var services stack.Services
	var gatewayAddress string
	var yamlGateway string
	if len(yamlFile) > 0 && len(args) == 0 {
		parsedServices, err := stack.ParseYAMLFile(yamlFile, regex, filter, envsubst)
		if err != nil {
			return err
		}

		if parsedServices != nil {
			services = *parsedServices
			yamlGateway = services.Provider.GatewayURL
		}
	}

	gatewayAddress = getGatewayURL(gateway, defaultGateway, yamlGateway, os.Getenv(openFaaSURLEnvironment))

	cliAuth, err := proxy.NewCLIAuth(token, gatewayAddress)
	if err != nil {
		return err
	}
	transport := GetDefaultCLITransport(tlsInsecure, &commandTimeout)
	proxyclient, err := proxy.NewClient(cliAuth, gatewayAddress, transport, &commandTimeout)
	if err != nil {
		return err
	}
	ctx := context.Background()

	if len(services.Functions) > 0 {

		for k, function := range services.Functions {
			function.Namespace = getNamespace(functionNamespace, function.Namespace)
			function.Name = k
			fmt.Printf("Deleting: %s.%s\n", function.Name, function.Namespace)

			proxyclient.DeleteFunction(ctx, function.Name, function.Namespace)
		}
	} else {
		if len(args) < 1 {
			return fmt.Errorf("please provide the name of a function to delete")
		}

		functionName = args[0]
		fmt.Printf("Deleting: %s.%s\n", functionName, functionNamespace)
		err := proxyclient.DeleteFunction(ctx, functionName, functionNamespace)
		if err != nil {
			return err
		}
	}

	return nil
}
