// Copyright (c) Forge4Flow DAO LLC 2024. All rights reserved.
// Licensed under the MIT license.

package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/forge4flow/forge-cli/builder"
	"github.com/forge4flow/forge-cli/stack"
	"github.com/spf13/cobra"
)

var (
	list          bool
	quiet         bool
	memoryLimit   string
	cpuLimit      string
	memoryRequest string
	cpuRequest    string
)

const functionsFileName = "functions.yml"

func init() {
	newFunctionCmd.Flags().StringVar(&language, "lang", "", "Language or template to use")
	newFunctionCmd.Flags().StringVarP(&gateway, "gateway", "g", defaultGateway, "Gateway URL to store in YAML stack file")
	newFunctionCmd.Flags().StringVar(&handlerDir, "handler", "", "directory the handler will be written to")
	newFunctionCmd.Flags().StringVarP(&imagePrefix, "prefix", "p", "", "Set prefix for the function image")

	newFunctionCmd.Flags().StringVar(&memoryLimit, "memory-limit", "", "Set a limit for the memory")
	newFunctionCmd.Flags().StringVar(&cpuLimit, "cpu-limit", "", "Set a limit for the CPU")

	newFunctionCmd.Flags().StringVar(&memoryRequest, "memory-request", "", "Set a request or the memory")
	newFunctionCmd.Flags().StringVar(&cpuRequest, "cpu-request", "", "Set a request value for the CPU")

	newFunctionCmd.Flags().BoolVar(&list, "list", false, "List available languages")
	newFunctionCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Skip template notes")

	forgeCmd.AddCommand(newFunctionCmd)
}

// newFunctionCmd displays newFunction information
var newFunctionCmd = &cobra.Command{
	Use:   "new FUNCTION_NAME --lang=FUNCTION_LANGUAGE [--gateway=http://host:port] | --list",
	Short: "Create a new template in the current folder with the name given as name",
	Long: `The new command creates a new function based upon hello-world in the given
language or type in --list for a list of languages available.`,
	Example: `  forge-cli new chatbot --lang node
  forge-cli new text-parser --lang python --quiet
  forge-cli new text-parser --lang python --gateway http://mydomain:8080
  forge-cli new --list`,
	PreRunE: preRunNewFunction,
	RunE:    runNewFunction,
}

// validateFunctionName provides least-common-denominator validation - i.e. only allows valid Kubernetes services names
func validateFunctionName(functionName string) error {
	// Regex for RFC-1123 validation:
	// 	k8s.io/kubernetes/pkg/util/validation/validation.go
	var validDNS = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	if matched := validDNS.MatchString(functionName); !matched {
		return fmt.Errorf(`function name can only contain a-z, 0-9 and dashes`)
	}
	return nil
}

// preRunNewFunction validates args & flags
func preRunNewFunction(cmd *cobra.Command, args []string) error {
	if list {
		return nil
	}

	language, _ = validateLanguageFlag(language)

	if len(language) == 0 && len(args) < 1 {
		cmd.Help()
		os.Exit(0)
	}
	if len(language) == 0 {
		return fmt.Errorf("you must supply a function language with the --lang flag")
	}

	if len(args) < 1 {
		return fmt.Errorf(`please provide a name for the function`)
	}

	functionName = args[0]

	if err := validateFunctionName(functionName); err != nil {
		return err
	}

	return nil
}

