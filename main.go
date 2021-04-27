// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	"github.com/mesh-for-data/openapi2crd/pkg/config"
	"github.com/mesh-for-data/openapi2crd/pkg/exporter"
	"github.com/mesh-for-data/openapi2crd/pkg/generator"
)

const (
	outputOption    = "output"
	resourcesOption = "input"
	gvkOption       = "gvk"
)

// RootCmd defines the root cli command
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "openapi2crd SPEC_FILE",
		Short:         "Generate CustomResourceDefinition from OpenAPI 3.0 document",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			specOptionValue := args[0]

			openapiSpec, err := config.LoadOpenAPI(specOptionValue)
			if err != nil {
				return err
			}

			crds := []*apiextensions.CustomResourceDefinition{}

			gvkOptionValues := viper.GetStringSlice(gvkOption)
			if len(gvkOptionValues) != 0 {
				loaded, err := config.GenerateCRDs(gvkOptionValues)
				if err != nil {
					return err
				}
				crds = append(crds, loaded...)
			}

			resourcesOptionValue := viper.GetString(resourcesOption)
			if resourcesOptionValue != "" {
				loaded, err := config.LoadCRDs(resourcesOptionValue)
				if err != nil {
					return err
				}
				crds = append(crds, loaded...)
			}

			if len(crds) == 0 {
				message := fmt.Sprintf("You must pass flags --%s or --%s", gvkOption, resourcesOption)
				if resourcesOptionValue != "" {
					message = fmt.Sprintf("Does directory %s include a YAML with CustomResourceDefinition?", resourcesOptionValue)
				}
				return fmt.Errorf("nothing to process. %s", message)
			}

			outputOptionValue := viper.GetString(outputOption)
			exporter, err := exporter.New(outputOptionValue)
			if err != nil {
				return err
			}

			generator := generator.New()
			for _, crd := range crds {
				modified, err := generator.Generate(crd, openapiSpec.Components.Schemas)
				if err != nil {
					return err
				}
				err = exporter.Export(modified)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringP(outputOption, "o", "", "Path to output file (required)")
	_ = cmd.MarkFlagRequired(outputOption)
	cmd.Flags().StringP(resourcesOption, "i", "",
		"Path to a directory with CustomResourceDefinition YAML files (required unless -g is used)")
	cmd.Flags().StringSliceP(gvkOption, "g", []string{},
		"The group/version/kind to create (can be specified zero or more times)")

	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	// Run the cli
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
