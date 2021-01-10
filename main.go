package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mesh-for-data/openapi2crd/pkg/config"
	"github.com/mesh-for-data/openapi2crd/pkg/exporter"
	"github.com/mesh-for-data/openapi2crd/pkg/generator"
)

const (
	outputOption    = "output"
	resourcesOption = "input"
)

// RootCmd defines the root cli command
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "openapi2crd FILE",
		Short:         "Generate CustomResourceDefinition from OpenAPI 3.0 document",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			specOptionValue := args[0]

			swagger, err := config.LoadSwagger(specOptionValue)
			if err != nil {
				return err
			}

			resourcesOptionValue := viper.GetString(resourcesOption)
			crds, err := config.LoadCRDs(resourcesOptionValue)
			if err != nil {
				return err
			}

			outputOptionValue := viper.GetString(outputOption)
			exporter, err := exporter.New(outputOptionValue)
			if err != nil {
				return err
			}

			generator := generator.New()
			for _, crd := range crds {
				modified := generator.Generate(crd, swagger.Components.Schemas)
				err := exporter.Export(modified)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringP(outputOption, "o", "", "Path to output file (required)")
	_ = cmd.MarkFlagRequired(outputOption)
	cmd.Flags().StringP(resourcesOption, "i", "", "Path to directory with CustomResourceDefinition YAML files (required)")
	_ = cmd.MarkFlagRequired(resourcesOption)

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
