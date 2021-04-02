/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/devrel-blox/drb/blox"
	"github.com/devrel-blox/drb/config"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		cobra.CheckErr(err)

		cobra.CheckErr(validateModels(cfg))
	},
}

func validateModels(cfg *config.BloxConfig) error {
	var errors error

	// We want to validate all the YAML for the models that we're aware of.
	for _, model := range blox.Models {
		// Attempt to decode all the YAML files with this directory as model
		fmt.Printf("Validating %s YAML files in %s\n", model.ID, path.Join(cfg.Base, cfg.Destination, model.Folder))

		filepath.Walk(path.Join(cfg.Base, cfg.Destination, model.Folder),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					// Squash, we've not even validated that it's a supported ext
					return nil
				}

				if info.IsDir() {
					return nil
				}

				ext := filepath.Ext(path)

				// if ext != cfg.DefaultExtension {
				// Should be SupportedExtensions?
				if ext != ".yaml" && ext != ".yml" {
					return nil
				}

				switch model.ID {
				case "profile":
					{
						// Check for Replace
						profileSchema := blox.ProfileCue
						if replace, ok := cfg.SchemaOverrides.Replace["profile"]; ok {
							profileSchema = replace
						}
						// Check for Extend
						// No access to Merge through cuego. Merge through other APIs
						// is deprecated. Need to investigate
						_, err := blox.ProfileFromYAML(path, profileSchema)
						if err != nil {
							errors = multierror.Append(errors, multierror.Prefix(err, path))
							return nil
						}

					}
				case "article":
					{
						_, err := blox.ArticleFromYAML(path)

						if err != nil {
							errors = multierror.Append(errors, multierror.Prefix(err, path))
							return nil
						}
					}
				case "category":
					{
						_, err := blox.CategoryFromYAML(path)

						if err != nil {
							errors = multierror.Append(errors, multierror.Prefix(err, path))
							return nil
						}
					}
				}

				return nil

			})
	}

	return errors
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