func runNewFunction(cmd *cobra.Command, args []string) error {
	if list == true {
		var availableTemplates []string

		templateFolders, err := ioutil.ReadDir(templateDirectory)
		if err != nil {
			return fmt.Errorf(`no language templates were found.

Download templates:
  forge-cli template pull           download the default templates
  forge-cli template store list     view the community template store`)
		}

		for _, file := range templateFolders {
			if file.IsDir() {
				availableTemplates = append(availableTemplates, file.Name())
			}
		}

		fmt.Printf("Languages available as templates:\n%s\n", printAvailableTemplates(availableTemplates))

		return nil
	}

	templateAddress := getTemplateURL("", os.Getenv(templateURLEnvironment), DefaultTemplateRepository)
	pullTemplates(templateAddress)

	if !stack.IsValidTemplate(language) {
		return fmt.Errorf("template: \"%s\" was not found in the templates directory", language)
	}

	var fileName, outputMsg string

	// Verify handerDir is set
	if handlerDir == "" {
		handlerDir = functionName
	}

	if _, err := os.Stat(handlerDir); err == nil {
		return fmt.Errorf("folder: %s already exists", handlerDir)
	}

	if _, err := os.Stat(functionsFileName); os.IsNotExist(err) {
		// If functions.yml doesn't exist, create it
		fileName = functionsFileName
		outputMsg = fmt.Sprintf("Stack file written: %s\n", fileName)
	} else {
		fileName = functionsFileName
		outputMsg = fmt.Sprintf("Stack file updated: %s\n", fileName)
	}

	if err := os.Mkdir(handlerDir, 0700); err != nil {
		return fmt.Errorf("folder: could not create %s : %s", handlerDir, err)
	}
	fmt.Printf("Folder: %s created.\n", handlerDir)

	if err := updateGitignore(); err != nil {
		return fmt.Errorf("got unexpected error while updating .gitignore file: %s", err)
	}

	pathToTemplateYAML := fmt.Sprintf("./template/%s/template.yml", language)
	if _, err := os.Stat(pathToTemplateYAML); os.IsNotExist(err) {
		return err
	}

	langTemplate, err := stack.ParseYAMLForLanguageTemplate(pathToTemplateYAML)
	if err != nil {
		return fmt.Errorf("error reading language template: %s", err.Error())
	}

	templateHandlerFolder := "function"
	if len(langTemplate.HandlerFolder) > 0 {
		templateHandlerFolder = langTemplate.HandlerFolder
	}

	fromTemplateHandler := filepath.Join("template", language, templateHandlerFolder)

	// Create function directory from template.
	builder.CopyFiles(fromTemplateHandler, handlerDir)
	printLogo()
	fmt.Printf("\nFunction created in folder: %s\n", handlerDir)

	imageName := fmt.Sprintf("%s:latest", functionName)

	imagePrefixVal := getPrefixValue()

	if imagePrefixVal = strings.TrimSpace(imagePrefixVal); len(imagePrefixVal) > 0 {
		imageName = fmt.Sprintf("%s/%s", imagePrefixVal, imageName)
	}

	function := stack.Function{
		Name:     functionName,
		Handler:  "./" + handlerDir,
		Language: language,
		Image:    imageName,
	}

	if len(memoryLimit) > 0 || len(cpuLimit) > 0 {
		function.Limits = &stack.FunctionResources{
			CPU:    cpuLimit,
			Memory: memoryLimit,
		}
	}

	if len(memoryRequest) > 0 || len(cpuRequest) > 0 {
		function.Requests = &stack.FunctionResources{
			CPU:    cpuRequest,
			Memory: memoryRequest,
		}
	}

	yamlContent := prepareYAMLContent(gateway, &function)

	f, err := os.OpenFile("./"+fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("could not open file '%s' %s", fileName, err)
	}

	_, stackWriteErr := f.Write([]byte(yamlContent))
	if stackWriteErr != nil {
		return fmt.Errorf("error writing stack file %s", stackWriteErr)
	}

	fmt.Print(outputMsg)

	if !quiet {
		languageTemplate, _ := stack.LoadLanguageTemplate(language)

		if languageTemplate.WelcomeMessage != "" {
			fmt.Printf("\nNotes:\n")
			fmt.Printf("%s\n", languageTemplate.WelcomeMessage)
		}
	}

	return nil
}

func getPrefixValue() string {
	prefix := ""
	if len(imagePrefix) > 0 {
		return imagePrefix
	}

	if val, ok := os.LookupEnv("OPENFAAS_PREFIX"); ok && len(val) > 0 {
		prefix = val
	}
	return prefix
}

func prepareYAMLContent(gateway string, function *stack.Function) (yamlContent string) {

	yamlContent = `  ` + function.Name + `:
    lang: ` + function.Language + `
    handler: ` + function.Handler + `
    image: ` + function.Image + `
`

	if function.Requests != nil && (len(function.Requests.CPU) > 0 || len(function.Requests.Memory) > 0) {
		yamlContent += "    requests:\n"
		if len(function.Requests.CPU) > 0 {
			yamlContent += `      cpu: ` + function.Requests.CPU + "\n"
		}

		if len(function.Requests.Memory) > 0 {
			yamlContent += `      memory: ` + function.Requests.Memory + "\n"
		}
	}

	if function.Limits != nil && (len(function.Limits.CPU) > 0 || len(function.Limits.Memory) > 0) {
		yamlContent += "    limits:\n"
		if len(function.Limits.CPU) > 0 {
			yamlContent += `      cpu: ` + function.Limits.CPU + "\n"
		}

		if len(function.Limits.Memory) > 0 {
			yamlContent += `      memory: ` + function.Limits.Memory + "\n"
		}
	}

	yamlContent += "\n"
	yamlContent = `version: ` + defaultSchemaVersion + `
provider:
  name: functions4flow
  gateway: ` + gateway + `
functions:
` + yamlContent

	return yamlContent
}

func printAvailableTemplates(availableTemplates []string) string {
	var result string
	sort.Slice(availableTemplates, func(i, j int) bool {
		return availableTemplates[i] < availableTemplates[j]
	})
	for _, template := range availableTemplates {
		result += fmt.Sprintf("- %s\n", template)
	}
	return result
}

func duplicateFunctionName(functionName string, appendFile string) error {
	fileBytes, readErr := os.ReadFile(appendFile)
	if readErr != nil {
		return fmt.Errorf("unable to read %s to append, %s", appendFile, readErr)
	}

	services, parseErr := stack.ParseYAMLData(fileBytes, "", "", envsubst)

	if parseErr != nil {
		return fmt.Errorf("Error parsing %s yml file", appendFile)
	}

	if _, exists := services.Functions[functionName]; exists {
		return fmt.Errorf(`
Function %s already exists in %s file. 
Cannot have duplicate function names in same yaml file`, functionName, appendFile)
	}

	return nil
}
