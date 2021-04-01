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
	"strings"

	"github.com/devrel-blox/drb/blox"
	"github.com/devrel-blox/drb/config"
	"github.com/devrel-blox/drb/encoding/markdown"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert markdown to yaml",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Load()
		cobra.CheckErr(err)
		cobra.CheckErr(convertModels(cfg))
	},
}

func convertModels(cfg *config.BloxConfig) error {
	var errors error

	for _, model := range blox.Models {
		// Attempt to decode all the YAML files with this directory as model
		fmt.Printf("Checking for %s markdown files in %s\n", model.ID, path.Join(cfg.Base, cfg.Source, model.Folder))

		filepath.Walk(path.Join(cfg.Base, cfg.Source, model.Folder),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					// Squash, we've not even validated that it's a supported ext
					return nil
				}

				if info.IsDir() {
					return nil
				}

				ext := filepath.Ext(path)
				slug := strings.Replace(filepath.Base(path), ext, "", -1)
				// if ext != cfg.DefaultExtension {
				// Should be SupportedExtensions?
				if ext != ".md" && ext != ".mdx" {
					return nil
				}
				f, err := os.Open(path)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, path))
					return nil
				}
				bb, err := os.ReadFile(path)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, path))
					return nil
				}
				f.Close()
				md, err := markdown.ToYAML(string(bb))
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, path))
					return nil
				}
				err = os.MkdirAll(model.DestinationContentPath(), 0755)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, path))
					return nil
				}
				mdf, err := os.Create(model.DestinationFilePath(slug))
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, path))
					return nil
				}
				_, err = mdf.WriteString(md)
				if err != nil {
					errors = multierror.Append(errors, multierror.Prefix(err, path))
					return nil
				}
				mdf.Close()
				/*	profile, err := profile.LoadFromYAML(path)
					if err != nil {
						failedModels[path] = err
						return nil
					}

				*/

				fmt.Println(fmt.Sprintf("Markdown file '%s' converted", slug))

				return nil

				// modelYaml, err := ioutil.ReadFile(path)

				// if err != nil {
				// 	failedModels[path] = err
				// 	return nil
				// }

				// var profile blox.Profile

				// err = yaml.Unmarshal(modelYaml, &profile)

				// if err != nil {
				// 	failedModels[path] = err
				// 	return nil
				// }

				// if err := profile.Validate(); err != nil {
				// 	failedModels[path] = err
				// 	return nil
				// }

				// fmt.Println("Valid ", model.Name, ": ", path)
				return nil
			})
	}

	return errors
}
func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
