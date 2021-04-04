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
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/devrel-blox/drb/blox"
	"github.com/devrel-blox/drb/config"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build content source into a JSON file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		cobra.CheckErr(err)
		fmt.Println("Preparing data output...")

		// convert markdown to yaml
		cobra.CheckErr(convertModels(cfg))

		// validate yaml
		cobra.CheckErr(validateModels(cfg))

		data, err := aggregateModels(cfg)
		cobra.CheckErr(err)

		f, err := os.Create(path.Join(cfg.Base, cfg.Destination, "data.json"))
		cobra.CheckErr(err)

		defer f.Close()

		bb, err := json.Marshal(data)
		cobra.CheckErr(err)

		_, err = f.Write(bb)
		cobra.CheckErr(err)

		fmt.Printf("wrote: %s\n", path.Join(cfg.Base, cfg.Destination, "data.json"))
	},
}

func aggregateModels(cfg *config.BloxConfig) (map[string][]interface{}, error) {
	data := make(map[string][]interface{})
	fmt.Printf("Loading YAML files...\n")

	for _, model := range blox.GetAllModels() {
		// Attempt to decode all the YAML files with this directory as model
		fmt.Printf("\t model %s: \n\t\tsource: %s\n", model.ID, path.Join(cfg.Base, cfg.Destination, model.Folder))

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
				slug := strings.Replace(filepath.Base(path), ext, "", -1)

				// if ext != cfg.DefaultExtension {
				// Should be SupportedExtensions?
				if ext != ".yaml" && ext != ".yml" {
					return nil
				}

				entity, err := blox.FromYAML(path, model.ID, model.Cue)
				if err != nil {
					return err
				}
				data[model.Folder] = append(data[model.ID], entity)
				fmt.Printf("\t\t\t%s '%s' loaded\n", model.ID, slug)

				return nil

			})
	}

	return data, nil
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
